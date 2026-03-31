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

// RegisterTools registers all tools from all categories with the MCP server.
// When readOnly is true, mutating tools (see IsMutatingTool) are not registered.
func RegisterTools(mcp *server.MCPServer, n *NATSServerTools, readOnly bool) {
	for _, category := range n.toolCategories() {
		for _, tool := range category.GetTools() {
			if readOnly && IsMutatingTool(tool.Tool.Name) {
				continue
			}
			tool.Register(mcp)
		}
	}
}
