package tools

import (
	"context"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// ServerTools represents all NATS server-related tools
type ServerTools struct {
	nats *NATSServerTools
}

// NewServerTools creates a new ServerTools instance
func NewServerTools(nats *NATSServerTools) *ServerTools {
	return &ServerTools{
		nats: nats,
	}
}

// GetTools implements the ToolCategory interface
func (s *ServerTools) GetTools() []Tool {
	return []Tool{
		{
			Tool: mcp.Tool{
				Name:        "server_list",
				Description: "Get NATS known server list",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"expect": map[string]interface{}{
							"type":        "integer",
							"description": "How many servers to expect",
						},
					},
					Required: []string{},
				},
			},
			Handler: s.serverListHandler(),
		},
		{
			Tool: mcp.Tool{
				Name:        "server_info",
				Description: "Get NATS server info",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"server": map[string]interface{}{
							"type":        "string",
							"description": "Server ID or Name to inspect",
						},
					},
					Required: []string{},
				},
			},
			Handler: s.serverInfoHandler(),
		},
		{
			Tool: mcp.Tool{
				Name:        "server_ping",
				Description: "Ping NATS server",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"expect": map[string]interface{}{
							"type":        "integer",
							"description": "How many servers to expect",
						},
					},
					Required: []string{},
				},
			},
			Handler: s.serverPingHandler(),
		},
	}
}

// nats server list
// Args:
//
//	[<expect>]  How many servers to expect
func (s *ServerTools) serverListHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var args []string
		args = append(args, "server", "list")

		if expect, ok := request.Params.Arguments["expect"].(int); ok {
			args = append(args, strconv.Itoa(expect))
		}

		output, err := s.nats.GetSysExecutor().ExecuteCommand(args...)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}

// nats server info
// Args:
//
//	[<server>]  Server ID or Name to inspect
func (s *ServerTools) serverInfoHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var args []string
		args = append(args, "server", "info")

		if server, ok := request.Params.Arguments["server"].(string); ok {
			args = append(args, server)
		}
		output, err := s.nats.GetSysExecutor().ExecuteCommand(args...)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}

// nats server ping
// Args:
//
//	[<expect>]  How many servers to expect
func (s *ServerTools) serverPingHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var args []string
		args = append(args, "server", "ping")

		if expect, ok := request.Params.Arguments["expect"].(string); ok {
			args = append(args, expect)
		}
		output, err := s.nats.GetSysExecutor().ExecuteCommand(args...)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}
