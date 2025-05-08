package tools

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type Tool struct {
	Tool    mcp.Tool
	Handler server.ToolHandlerFunc
}

// List of available tools.
var Tools []Tool

func (t *Tool) Register(mcp *server.MCPServer) {
	mcp.AddTool(t.Tool, t.Handler)
}

func RegisterTools(mcp *server.MCPServer, tools []Tool) {
	for _, tool := range tools {
		tool.Register(mcp)
	}
}
