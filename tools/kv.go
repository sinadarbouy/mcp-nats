package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// KVTools represents all NATS KV-related tools
type KVTools struct {
	nats *NATSServerTools
}

// NewKVTools creates a new KVTools instance
func NewKVTools(nats *NATSServerTools) *KVTools {
	return &KVTools{
		nats: nats,
	}
}

// GetTools implements the ToolCategory interface
func (k *KVTools) GetTools() []Tool {
	return []Tool{
		{
			Tool: mcp.Tool{
				Name:        "kv_add",
				Description: "Adds a new KV Store Bucket",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"account_name": map[string]interface{}{
							"type":        "string",
							"description": "The NATS account to use",
						},
						"bucket": map[string]interface{}{
							"type":        "string",
							"description": "The name of the bucket to create",
						},
						"history": map[string]interface{}{
							"type":        "integer",
							"description": "How many historic values to keep per key",
							"default":     1,
						},
						"ttl": map[string]interface{}{
							"type":        "string",
							"description": "How long to keep values for",
						},
						"replicas": map[string]interface{}{
							"type":        "integer",
							"description": "How many replicas of the data to store",
							"default":     1,
						},
						"max_value_size": map[string]interface{}{
							"type":        "integer",
							"description": "Maximum size for any single value in bytes",
						},
						"max_bucket_size": map[string]interface{}{
							"type":        "integer",
							"description": "Maximum size for the bucket in bytes",
						},
						"description": map[string]interface{}{
							"type":        "string",
							"description": "A description for the bucket",
						},
						"storage": map[string]interface{}{
							"type":        "string",
							"description": "Storage backend to use (file, memory)",
							"enum":        []string{"file", "memory"},
						},
						"compress": map[string]interface{}{
							"type":        "boolean",
							"description": "Whether to compress the bucket data",
						},
						"tags": map[string]interface{}{
							"type":        "array",
							"items":       map[string]interface{}{"type": "string"},
							"description": "Place the bucket on servers that has specific tags",
						},
						"cluster": map[string]interface{}{
							"type":        "string",
							"description": "Place the bucket on a specific cluster",
						},
						"republish_source": map[string]interface{}{
							"type":        "string",
							"description": "Republish messages to republish_destination",
						},
						"republish_destination": map[string]interface{}{
							"type":        "string",
							"description": "Republish destination for messages in republish_source",
						},
						"republish_headers": map[string]interface{}{
							"type":        "boolean",
							"description": "Republish only message headers, no bodies",
						},
						"mirror": map[string]interface{}{
							"type":        "string",
							"description": "Creates a mirror of a different bucket",
						},
						"mirror_domain": map[string]interface{}{
							"type":        "string",
							"description": "When mirroring find the bucket in a different domain",
						},
						"source": map[string]interface{}{
							"type":        "array",
							"items":       map[string]interface{}{"type": "string"},
							"description": "Source from a different bucket",
						},
					},
					Required: []string{"account_name", "bucket"},
				},
			},
			Handler: k.kvAddHandler(),
		},
		{
			Tool: mcp.Tool{
				Name:        "kv_put",
				Description: "Puts a value into a key",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"account_name": map[string]interface{}{
							"type":        "string",
							"description": "The NATS account to use",
						},
						"bucket": map[string]interface{}{
							"type":        "string",
							"description": "The bucket name",
						},
						"key": map[string]interface{}{
							"type":        "string",
							"description": "The key to set",
						},
						"value": map[string]interface{}{
							"type":        "string",
							"description": "The value to store, when empty reads STDIN",
						},
						"stdin": map[string]interface{}{
							"type":        "string",
							"description": "Value to use as STDIN when no value is provided",
						},
						"flags": map[string]interface{}{
							"type":        "array",
							"items":       map[string]interface{}{"type": "string"},
							"description": "Optional flags to pass to the command",
						},
					},
					Required: []string{"account_name", "bucket", "key"},
				},
			},
			Handler: k.kvPutHandler(),
		},
		{
			Tool: mcp.Tool{
				Name:        "kv_get",
				Description: "Gets a value for a key",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"account_name": map[string]interface{}{
							"type":        "string",
							"description": "The NATS account to use",
						},
						"bucket": map[string]interface{}{
							"type":        "string",
							"description": "The bucket name",
						},
						"key": map[string]interface{}{
							"type":        "string",
							"description": "The key to get",
						},
						"revision": map[string]interface{}{
							"type":        "string",
							"description": "Gets a specific revision",
						},
						"raw": map[string]interface{}{
							"type":        "boolean",
							"description": "Show only the value string",
						},
						"flags": map[string]interface{}{
							"type":        "array",
							"items":       map[string]interface{}{"type": "string"},
							"description": "Optional flags to pass to the command",
						},
					},
					Required: []string{"account_name", "bucket", "key"},
				},
			},
			Handler: k.kvGetHandler(),
		},
		{
			Tool: mcp.Tool{
				Name:        "kv_create",
				Description: "Puts a value into a key only if the key is new or its last operation was a delete",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"account_name": map[string]interface{}{
							"type":        "string",
							"description": "The NATS account to use",
						},
						"bucket": map[string]interface{}{
							"type":        "string",
							"description": "The bucket name",
						},
						"key": map[string]interface{}{
							"type":        "string",
							"description": "The key to create",
						},
						"value": map[string]interface{}{
							"type":        "string",
							"description": "The value to store, when empty reads STDIN",
						},
						"stdin": map[string]interface{}{
							"type":        "string",
							"description": "Value to use as STDIN when no value is provided",
						},
						"flags": map[string]interface{}{
							"type":        "array",
							"items":       map[string]interface{}{"type": "string"},
							"description": "Optional flags to pass to the command",
						},
					},
					Required: []string{"account_name", "bucket", "key"},
				},
			},
			Handler: k.kvCreateHandler(),
		},
		{
			Tool: mcp.Tool{
				Name:        "kv_update",
				Description: "Updates a key with a new value if the previous value matches the given revision",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"account_name": map[string]interface{}{
							"type":        "string",
							"description": "The NATS account to use",
						},
						"bucket": map[string]interface{}{
							"type":        "string",
							"description": "The bucket name",
						},
						"key": map[string]interface{}{
							"type":        "string",
							"description": "The key to update",
						},
						"value": map[string]interface{}{
							"type":        "string",
							"description": "The value to store, when empty reads STDIN",
						},
						"stdin": map[string]interface{}{
							"type":        "string",
							"description": "Value to use as STDIN when no value is provided",
						},
						"revision": map[string]interface{}{
							"type":        "string",
							"description": "The revision of the previous value in the bucket",
						},
						"flags": map[string]interface{}{
							"type":        "array",
							"items":       map[string]interface{}{"type": "string"},
							"description": "Optional flags to pass to the command",
						},
					},
					Required: []string{"account_name", "bucket", "key"},
				},
			},
			Handler: k.kvUpdateHandler(),
		},
		{
			Tool: mcp.Tool{
				Name:        "kv_del",
				Description: "Deletes a key or the entire bucket",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"account_name": map[string]interface{}{
							"type":        "string",
							"description": "The NATS account to use",
						},
						"bucket": map[string]interface{}{
							"type":        "string",
							"description": "The bucket to act on",
						},
						"key": map[string]interface{}{
							"type":        "string",
							"description": "The key to act on",
						},
						"force": map[string]interface{}{
							"type":        "boolean",
							"description": "Act without confirmation",
							"default":     false,
						},
						"flags": map[string]interface{}{
							"type":        "array",
							"items":       map[string]interface{}{"type": "string"},
							"description": "Optional flags to pass to the command",
						},
					},
					Required: []string{"account_name", "bucket"},
				},
			},
			Handler: k.kvDelHandler(),
		},
		{
			Tool: mcp.Tool{
				Name:        "kv_purge",
				Description: "Deletes a key from the bucket, clearing history before creating a delete marker",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"account_name": map[string]interface{}{
							"type":        "string",
							"description": "The NATS account to use",
						},
						"bucket": map[string]interface{}{
							"type":        "string",
							"description": "The bucket to act on",
						},
						"key": map[string]interface{}{
							"type":        "string",
							"description": "The key to act on",
						},
						"force": map[string]interface{}{
							"type":        "boolean",
							"description": "Act without confirmation",
							"default":     false,
						},
						"flags": map[string]interface{}{
							"type":        "array",
							"items":       map[string]interface{}{"type": "string"},
							"description": "Optional flags to pass to the command",
						},
					},
					Required: []string{"account_name", "bucket", "key"},
				},
			},
			Handler: k.kvPurgeHandler(),
		},
		{
			Tool: mcp.Tool{
				Name:        "kv_history",
				Description: "Shows the full history for a key",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"account_name": map[string]interface{}{
							"type":        "string",
							"description": "The NATS account to use",
						},
						"bucket": map[string]interface{}{
							"type":        "string",
							"description": "The bucket to act on",
						},
						"key": map[string]interface{}{
							"type":        "string",
							"description": "The key to act on",
						},
						"flags": map[string]interface{}{
							"type":        "array",
							"items":       map[string]interface{}{"type": "string"},
							"description": "Optional flags to pass to the command",
						},
					},
					Required: []string{"account_name", "bucket", "key"},
				},
			},
			Handler: k.kvHistoryHandler(),
		},
		{
			Tool: mcp.Tool{
				Name:        "kv_ls",
				Description: "List available buckets or the keys in a bucket",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"account_name": map[string]interface{}{
							"type":        "string",
							"description": "The NATS account to use",
						},
						"bucket": map[string]interface{}{
							"type":        "string",
							"description": "The bucket to list the keys",
						},
						"names": map[string]interface{}{
							"type":        "boolean",
							"description": "Show just the bucket names",
							"default":     false,
						},
						"verbose": map[string]interface{}{
							"type":        "boolean",
							"description": "Show detailed info about the key",
							"default":     false,
						},
						"display_value": map[string]interface{}{
							"type":        "boolean",
							"description": "Display value in verbose output (has no effect without 'verbose')",
							"default":     false,
						},
						"flags": map[string]interface{}{
							"type":        "array",
							"items":       map[string]interface{}{"type": "string"},
							"description": "Optional flags to pass to the command",
						},
					},
					Required: []string{"account_name"},
				},
			},
			Handler: k.kvLsHandler(),
		},
		{
			Tool: mcp.Tool{
				Name:        "kv_watch",
				Description: "Watch the bucket or a specific key for updated",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"account_name": map[string]interface{}{
							"type":        "string",
							"description": "The NATS account to use",
						},
						"bucket": map[string]interface{}{
							"type":        "string",
							"description": "The bucket to act on",
						},
						"key": map[string]interface{}{
							"type":        "string",
							"description": "The key to act on",
						},
						"flags": map[string]interface{}{
							"type":        "array",
							"items":       map[string]interface{}{"type": "string"},
							"description": "Optional flags to pass to the command",
						},
					},
					Required: []string{"account_name", "bucket"},
				},
			},
			Handler: k.kvWatchHandler(),
		},
		{
			Tool: mcp.Tool{
				Name:        "kv_info",
				Description: "View the status of a KV store",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"account_name": map[string]interface{}{
							"type":        "string",
							"description": "The NATS account to use",
						},
						"bucket": map[string]interface{}{
							"type":        "string",
							"description": "The bucket to act on",
						},
						"flags": map[string]interface{}{
							"type":        "array",
							"items":       map[string]interface{}{"type": "string"},
							"description": "Optional flags to pass to the command",
						},
					},
					Required: []string{"account_name"},
				},
			},
			Handler: k.kvInfoHandler(),
		},
		{
			Tool: mcp.Tool{
				Name:        "kv_compact",
				Description: "Reclaim space used by deleted keys",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"account_name": map[string]interface{}{
							"type":        "string",
							"description": "The NATS account to use",
						},
						"bucket": map[string]interface{}{
							"type":        "string",
							"description": "The bucket to act on",
						},
						"force": map[string]interface{}{
							"type":        "boolean",
							"description": "Act without confirmation",
							"default":     false,
						},
						"flags": map[string]interface{}{
							"type":        "array",
							"items":       map[string]interface{}{"type": "string"},
							"description": "Optional flags to pass to the command",
						},
					},
					Required: []string{"account_name", "bucket"},
				},
			},
			Handler: k.kvCompactHandler(),
		},
	}
}

// Handler implementations
func (k *KVTools) kvAddHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}

		bucket, ok := request.Params.Arguments["bucket"].(string)
		if !ok {
			return nil, fmt.Errorf("missing bucket")
		}

		executor, err := k.nats.GetExecutor(ctx, accountName)
		if err != nil {
			return nil, err
		}

		args := []string{"kv", "add", bucket}

		// Add optional flags
		if history, ok := request.Params.Arguments["history"].(float64); ok {
			args = append(args, fmt.Sprintf("--history=%d", int(history)))
		}
		if ttl, ok := request.Params.Arguments["ttl"].(string); ok {
			args = append(args, fmt.Sprintf("--ttl=%s", ttl))
		}
		if replicas, ok := request.Params.Arguments["replicas"].(float64); ok {
			args = append(args, fmt.Sprintf("--replicas=%d", int(replicas)))
		}
		if maxValueSize, ok := request.Params.Arguments["max_value_size"].(float64); ok {
			args = append(args, fmt.Sprintf("--max-value-size=%d", int(maxValueSize)))
		}
		if maxBucketSize, ok := request.Params.Arguments["max_bucket_size"].(float64); ok {
			args = append(args, fmt.Sprintf("--max-bucket-size=%d", int(maxBucketSize)))
		}
		if description, ok := request.Params.Arguments["description"].(string); ok {
			args = append(args, fmt.Sprintf("--description=%s", description))
		}
		if storage, ok := request.Params.Arguments["storage"].(string); ok {
			args = append(args, fmt.Sprintf("--storage=%s", storage))
		}
		if compress, ok := request.Params.Arguments["compress"].(bool); ok {
			if compress {
				args = append(args, "--compress")
			} else {
				args = append(args, "--no-compress")
			}
		}
		if tags, ok := request.Params.Arguments["tags"].([]interface{}); ok {
			for _, tag := range tags {
				if strTag, ok := tag.(string); ok {
					args = append(args, fmt.Sprintf("--tags=%s", strTag))
				}
			}
		}
		if cluster, ok := request.Params.Arguments["cluster"].(string); ok {
			args = append(args, fmt.Sprintf("--cluster=%s", cluster))
		}
		if republishSource, ok := request.Params.Arguments["republish_source"].(string); ok {
			args = append(args, fmt.Sprintf("--republish-source=%s", republishSource))
		}
		if republishDest, ok := request.Params.Arguments["republish_destination"].(string); ok {
			args = append(args, fmt.Sprintf("--republish-destination=%s", republishDest))
		}
		if republishHeaders, ok := request.Params.Arguments["republish_headers"].(bool); ok && republishHeaders {
			args = append(args, "--republish-headers")
		}
		if mirror, ok := request.Params.Arguments["mirror"].(string); ok {
			args = append(args, fmt.Sprintf("--mirror=%s", mirror))
		}
		if mirrorDomain, ok := request.Params.Arguments["mirror_domain"].(string); ok {
			args = append(args, fmt.Sprintf("--mirror-domain=%s", mirrorDomain))
		}
		if sources, ok := request.Params.Arguments["source"].([]interface{}); ok {
			for _, source := range sources {
				if strSource, ok := source.(string); ok {
					args = append(args, fmt.Sprintf("--source=%s", strSource))
				}
			}
		}

		// Add any additional flags passed
		if flags := getFlags(request.Params.Arguments); flags != nil {
			args = append(args, flags...)
		}

		output, err := executor.ExecuteCommand(args...)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}

func (k *KVTools) kvPutHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}

		bucket, ok := request.Params.Arguments["bucket"].(string)
		if !ok {
			return nil, fmt.Errorf("missing bucket")
		}

		key, ok := request.Params.Arguments["key"].(string)
		if !ok {
			return nil, fmt.Errorf("missing key")
		}

		executor, err := k.nats.GetExecutor(ctx, accountName)
		if err != nil {
			return nil, err
		}

		args := []string{"kv", "put", bucket, key}

		// If value is provided, add it to args
		if value, ok := request.Params.Arguments["value"].(string); ok {
			args = append(args, value)
		} else if stdin, ok := request.Params.Arguments["stdin"].(string); ok {
			// If no value but stdin is provided, use it as input
			executor.SetStdin(stdin)
		}

		// Add any additional flags passed
		if flags := getFlags(request.Params.Arguments); flags != nil {
			args = append(args, flags...)
		}

		output, err := executor.ExecuteCommand(args...)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}

func (k *KVTools) kvGetHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}

		bucket, ok := request.Params.Arguments["bucket"].(string)
		if !ok {
			return nil, fmt.Errorf("missing bucket")
		}

		key, ok := request.Params.Arguments["key"].(string)
		if !ok {
			return nil, fmt.Errorf("missing key")
		}

		executor, err := k.nats.GetExecutor(ctx, accountName)
		if err != nil {
			return nil, err
		}

		args := []string{"kv", "get", bucket, key}

		// Add revision flag if provided
		if revision, ok := request.Params.Arguments["revision"].(string); ok {
			args = append(args, fmt.Sprintf("--revision=%s", revision))
		}

		// Add raw flag if true
		if raw, ok := request.Params.Arguments["raw"].(bool); ok && raw {
			args = append(args, "--raw")
		}

		// Add any additional flags passed
		if flags := getFlags(request.Params.Arguments); flags != nil {
			args = append(args, flags...)
		}

		output, err := executor.ExecuteCommand(args...)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}

func (k *KVTools) kvCreateHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}

		bucket, ok := request.Params.Arguments["bucket"].(string)
		if !ok {
			return nil, fmt.Errorf("missing bucket")
		}

		key, ok := request.Params.Arguments["key"].(string)
		if !ok {
			return nil, fmt.Errorf("missing key")
		}

		executor, err := k.nats.GetExecutor(ctx, accountName)
		if err != nil {
			return nil, err
		}

		args := []string{"kv", "create", bucket, key}

		// If value is provided, add it to args
		if value, ok := request.Params.Arguments["value"].(string); ok {
			args = append(args, value)
		} else if stdin, ok := request.Params.Arguments["stdin"].(string); ok {
			// If no value but stdin is provided, use it as input
			executor.SetStdin(stdin)
		}

		// Add any additional flags passed
		if flags := getFlags(request.Params.Arguments); flags != nil {
			args = append(args, flags...)
		}

		output, err := executor.ExecuteCommand(args...)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}

func (k *KVTools) kvUpdateHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}

		bucket, ok := request.Params.Arguments["bucket"].(string)
		if !ok {
			return nil, fmt.Errorf("missing bucket")
		}

		key, ok := request.Params.Arguments["key"].(string)
		if !ok {
			return nil, fmt.Errorf("missing key")
		}

		executor, err := k.nats.GetExecutor(ctx, accountName)
		if err != nil {
			return nil, err
		}

		args := []string{"kv", "update", bucket, key}

		// If value is provided, add it to args
		if value, ok := request.Params.Arguments["value"].(string); ok {
			args = append(args, value)
		} else if stdin, ok := request.Params.Arguments["stdin"].(string); ok {
			// If no value but stdin is provided, use it as input
			executor.SetStdin(stdin)
		}

		// Add revision if provided
		if revision, ok := request.Params.Arguments["revision"].(string); ok {
			args = append(args, revision)
		}

		// Add any additional flags passed
		if flags := getFlags(request.Params.Arguments); flags != nil {
			args = append(args, flags...)
		}

		output, err := executor.ExecuteCommand(args...)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}

func (k *KVTools) kvDelHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}

		bucket, ok := request.Params.Arguments["bucket"].(string)
		if !ok {
			return nil, fmt.Errorf("missing bucket")
		}

		executor, err := k.nats.GetExecutor(ctx, accountName)
		if err != nil {
			return nil, err
		}

		args := []string{"kv", "del", bucket}

		// Add key if provided
		if key, ok := request.Params.Arguments["key"].(string); ok {
			args = append(args, key)
		}

		// Add force flag if true
		if force, ok := request.Params.Arguments["force"].(bool); ok && force {
			args = append(args, "--force")
		}

		// Add any additional flags passed
		if flags := getFlags(request.Params.Arguments); flags != nil {
			args = append(args, flags...)
		}

		output, err := executor.ExecuteCommand(args...)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}

func (k *KVTools) kvPurgeHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}

		bucket, ok := request.Params.Arguments["bucket"].(string)
		if !ok {
			return nil, fmt.Errorf("missing bucket")
		}

		key, ok := request.Params.Arguments["key"].(string)
		if !ok {
			return nil, fmt.Errorf("missing key")
		}

		executor, err := k.nats.GetExecutor(ctx, accountName)
		if err != nil {
			return nil, err
		}

		args := []string{"kv", "purge", bucket, key}

		// Add force flag if true
		if force, ok := request.Params.Arguments["force"].(bool); ok && force {
			args = append(args, "--force")
		}

		// Add any additional flags passed
		if flags := getFlags(request.Params.Arguments); flags != nil {
			args = append(args, flags...)
		}

		output, err := executor.ExecuteCommand(args...)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}

func (k *KVTools) kvHistoryHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}

		bucket, ok := request.Params.Arguments["bucket"].(string)
		if !ok {
			return nil, fmt.Errorf("missing bucket")
		}

		key, ok := request.Params.Arguments["key"].(string)
		if !ok {
			return nil, fmt.Errorf("missing key")
		}

		executor, err := k.nats.GetExecutor(ctx, accountName)
		if err != nil {
			return nil, err
		}

		args := []string{"kv", "history", bucket, key}
		if flags := getFlags(request.Params.Arguments); flags != nil {
			args = append(args, flags...)
		}

		output, err := executor.ExecuteCommand(args...)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}

func (k *KVTools) kvLsHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}

		executor, err := k.nats.GetExecutor(ctx, accountName)
		if err != nil {
			return nil, err
		}

		args := []string{"kv", "ls"}

		// Add bucket if provided
		if bucket, ok := request.Params.Arguments["bucket"].(string); ok {
			args = append(args, bucket)
		}

		// Add names flag if true
		if names, ok := request.Params.Arguments["names"].(bool); ok && names {
			args = append(args, "--names")
		}

		// Add verbose flag if true
		if verbose, ok := request.Params.Arguments["verbose"].(bool); ok && verbose {
			args = append(args, "--verbose")
		}

		// Add display-value flag if true
		if displayValue, ok := request.Params.Arguments["display_value"].(bool); ok && displayValue {
			args = append(args, "--display-value")
		}

		// Add any additional flags passed
		if flags := getFlags(request.Params.Arguments); flags != nil {
			args = append(args, flags...)
		}

		output, err := executor.ExecuteCommand(args...)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}

func (k *KVTools) kvWatchHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}

		bucket, ok := request.Params.Arguments["bucket"].(string)
		if !ok {
			return nil, fmt.Errorf("missing bucket")
		}

		executor, err := k.nats.GetExecutor(ctx, accountName)
		if err != nil {
			return nil, err
		}

		args := []string{"kv", "watch", bucket}
		if key, ok := request.Params.Arguments["key"].(string); ok {
			args = append(args, key)
		}
		if flags := getFlags(request.Params.Arguments); flags != nil {
			args = append(args, flags...)
		}

		output, err := executor.ExecuteCommand(args...)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}

func (k *KVTools) kvInfoHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}

		executor, err := k.nats.GetExecutor(ctx, accountName)
		if err != nil {
			return nil, err
		}

		args := []string{"kv", "info"}

		// Add bucket if provided
		if bucket, ok := request.Params.Arguments["bucket"].(string); ok {
			args = append(args, bucket)
		}

		// Add any additional flags passed
		if flags := getFlags(request.Params.Arguments); flags != nil {
			args = append(args, flags...)
		}

		output, err := executor.ExecuteCommand(args...)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}

func (k *KVTools) kvCompactHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}

		bucket, ok := request.Params.Arguments["bucket"].(string)
		if !ok {
			return nil, fmt.Errorf("missing bucket")
		}

		executor, err := k.nats.GetExecutor(ctx, accountName)
		if err != nil {
			return nil, err
		}

		args := []string{"kv", "compact", bucket}

		// Add force flag if true
		if force, ok := request.Params.Arguments["force"].(bool); ok && force {
			args = append(args, "--force")
		}

		// Add any additional flags passed
		if flags := getFlags(request.Params.Arguments); flags != nil {
			args = append(args, flags...)
		}

		output, err := executor.ExecuteCommand(args...)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}
