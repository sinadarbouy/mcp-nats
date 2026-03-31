package tools

// Mutating tools change JetStream/KV/object state, publish messages, or perform
// account backup/restore (filesystem / cluster mutation). When adding a new
// tool that performs writes, register its name here and keep mutatingToolNames
// in sync.
var mutatingToolNames = map[string]struct{}{
	"publish": {},

	"kv_add":     {},
	"kv_put":     {},
	"kv_create":  {},
	"kv_update":  {},
	"kv_del":     {},
	"kv_purge":   {},
	"kv_compact": {},

	"object_add":  {},
	"object_put":  {},
	"object_del":  {},
	"object_seal": {},

	"account_backup":  {},
	"account_restore": {},
}

// IsMutatingTool reports whether the named MCP tool mutates NATS/JetStream
// state, publishes data, or writes local backup output.
func IsMutatingTool(name string) bool {
	_, ok := mutatingToolNames[name]
	return ok
}

// MutatingToolCount returns the number of registered mutating tool names.
func MutatingToolCount() int {
	return len(mutatingToolNames)
}
