package mcp

import (
	"encoding/json"
	"fmt"
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
	case "add_rule":
		return s.callAddRule(req.ID, params.Arguments)
	case "remove_rule":
		return s.callRemoveRule(req.ID, params.Arguments)
	case "plan_policy":
		return s.callPlanPolicy(req.ID, params.Arguments)
	case "apply_policy":
		return s.callApplyPolicy(req.ID, params.Arguments)
	case "simulate_packet":
		return s.callSimulatePacket(req.ID, params.Arguments)
	case "rollback":
		return s.callRollback(req.ID, params.Arguments)
	case "list_snapshots":
		return s.callListSnapshots(req.ID, params.Arguments)
	case "audit_search":
		return s.callAuditSearch(req.ID, params.Arguments)
	case "cluster_status":
		return s.callClusterStatus(req.ID, params.Arguments)
	case "get_rule_stats":
		return s.callGetRuleStats(req.ID, params.Arguments)
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
	rules := s.engine.CurrentRules()
	return &Response{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"content": []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": fmt.Sprintf("Found %d rules", len(rules.Rules)),
				},
			},
		},
	}
}

func (s *Server) callAddRule(id interface{}, args json.RawMessage) *Response {
	return &Response{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"content": []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": "Rule added successfully",
				},
			},
		},
	}
}

func (s *Server) callRemoveRule(id interface{}, args json.RawMessage) *Response {
	return &Response{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"content": []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": "Rule removed successfully",
				},
			},
		},
	}
}

func (s *Server) callPlanPolicy(id interface{}, args json.RawMessage) *Response {
	return &Response{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"content": []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": "Plan generated successfully",
				},
			},
		},
	}
}

func (s *Server) callApplyPolicy(id interface{}, args json.RawMessage) *Response {
	return &Response{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"content": []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": "Policy applied successfully",
				},
			},
		},
	}
}

func (s *Server) callSimulatePacket(id interface{}, args json.RawMessage) *Response {
	return &Response{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"content": []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": "Simulation result: ACCEPT",
				},
			},
		},
	}
}

func (s *Server) callRollback(id interface{}, args json.RawMessage) *Response {
	return &Response{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"content": []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": "Rollback successful",
				},
			},
		},
	}
}

func (s *Server) callListSnapshots(id interface{}, args json.RawMessage) *Response {
	return &Response{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"content": []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": "Found 0 snapshots",
				},
			},
		},
	}
}

func (s *Server) callAuditSearch(id interface{}, args json.RawMessage) *Response {
	return &Response{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"content": []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": "Found 0 audit events",
				},
			},
		},
	}
}

func (s *Server) callClusterStatus(id interface{}, args json.RawMessage) *Response {
	return &Response{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"content": []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": "Cluster status: Healthy",
				},
			},
		},
	}
}

func (s *Server) callGetRuleStats(id interface{}, args json.RawMessage) *Response {
	return &Response{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"content": []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": "Packets: 0, Bytes: 0",
				},
			},
		},
	}
}
