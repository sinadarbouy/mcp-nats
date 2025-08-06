package tools

import (
	"context"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sinadarbouy/mcp-nats/tools/common"
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

// isAccountNameRequired determines if account_name is required based on auth strategy
func (s *ServerTools) isAccountNameRequired() bool {
	return common.IsAccountNameRequired()
}

// GetTools implements the ToolCategory interface
func (s *ServerTools) GetTools() []Tool {
	// Determine if we need account_name based on authentication strategy
	needsAccountName := s.isAccountNameRequired()

	return []Tool{
		{
			Tool: mcp.Tool{
				Name:        "server_list",
				Description: "Get NATS known server list",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: func() map[string]interface{} {
						props := map[string]interface{}{
							"expect": map[string]interface{}{
								"type":        "integer",
								"description": "How many servers to expect",
							},
						}
						if needsAccountName {
							props["account_name"] = map[string]interface{}{
								"type":        "string",
								"description": "The NATS account to use (required for credentials-based authentication)",
							}
						}
						return props
					}(),
					Required: func() []string {
						if needsAccountName {
							return []string{"account_name"}
						}
						return []string{}
					}(),
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
					Properties: func() map[string]interface{} {
						props := map[string]interface{}{
							"server": map[string]interface{}{
								"type":        "string",
								"description": "Server ID or Name to inspect",
							},
						}
						if needsAccountName {
							props["account_name"] = map[string]interface{}{
								"type":        "string",
								"description": "The NATS account to use (required for credentials-based authentication)",
							}
						}
						return props
					}(),
					Required: func() []string {
						if needsAccountName {
							return []string{"account_name"}
						}
						return []string{}
					}(),
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
					Properties: func() map[string]interface{} {
						props := map[string]interface{}{
							"expect": map[string]interface{}{
								"type":        "integer",
								"description": "How many servers to expect",
							},
						}
						if needsAccountName {
							props["account_name"] = map[string]interface{}{
								"type":        "string",
								"description": "The NATS account to use (required for credentials-based authentication)",
							}
						}
						return props
					}(),
					Required: func() []string {
						if needsAccountName {
							return []string{"account_name"}
						}
						return []string{}
					}(),
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
		accountName, err := common.DetermineAccountName(request.Params.Arguments)
		if err != nil {
			return nil, err
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
		accountName, err := common.DetermineAccountName(request.Params.Arguments)
		if err != nil {
			return nil, err
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
		accountName, err := common.DetermineAccountName(request.Params.Arguments)
		if err != nil {
			return nil, err
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
