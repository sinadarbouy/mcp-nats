package tools

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"text/template"
	"time"

	"github.com/google/uuid"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sinadarbouy/mcp-nats/tools/common"
)

// PublishTools represents all NATS publish-related tools
type PublishTools struct {
	nats *NATSServerTools
}

// NewPublishTools creates a new PublishTools instance
func NewPublishTools(nats *NATSServerTools) *PublishTools {
	return &PublishTools{
		nats: nats,
	}
}

// isAccountNameRequired determines if account_name is required based on auth strategy
func (p *PublishTools) isAccountNameRequired() bool {
	return common.IsAccountNameRequired()
}

// GetTools implements the ToolCategory interface
func (p *PublishTools) GetTools() []Tool {
	// Determine if we need account_name based on authentication strategy
	needsAccountName := p.isAccountNameRequired()

	return []Tool{
		{
			Tool: mcp.Tool{
				Name:        "publish",
				Description: "Generic data publish utility",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: func() map[string]interface{} {
						props := map[string]interface{}{
							"subject": map[string]interface{}{
								"type":        "string",
								"description": "Subject to publish to",
							},
							"body": map[string]interface{}{
								"type":        "string",
								"description": "Message body",
							},
							"reply": map[string]interface{}{
								"type":        "string",
								"description": "Sets a custom reply to subject",
							},
							"header": map[string]interface{}{
								"type":        "array",
								"items":       map[string]interface{}{"type": "string"},
								"description": "Adds headers to the message",
							},
							"count": map[string]interface{}{
								"type":        "integer",
								"description": "Publish multiple messages",
								"default":     1,
							},
							"sleep": map[string]interface{}{
								"type":        "string",
								"description": "When publishing multiple messages, sleep between publishes",
							},
							"force_stdin": map[string]interface{}{
								"type":        "boolean",
								"description": "Force reading from stdin",
								"default":     false,
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
							return []string{"account_name", "subject"}
						}
						return []string{"subject"}
					}(),
				},
			},
			Handler: p.publishHandler(),
		},
	}
}

// Template functions for message generation
type templateData struct {
	Count     int
	TimeStamp string
	Unix      int64
	UnixNano  int64
	Time      time.Time
	ID        string
}

func generateRandomString(min, max int) string {
	length := rand.Intn(max-min+1) + min
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func (p *PublishTools) publishHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, err := common.DetermineAccountName(request.Params.Arguments)
		if err != nil {
			return nil, err
		}

		subject, ok := request.Params.Arguments["subject"].(string)
		if !ok {
			return nil, fmt.Errorf("missing subject")
		}

		executor, err := p.nats.GetExecutor(ctx, accountName)
		if err != nil {
			return nil, fmt.Errorf("failed to get NATS executor: %w", err)
		}

		body := ""
		if b, ok := request.Params.Arguments["body"].(string); ok {
			body = b
		}

		count := 1
		if c, ok := request.Params.Arguments["count"].(float64); ok {
			count = int(c)
		}

		var sleep time.Duration
		if s, ok := request.Params.Arguments["sleep"].(string); ok {
			sleep, err = time.ParseDuration(s)
			if err != nil {
				return nil, fmt.Errorf("invalid sleep duration: %w", err)
			}
		}

		// Create template functions map
		funcMap := template.FuncMap{
			"Random": generateRandomString,
		}

		// Create template
		tmpl, err := template.New("message").Funcs(funcMap).Parse(body)
		if err != nil {
			return nil, fmt.Errorf("failed to parse message template: %w", err)
		}

		for i := 0; i < count; i++ {
			now := time.Now()
			data := templateData{
				Count:     i + 1,
				TimeStamp: now.Format(time.RFC3339),
				Unix:      now.Unix(),
				UnixNano:  now.UnixNano(),
				Time:      now,
				ID:        uuid.New().String(),
			}

			var msgBuilder strings.Builder
			if err := tmpl.Execute(&msgBuilder, data); err != nil {
				return nil, fmt.Errorf("failed to execute message template: %w", err)
			}

			msg := msgBuilder.String()

			// Build command arguments
			args := []string{"pub"}

			// Add reply if specified
			if r, ok := request.Params.Arguments["reply"].(string); ok && r != "" {
				args = append(args, "--reply", r)
			}

			// Add headers if specified
			if h, ok := request.Params.Arguments["header"].([]interface{}); ok {
				for _, header := range h {
					args = append(args, "--header", header.(string))
				}
			}

			// Add subject and message
			args = append(args, subject, msg)

			// Execute the command
			if _, err := executor.ExecuteCommand(args...); err != nil {
				return nil, fmt.Errorf("failed to publish message: %w", err)
			}

			if i < count-1 && sleep > 0 {
				time.Sleep(sleep)
			}
		}

		return mcp.NewToolResultText(fmt.Sprintf("Published %d message(s) to %s", count, subject)), nil
	}
}
