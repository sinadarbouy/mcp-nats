package common

import (
	"fmt"
	"os/exec"
)

// NATSExecutor provides common NATS command execution functionality
type NATSExecutor struct {
	URL       string
	CredsPath string
}

// ExecuteCommand executes a NATS CLI command with the configured credentials
func (e *NATSExecutor) ExecuteCommand(args ...string) (string, error) {
	baseArgs := []string{"-s", e.URL, "--creds", e.CredsPath}
	args = append(baseArgs, args...)

	cmd := exec.Command("nats", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("NATS command failed: %v, output: %s", err, string(output))
	}
	return string(output), nil
}
