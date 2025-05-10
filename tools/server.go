package tools

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// GetServerTools returns all NATS server tools
func (n *NATSServerTools) GetServerTools() []Tool {
	return []Tool{
		{
			Tool: mcp.Tool{
				Name:        "server_list",
				Description: "Get NATS known server list",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"random_string": map[string]interface{}{
							"type":        "string",
							"description": "Dummy parameter for no-parameter tools",
						},
					},
					Required: []string{"random_string"},
				},
			},
			Handler: n.serverListHandler(),
		},
		{
			Tool: mcp.Tool{
				Name:        "server_check",
				Description: "Check NATS server health",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"random_string": map[string]interface{}{
							"type":        "string",
							"description": "Dummy parameter for no-parameter tools",
						},
					},
					Required: []string{"random_string"},
				},
			},
			Handler: n.serverCheckHandler(),
		},
		{
			Tool: mcp.Tool{
				Name:        "server_ping",
				Description: "Ping NATS server",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"random_string": map[string]interface{}{
							"type":        "string",
							"description": "Dummy parameter for no-parameter tools",
						},
					},
					Required: []string{"random_string"},
				},
			},
			Handler: n.serverPingHandler(),
		},
	}
}

func (n *NATSServerTools) executeNATSCommand(args ...string) (string, error) {
	baseArgs := []string{"-s", n.natsURL, "--creds", n.natsCredsPath}
	args = append(baseArgs, args...)

	cmd := exec.Command("nats", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("NATS command failed: %v, output: %s", err, string(output))
	}
	return string(output), nil
}

func (n *NATSServerTools) serverListHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		output, err := n.executeNATSCommand("server", "list")
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}

func (n *NATSServerTools) serverCheckHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		output, err := n.executeNATSCommand("server", "check")
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}

func (n *NATSServerTools) serverPingHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		output, err := n.executeNATSCommand("server", "ping")
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}
