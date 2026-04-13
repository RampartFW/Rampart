package cluster

import (
	"context"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/rampartfw/rampart/internal/model"
)

// RaftNode represents a single node in the Raft cluster.
type RaftNode struct {
	mu sync.RWMutex

	id          string
	state       model.NodeState
	currentTerm uint64
	votedFor    string
	log         *Log
	commitIndex uint64
	lastApplied uint64

	// Leader-only state
	nextIndex  map[string]uint64
	matchIndex map[string]uint64

	transport Transport
	peers     map[string]string // id -> address
	fsm       FSM

	electionTimer  *time.Timer
	heartbeatTimer *time.Timer

	proposals map[uint64]chan error

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// Transport defines the interface for Raft RPCs.
type Transport interface {
	SendRequestVote(ctx context.Context, target string, req RequestVoteRequest) (RequestVoteResponse, error)
	SendAppendEntries(ctx context.Context, target string, req AppendEntriesRequest) (AppendEntriesResponse, error)
	Listen(address string, handler RPCHandler) error
	Close() error
}

// RPCHandler handles incoming Raft RPCs.
type RPCHandler interface {
	HandleRequestVote(req RequestVoteRequest) (RequestVoteResponse, error)
	HandleAppendEntries(req AppendEntriesRequest) (AppendEntriesResponse, error)
}

// FSM defines the interface for the Raft finite state machine.
type FSM interface {
	Apply(entry model.LogEntry) error
	Snapshot() ([]byte, error)
	Restore(snapshot []byte) error
}

// RPC Request/Response structures

type RequestVoteRequest struct {
	Term         uint64
	CandidateID  string
	LastLogIndex uint64
	LastLogTerm  uint64
}

type RequestVoteResponse struct {
	Term        uint64
	VoteGranted bool
}

type AppendEntriesRequest struct {
	Term         uint64
	LeaderID     string
	PrevLogIndex uint64
	PrevLogTerm  uint64
	Entries      []model.LogEntry
	LeaderCommit uint64
}

type AppendEntriesResponse struct {
	Term    uint64
	Success bool
}

func init() {
	gob.Register(RequestVoteRequest{})
	gob.Register(RequestVoteResponse{})
	gob.Register(AppendEntriesRequest{})
	gob.Register(AppendEntriesResponse{})
	gob.Register(model.LogEntry{})
}

// NewRaftNode creates a new Raft node.
func NewRaftNode(id string, peers map[string]string, transport Transport, log *Log, fsm FSM) *RaftNode {
	ctx, cancel := context.WithCancel(context.Background())
	return &RaftNode{
		id:         id,
		state:      model.StateFollower,
		peers:      peers,
		transport:  transport,
		log:        log,
		fsm:        fsm,
		nextIndex:  make(map[string]uint64),
		matchIndex: make(map[string]uint64),
		proposals:  make(map[uint64]chan error),
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Start starts the Raft node.
func (rn *RaftNode) Start(address string) error {
	if err := rn.transport.Listen(address, rn); err != nil {
		return err
	}

	rn.mu.Lock()
	rn.resetElectionTimer()
	rn.mu.Unlock()

	rn.wg.Add(1)
	go rn.run()
	return nil
}

func (rn *RaftNode) run() {
	defer rn.wg.Done()
	for {
		select {
		case <-rn.ctx.Done():
			return
		case <-rn.electionTimer.C:
			rn.startElection()
		}
	}
}

func (rn *RaftNode) resetElectionTimer() {
	if rn.electionTimer != nil {
		rn.electionTimer.Stop()
	}
	timeout := time.Duration(150+rn.randomInt(150)) * time.Millisecond
	rn.electionTimer = time.NewTimer(timeout)
}

func (rn *RaftNode) randomInt(n int64) int64 {
	val, _ := rand.Int(rand.Reader, big.NewInt(n))
	return val.Int64()
}

func (rn *RaftNode) startElection() {
	rn.mu.Lock()
	rn.state = model.StateCandidate
	rn.currentTerm++
	rn.votedFor = rn.id
	term := rn.currentTerm
	lastLog := rn.log.LastEntry()
	rn.resetElectionTimer()
	rn.mu.Unlock()

	votes := 1
	var votesMu sync.Mutex

	for peerID, addr := range rn.peers {
		if peerID == rn.id {
			continue
		}

		rn.wg.Add(1)
		go func(peerID, addr string) {
			defer rn.wg.Done()
			req := RequestVoteRequest{
				Term:         term,
				CandidateID:  rn.id,
				LastLogIndex: lastLog.Index,
				LastLogTerm:  lastLog.Term,
			}

			resp, err := rn.transport.SendRequestVote(rn.ctx, addr, req)
			if err != nil {
				return
			}

			rn.mu.Lock()
			defer rn.mu.Unlock()
			if resp.Term > rn.currentTerm {
				rn.stepDown(resp.Term)
				return
			}

			if rn.state == model.StateCandidate && rn.currentTerm == term && resp.VoteGranted {
				votesMu.Lock()
				votes++
				if votes > (len(rn.peers)/2) {
					rn.becomeLeader()
				}
				votesMu.Unlock()
			}
		}(peerID, addr)
	}
}

func (rn *RaftNode) stepDown(term uint64) {
	rn.state = model.StateFollower
	rn.currentTerm = term
	rn.votedFor = ""
	rn.resetElectionTimer()
	
	// Notify all pending proposals that we are no longer leader
	for idx, ch := range rn.proposals {
		ch <- fmt.Errorf("stepped down from leader")
		delete(rn.proposals, idx)
	}
}

func (rn *RaftNode) becomeLeader() {
	if rn.state == model.StateLeader {
		return
	}
	rn.state = model.StateLeader
	lastLog := rn.log.LastEntry()
	for peerID := range rn.peers {
		rn.nextIndex[peerID] = lastLog.Index + 1
		rn.matchIndex[peerID] = 0
	}

	if rn.heartbeatTimer != nil {
		rn.heartbeatTimer.Stop()
	}
	rn.heartbeatTimer = time.NewTimer(50 * time.Millisecond)

	rn.wg.Add(1)
	go rn.leaderLoop()
}

func (rn *RaftNode) leaderLoop() {
	defer rn.wg.Done()
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-rn.ctx.Done():
			return
		case <-ticker.C:
			rn.mu.RLock()
			if rn.state != model.StateLeader {
				rn.mu.RUnlock()
				return
			}
			rn.replicateLog()
			rn.mu.RUnlock()
		}
	}
}

func (rn *RaftNode) replicateLog() {
	for peerID, addr := range rn.peers {
		if peerID == rn.id {
			continue
		}

		rn.wg.Add(1)
		go func(peerID, addr string) {
			defer rn.wg.Done()
			rn.mu.RLock()
			if rn.state != model.StateLeader {
				rn.mu.RUnlock()
				return
			}

			ni := rn.nextIndex[peerID]
			entries := rn.log.Entries(ni, 0)
			prevIndex := ni - 1
			prevTerm := uint64(0)
			if prevIndex > 0 {
				if entry, ok := rn.log.Get(prevIndex); ok {
					prevTerm = entry.Term
				}
			}

			req := AppendEntriesRequest{
				Term:         rn.currentTerm,
				LeaderID:     rn.id,
				PrevLogIndex: prevIndex,
				PrevLogTerm:  prevTerm,
				Entries:      entries,
				LeaderCommit: rn.commitIndex,
			}
			rn.mu.RUnlock()

			resp, err := rn.transport.SendAppendEntries(rn.ctx, addr, req)
			if err != nil {
				return
			}

			rn.mu.Lock()
			defer rn.mu.Unlock()

			if resp.Term > rn.currentTerm {
				rn.stepDown(resp.Term)
				return
			}

			if rn.state == model.StateLeader && rn.currentTerm == req.Term {
				if resp.Success {
					rn.matchIndex[peerID] = req.PrevLogIndex + uint64(len(entries))
					rn.nextIndex[peerID] = rn.matchIndex[peerID] + 1
					rn.updateCommitIndex()
				} else {
					if rn.nextIndex[peerID] > 1 {
						rn.nextIndex[peerID]--
					}
				}
			}
		}(peerID, addr)
	}
}

func (rn *RaftNode) updateCommitIndex() {
	for n := rn.log.LastEntry().Index; n > rn.commitIndex; n-- {
		count := 1
		for peerID := range rn.peers {
			if peerID != rn.id && rn.matchIndex[peerID] >= n {
				count++
			}
		}

		if count > len(rn.peers)/2 {
			if entry, ok := rn.log.Get(n); ok && entry.Term == rn.currentTerm {
				rn.commitIndex = n
				rn.applyEntries()
				break
			}
		}
	}
}

func (rn *RaftNode) applyEntries() {
	for rn.lastApplied < rn.commitIndex {
		rn.lastApplied++
		if entry, ok := rn.log.Get(rn.lastApplied); ok {
			err := rn.fsm.Apply(entry)
			
			// Notify proposer if we are the leader
			if rn.state == model.StateLeader {
				if ch, ok := rn.proposals[entry.Index]; ok {
					ch <- err
					delete(rn.proposals, entry.Index)
				}
			}
		}
	}
}

// RPCHandler implementation

func (rn *RaftNode) HandleRequestVote(req RequestVoteRequest) (RequestVoteResponse, error) {
	rn.mu.Lock()
	defer rn.mu.Unlock()

	resp := RequestVoteResponse{
		Term:        rn.currentTerm,
		VoteGranted: false,
	}

	if req.Term < rn.currentTerm {
		return resp, nil
	}

	if req.Term > rn.currentTerm {
		rn.stepDown(req.Term)
	}

	lastLog := rn.log.LastEntry()
	upToDate := req.LastLogTerm > lastLog.Term || (req.LastLogTerm == lastLog.Term && req.LastLogIndex >= lastLog.Index)

	if (rn.votedFor == "" || rn.votedFor == req.CandidateID) && upToDate {
		rn.votedFor = req.CandidateID
		resp.VoteGranted = true
		rn.resetElectionTimer()
	}

	resp.Term = rn.currentTerm
	return resp, nil
}

func (rn *RaftNode) HandleAppendEntries(req AppendEntriesRequest) (AppendEntriesResponse, error) {
	rn.mu.Lock()
	defer rn.mu.Unlock()

	resp := AppendEntriesResponse{
		Term:    rn.currentTerm,
		Success: false,
	}

	if req.Term < rn.currentTerm {
		return resp, nil
	}

	if req.Term > rn.currentTerm || rn.state != model.StateFollower {
		rn.stepDown(req.Term)
	}
	rn.resetElectionTimer()

	// Log consistency check
	if req.PrevLogIndex > 0 {
		entry, ok := rn.log.Get(req.PrevLogIndex)
		if !ok || entry.Term != req.PrevLogTerm {
			resp.Term = rn.currentTerm
			return resp, nil
		}
	}

	// Append entries
	if len(req.Entries) > 0 {
		for i, entry := range req.Entries {
			existing, ok := rn.log.Get(entry.Index)
			if ok && existing.Term != entry.Term {
				rn.log.Truncate(entry.Index)
			}
			if !ok || existing.Term != entry.Term {
				rn.log.Append(req.Entries[i:]...)
				break
			}
		}
	}

	if req.LeaderCommit > rn.commitIndex {
		lastIdx := rn.log.LastEntry().Index
		if req.LeaderCommit < lastIdx {
			rn.commitIndex = req.LeaderCommit
		} else {
			rn.commitIndex = lastIdx
		}
		rn.applyEntries()
	}

	resp.Term = rn.currentTerm
	resp.Success = true
	return resp, nil
}

// Propose proposes a new entry to the cluster.
func (rn *RaftNode) Propose(entryType model.EntryType, data []byte) error {
	rn.mu.Lock()
	if rn.state != model.StateLeader {
		rn.mu.Unlock()
		return fmt.Errorf("not the leader")
	}

	lastLog := rn.log.LastEntry()
	index := lastLog.Index + 1
	entry := model.LogEntry{
		Term:      rn.currentTerm,
		Index:     index,
		Type:      entryType,
		Data:      data,
		Timestamp: time.Now(),
	}

	if err := rn.log.Append(entry); err != nil {
		rn.mu.Unlock()
		return err
	}

	ch := make(chan error, 1)
	rn.proposals[index] = ch
	rn.mu.Unlock()

	// Wait for entry to be committed
	select {
	case err := <-ch:
		return err
	case <-rn.ctx.Done():
		return rn.ctx.Err()
	case <-time.After(5 * time.Second):
		rn.mu.Lock()
		delete(rn.proposals, index)
		rn.mu.Unlock()
		return fmt.Errorf("proposal timed out")
	}
}

// Status returns the current status of the node.
func (rn *RaftNode) Status() model.NodeStatus {
	rn.mu.RLock()
	defer rn.mu.RUnlock()

	return model.NodeStatus{
		ID:        rn.id,
		State:     rn.state,
		IsHealthy: true, // Simplified
	}
}

// Close stops the Raft node.
func (rn *RaftNode) Close() error {
	rn.cancel()
	rn.mu.Lock()
	defer rn.mu.Unlock()
	if rn.electionTimer != nil {
		rn.electionTimer.Stop()
	}
	if rn.heartbeatTimer != nil {
		rn.heartbeatTimer.Stop()
	}
	
	// Cleanup proposals
	for idx, ch := range rn.proposals {
		ch <- fmt.Errorf("closing node")
		delete(rn.proposals, idx)
	}

	rn.wg.Wait()
	return rn.transport.Close()
}
