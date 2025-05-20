package tools

import (
	"context"
	"fmt"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RTTTools represents all NATS RTT-related tools
type RTTTools struct {
	nats *NATSServerTools
}

// NewRTTTools creates a new RTTTools instance
func NewRTTTools(nats *NATSServerTools) *RTTTools {
	return &RTTTools{
		nats: nats,
	}
}

// GetTools implements the ToolCategory interface
func (r *RTTTools) GetTools() []Tool {
	return []Tool{
		{
			Tool: mcp.Tool{
				Name:        "rtt",
				Description: "Compute round-trip time to NATS server",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"account_name": map[string]interface{}{
							"type":        "string",
							"description": "The NATS account to use",
						},
						"iterations": map[string]interface{}{
							"type":        "integer",
							"description": "How many round trips to do when testing",
						},
						"json": map[string]interface{}{
							"type":        "boolean",
							"description": "Produce JSON output",
							"default":     false,
						},
					},
					Required: []string{"account_name"},
				},
			},
			Handler: r.rttHandler(),
		},
	}
}

func (r *RTTTools) rttHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}

		executor, err := r.nats.GetExecutor(ctx, accountName)
		if err != nil {
			return nil, err
		}

		args := []string{"rtt"}

		// Add json flag if specified
		if json, ok := request.Params.Arguments["json"].(bool); ok && json {
			args = append(args, "--json")
		}

		// Add iterations if specified
		if iterations, ok := request.Params.Arguments["iterations"].(int); ok {
			args = append(args, strconv.Itoa(iterations))
		}

		output, err := executor.ExecuteCommand(args...)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}
