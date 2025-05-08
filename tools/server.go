package tools

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func GetServerlistToolHandler(natsURL string, natsCredsPath string) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		cmd := exec.Command("nats", "-s", natsURL, "--creds", natsCredsPath, "server", "list")

		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		return mcp.NewToolResultText(string(output)), nil
	}
}
