package cluster

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

// TCPTransport implements the Transport interface using TCP and TLS.
type TCPTransport struct {
	mu sync.Mutex

	listener net.Listener
	handler  RPCHandler
	conns    map[string]net.Conn
	tlsConfig *tls.Config
}

// NewTCPTransport creates a new TCP transport with TLS.
func NewTCPTransport(certFile, keyFile, caFile string) (*TCPTransport, error) {
	// Load node certificate and key
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load node certificate: %w", err)
	}

	// Load CA certificate
	caCert, err := os.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %w", err)
	}
	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to append CA certificate to pool")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caPool,
		ClientCAs:    caPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS12,
	}

	return &TCPTransport{
		conns:     make(map[string]net.Conn),
		tlsConfig: tlsConfig,
	}, nil
}

// Listen starts the TCP listener and handles incoming connections.
func (t *TCPTransport) Listen(address string, handler RPCHandler) error {
	l, err := tls.Listen("tcp", address, t.tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", address, err)
	}

	t.listener = l
	t.handler = handler

	go t.acceptLoop()
	return nil
}

func (t *TCPTransport) acceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			return
		}
		go t.handleConnection(conn)
	}
}

func (t *TCPTransport) handleConnection(conn net.Conn) {
	defer conn.Close()
	decoder := gob.NewDecoder(conn)
	encoder := gob.NewEncoder(conn)

	for {
		var rpcType string
		if err := decoder.Decode(&rpcType); err != nil {
			return
		}

		switch rpcType {
		case "RequestVote":
			var req RequestVoteRequest
			if err := decoder.Decode(&req); err != nil {
				return
			}
			resp, err := t.handler.HandleRequestVote(req)
			if err != nil {
				return
			}
			if err := encoder.Encode(resp); err != nil {
				return
			}
		case "AppendEntries":
			var req AppendEntriesRequest
			if err := decoder.Decode(&req); err != nil {
				return
			}
			resp, err := t.handler.HandleAppendEntries(req)
			if err != nil {
				return
			}
			if err := encoder.Encode(resp); err != nil {
				return
			}
		}
	}
}

// SendRequestVote sends a RequestVote RPC to a target node.
func (t *TCPTransport) SendRequestVote(ctx context.Context, target string, req RequestVoteRequest) (RequestVoteResponse, error) {
	conn, err := t.getConn(target)
	if err != nil {
		return RequestVoteResponse{}, err
	}

	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	if err := encoder.Encode("RequestVote"); err != nil {
		t.closeConn(target)
		return RequestVoteResponse{}, err
	}
	if err := encoder.Encode(req); err != nil {
		t.closeConn(target)
		return RequestVoteResponse{}, err
	}

	var resp RequestVoteResponse
	if err := decoder.Decode(&resp); err != nil {
		t.closeConn(target)
		return RequestVoteResponse{}, err
	}

	return resp, nil
}

// SendAppendEntries sends an AppendEntries RPC to a target node.
func (t *TCPTransport) SendAppendEntries(ctx context.Context, target string, req AppendEntriesRequest) (AppendEntriesResponse, error) {
	conn, err := t.getConn(target)
	if err != nil {
		return AppendEntriesResponse{}, err
	}

	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	if err := encoder.Encode("AppendEntries"); err != nil {
		t.closeConn(target)
		return AppendEntriesResponse{}, err
	}
	if err := encoder.Encode(req); err != nil {
		t.closeConn(target)
		return AppendEntriesResponse{}, err
	}

	var resp AppendEntriesResponse
	if err := decoder.Decode(&resp); err != nil {
		t.closeConn(target)
		return AppendEntriesResponse{}, err
	}

	return resp, nil
}

func (t *TCPTransport) getConn(target string) (net.Conn, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if conn, ok := t.conns[target]; ok {
		return conn, nil
	}

	dialer := &net.Dialer{Timeout: 5 * time.Second}
	conn, err := tls.DialWithDialer(dialer, "tcp", target, t.tlsConfig)
	if err != nil {
		return nil, err
	}

	t.conns[target] = conn
	return conn, nil
}

func (t *TCPTransport) closeConn(target string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if conn, ok := t.conns[target]; ok {
		conn.Close()
		delete(t.conns, target)
	}
}

// Close closes the listener and all active connections.
func (t *TCPTransport) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.listener != nil {
		t.listener.Close()
	}

	for _, conn := range t.conns {
		conn.Close()
	}
	t.conns = make(map[string]net.Conn)

	return nil
}
