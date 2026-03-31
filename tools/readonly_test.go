package tools

import "testing"

func TestIsMutatingTool_mutatingSet(t *testing.T) {
	for name := range mutatingToolNames {
		if !IsMutatingTool(name) {
			t.Errorf("IsMutatingTool(%q) = false, want true", name)
		}
	}
}

func TestIsMutatingTool_readOnlySafe(t *testing.T) {
	safe := []string{
		"stream_list",
		"kv_get",
		"account_info",
		"server_list",
		"rtt",
		"object_ls",
	}
	for _, name := range safe {
		if IsMutatingTool(name) {
			t.Errorf("IsMutatingTool(%q) = true, want false", name)
		}
	}
}

func TestMutatingToolCount(t *testing.T) {
	if got, want := MutatingToolCount(), len(mutatingToolNames); got != want {
		t.Fatalf("MutatingToolCount() = %d, want %d", got, want)
	}
}

func TestToolCount_readOnlySkipsMutating(t *testing.T) {
	n, err := NewNATSServerTools()
	if err != nil {
		t.Fatalf("NewNATSServerTools: %v", err)
	}

	full := ToolCount(n, false)
	readOnly := ToolCount(n, true)
	if got := full - readOnly; got != MutatingToolCount() {
		t.Fatalf("full - readOnly = %d, want MutatingToolCount() = %d", got, MutatingToolCount())
	}

	seen := make(map[string]int)
	for _, cat := range n.toolCategories() {
		for _, tool := range cat.GetTools() {
			seen[tool.Tool.Name]++
		}
	}
	for name := range mutatingToolNames {
		if seen[name] != 1 {
			t.Fatalf("mutating tool %q appears %d times in catalog, want 1", name, seen[name])
		}
	}
}
