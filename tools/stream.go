package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func (n *NATSServerTools) streamInfoHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		stream, ok := request.Params.Arguments["stream"].(string)
		if !ok {
			return nil, fmt.Errorf("missing stream")
		}
		output, err := n.executeNATSCommand("stream", "info", stream)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}

func (n *NATSServerTools) GetStreamTools() []Tool {
	return []Tool{
		{
			Tool: mcp.Tool{
				Name:        "stream_info",
				Description: "Get information about a NATS stream",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"stream": map[string]interface{}{
							"type":        "string",
							"description": "The name of the stream to get information about",
						},
					},
					Required: []string{"stream"},
				},
			},
			Handler: n.streamInfoHandler(),
		},
	}
}
