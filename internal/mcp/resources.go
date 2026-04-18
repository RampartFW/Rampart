package mcp

import (
	"encoding/json"
	"fmt"
)

type Resource struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	MimeType    string `json:"mimeType,omitempty"`
}

func (s *Server) handleResourcesList(req *Request) *Response {
	resources := []Resource{
		{
			URI:         "rampart://policies/current",
			Name:        "Current Policy",
			Description: "Current active policy as YAML",
			MimeType:    "application/x-yaml",
		},
		{
			URI:         "rampart://rules",
			Name:        "Active Rules",
			Description: "Active rules in JSON format",
			MimeType:    "application/json",
		},
		{
			URI:         "rampart://audit/recent",
			Name:        "Recent Audit Events",
			Description: "Recent audit events",
			MimeType:    "application/json",
		},
		{
			URI:         "rampart://cluster/status",
			Name:        "Cluster Status",
			Description: "Cluster health and node list",
			MimeType:    "application/json",
		},
	}

	return &Response{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"resources": resources,
		},
	}
}

func (s *Server) handleResourcesRead(req *Request) *Response {
	var params struct {
		URI string `json:"uri"`
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

	switch params.URI {
	case "rampart://policies/current":
		return s.readCurrentPolicy(req.ID)
	case "rampart://rules":
		return s.readRules(req.ID)
	case "rampart://audit/recent":
		return s.readAuditRecent(req.ID)
	case "rampart://cluster/status":
		return s.readClusterStatus(req.ID)
	default:
		return &Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &RPCError{
				Code:    -32601,
				Message: fmt.Sprintf("Resource not found: %s", params.URI),
			},
		}
	}
}

func (s *Server) readCurrentPolicy(id interface{}) *Response {
	return &Response{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"contents": []interface{}{
				map[string]interface{}{
					"uri":      "rampart://policies/current",
					"mimeType": "application/x-yaml",
					"text":     "apiVersion: rampartfw.com/v1\nkind: PolicySet\nmetadata:\n  name: current\npolicies: []",
				},
			},
		},
	}
}

func (s *Server) readRules(id interface{}) *Response {
	rules := s.engine.CurrentRules()
	data, _ := json.MarshalIndent(rules, "", "  ")
	return &Response{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"contents": []interface{}{
				map[string]interface{}{
					"uri":      "rampart://rules",
					"mimeType": "application/json",
					"text":     string(data),
				},
			},
		},
	}
}

func (s *Server) readAuditRecent(id interface{}) *Response {
	return &Response{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"contents": []interface{}{
				map[string]interface{}{
					"uri":      "rampart://audit/recent",
					"mimeType": "application/json",
					"text":     "[]",
				},
			},
		},
	}
}

func (s *Server) readClusterStatus(id interface{}) *Response {
	return &Response{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"contents": []interface{}{
				map[string]interface{}{
					"uri":      "rampart://cluster/status",
					"mimeType": "application/json",
					"text":     "{\"status\": \"Healthy\", \"nodes\": []}",
				},
			},
		},
	}
}

