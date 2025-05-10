package tools

import (
	"context"
	"fmt"
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
						"account_name": map[string]interface{}{
							"type":        "string",
							"description": "The NATS account to use",
						},
						"expect": map[string]interface{}{
							"type":        "integer",
							"description": "How many servers to expect",
						},
					},
					Required: []string{"account_name"},
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
						"account_name": map[string]interface{}{
							"type":        "string",
							"description": "The NATS account to use",
						},
						"server": map[string]interface{}{
							"type":        "string",
							"description": "Server ID or Name to inspect",
						},
					},
					Required: []string{"account_name"},
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
						"account_name": map[string]interface{}{
							"type":        "string",
							"description": "The NATS account to use",
						},
						"expect": map[string]interface{}{
							"type":        "integer",
							"description": "How many servers to expect",
						},
					},
					Required: []string{"account_name"},
				},
			},
			Handler: s.serverPingHandler(),
		},
	}
}

// nats server list
// Args:
//
//	[<account_name>]  The NATS account to use
//	[<expect>]  How many servers to expect
func (s *ServerTools) serverListHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}

		executor, err := s.nats.GetExecutor(ctx, accountName)
		if err != nil {
			return nil, err
		}

		var args []string
		args = append(args, "server", "list")

		if expect, ok := request.Params.Arguments["expect"].(int); ok {
			args = append(args, strconv.Itoa(expect))
		}

		output, err := executor.ExecuteCommand(args...)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}

// nats server info
// Args:
//
//	[<account_name>]  The NATS account to use
//	[<server>]  Server ID or Name to inspect
func (s *ServerTools) serverInfoHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}
		executor, err := s.nats.GetExecutor(ctx, accountName)
		if err != nil {
			return nil, err
		}

		var args []string
		args = append(args, "server", "info")

		if server, ok := request.Params.Arguments["server"].(string); ok {
			args = append(args, server)
		}
		output, err := executor.ExecuteCommand(args...)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}

// nats server ping
// Args:
//
//	[<account_name>]  The NATS account to use
//	[<expect>]  How many servers to expect
func (s *ServerTools) serverPingHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}

		executor, err := s.nats.GetExecutor(ctx, accountName)
		if err != nil {
			return nil, err
		}

		var args []string
		args = append(args, "server", "ping")

		if expect, ok := request.Params.Arguments["expect"].(string); ok {
			args = append(args, expect)
		}
		output, err := executor.ExecuteCommand(args...)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}
