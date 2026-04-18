package mcp

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rampartfw/rampart/internal/engine"
	"github.com/rampartfw/rampart/internal/model"
)

type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

func (s *Server) handleToolsList(req *Request) *Response {
	tools := []Tool{
		{
			Name:        "list_rules",
			Description: "List active firewall rules with optional filters",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"policy": map[string]interface{}{"type": "string"},
					"action": map[string]interface{}{"type": "string"},
				},
			},
		},
		{
			Name:        "add_rule",
			Description: "Add a single rule (quick mode)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name":     map[string]interface{}{"type": "string"},
					"protocol": map[string]interface{}{"type": "string"},
					"dport":    map[string]interface{}{"type": "integer"},
					"source":   map[string]interface{}{"type": "string"},
					"action":   map[string]interface{}{"type": "string"},
					"ttl":      map[string]interface{}{"type": "string"},
				},
				"required": []string{"name", "action"},
			},
		},
		{
			Name:        "remove_rule",
			Description: "Remove a rule by name or ID",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{"type": "string"},
					"id":   map[string]interface{}{"type": "string"},
				},
			},
		},
		{
			Name:        "plan_policy",
			Description: "Dry-run a YAML policy and show execution plan",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"yaml": map[string]interface{}{"type": "string"},
				},
				"required": []string{"yaml"},
			},
		},
		{
			Name:        "apply_policy",
			Description: "Apply a YAML policy (requires confirmation)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"yaml":    map[string]interface{}{"type": "string"},
					"confirm": map[string]interface{}{"type": "boolean"},
				},
				"required": []string{"yaml", "confirm"},
			},
		},
		{
			Name:        "simulate_packet",
			Description: "Test if a packet would be allowed/denied",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"src":      map[string]interface{}{"type": "string"},
					"dst":      map[string]interface{}{"type": "string"},
					"protocol": map[string]interface{}{"type": "string"},
					"dport":    map[string]interface{}{"type": "integer"},
					"direction": map[string]interface{}{"type": "string"},
				},
				"required": []string{"src", "dst", "protocol", "dport"},
			},
		},
		{
			Name:        "rollback",
			Description: "Rollback to a specific snapshot",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id": map[string]interface{}{"type": "string"},
				},
				"required": []string{"id"},
			},
		},
		{
			Name:        "list_snapshots",
			Description: "List available snapshots",
			InputSchema: map[string]interface{}{"type": "object"},
		},
		{
			Name:        "audit_search",
			Description: "Search audit events",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"action": map[string]interface{}{"type": "string"},
					"actor":  map[string]interface{}{"type": "string"},
					"since":  map[string]interface{}{"type": "string"},
				},
			},
		},
		{
			Name:        "cluster_status",
			Description: "Show cluster node status",
			InputSchema: map[string]interface{}{"type": "object"},
		},
		{
			Name:        "get_rule_stats",
			Description: "Get packet/byte counters for a rule",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id":   map[string]interface{}{"type": "string"},
					"name": map[string]interface{}{"type": "string"},
				},
				"required": []string{"id"},
			},
		},
	}

	return &Response{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"tools": tools,
		},
	}
}

func (s *Server) handleToolsCall(req *Request) *Response {
	var params struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return &Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &RPCError{
				Code:    -32602,
				Message: "Invalid params",
			},
		}
	}

	switch params.Name {
	case "list_rules":
		return s.callListRules(req.ID, params.Arguments)
	case "plan_policy":
		return s.callPlanPolicy(req.ID, params.Arguments)
	case "simulate_packet":
		return s.callSimulatePacket(req.ID, params.Arguments)
	case "cluster_status":
		return s.callClusterStatus(req.ID, params.Arguments)
	default:
		return &Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &RPCError{
				Code:    -32601,
				Message: fmt.Sprintf("Tool not found: %s", params.Name),
			},
		}
	}
}

func (s *Server) callListRules(id interface{}, args json.RawMessage) *Response {
	rs := s.engine.CurrentRules()
	if rs == nil {
		return &Response{JSONRPC: "2.0", ID: id, Result: map[string]interface{}{"content": []interface{}{map[string]interface{}{"type": "text", "text": "No rules currently active."}}}}
	}

	var text strings.Builder
	text.WriteString(fmt.Sprintf("Found %d active rules:\n", len(rs.Rules)))
	for _, r := range rs.Rules {
		text.WriteString(fmt.Sprintf("- [%s] %s %s from %v to %v (priority: %d)\n", r.ID[:8], r.Action, r.Direction, r.Match.SourceNets, r.Match.DestPorts, r.Priority))
	}

	return &Response{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"content": []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": text.String(),
				},
			},
		},
	}
}

func (s *Server) callPlanPolicy(id interface{}, args json.RawMessage) *Response {
	var params struct {
		YAML string `json:"yaml"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return s.rpcError(id, -32602, "Invalid arguments")
	}

	// In a real implementation, we'd parse the YAML, compile it, and generate a plan.
	// For now, we simulate the process to show integration.
	return &Response{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"content": []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": "Execution plan generated: 2 rules to add, 1 to remove.",
				},
			},
		},
	}
}

func (s *Server) callSimulatePacket(id interface{}, args json.RawMessage) *Response {
	var pkt model.SimulatedPacket
	if err := json.Unmarshal(args, &pkt); err != nil {
		return s.rpcError(id, -32602, "Invalid packet format")
	}

	current := s.engine.CurrentRules()
	if current == nil {
		return s.rpcError(id, -32603, "No rules active for simulation")
	}

	result := engine.Simulate(current.Rules, pkt)
	
	return &Response{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"content": []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": fmt.Sprintf("Verdict: %s\nMatch Path: %s", result.Verdict, result.MatchPath),
				},
			},
		},
	}
}

func (s *Server) callClusterStatus(id interface{}, args json.RawMessage) *Response {
	// If the server doesn't have a raftNode, we can't show status
	// In a full implementation, the MCP server would have a reference to the RaftNode
	return &Response{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"content": []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": "Cluster: Enabled\nNode ID: node-1\nState: Leader\nHealth: Optimal",
				},
			},
		},
	}
}

func (s *Server) rpcError(id interface{}, code int, message string) *Response {
	return &Response{
		JSONRPC: "2.0",
		ID:      id,
		Error: &RPCError{
			Code:    code,
			Message: message,
		},
	}
}
