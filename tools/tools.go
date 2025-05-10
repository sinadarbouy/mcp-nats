package tools

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// ToolCategory represents a group of related NATS tools
type ToolCategory interface {
	GetTools() []Tool
}

// Tool combines an MCP tool definition with its handler
type Tool struct {
	Tool    mcp.Tool
	Handler server.ToolHandlerFunc
}

// Register registers a single tool with the MCP server
func (t *Tool) Register(mcp *server.MCPServer) {
	mcp.AddTool(t.Tool, t.Handler)
}

// RegisterTools registers all tools from all categories with the MCP server
func RegisterTools(mcp *server.MCPServer, n *NATSServerTools) {
	// Define tool categories
	categories := []ToolCategory{
		n.ServerTools(),
		n.StreamTools(),
		// Add new categories here as needed
	}

	// Register all tools from each category
	for _, category := range categories {
		for _, tool := range category.GetTools() {
			tool.Register(mcp)
		}
	}
}
