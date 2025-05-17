package common

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sinadarbouy/mcp-nats/internal/logger"
)

// NATSCreds represents NATS credentials for an account
type NATSCreds struct {
	AccountName string
	Creds       string // base64 encoded credentials
}

// GetCredsFromEnv gets all NATS credentials from environment variables
// Environment variables should be in the format NATS_<ACCOUNT>_CRED
func GetCredsFromEnv() (map[string]NATSCreds, error) {
	creds := make(map[string]NATSCreds)

	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]

		if strings.HasPrefix(key, "NATS_") && strings.HasSuffix(key, "_CREDS") {
			accountName := strings.TrimSuffix(strings.TrimPrefix(key, "NATS_"), "_CREDS")

			credValue := parts[1]

			logger.Debug("Found NATS credentials in environment",
				"account", accountName,
			)

			creds[accountName] = NATSCreds{
				AccountName: accountName,
				Creds:       credValue,
			}
		}
	}

	return creds, nil
}

// NATSExecutor provides common NATS command execution functionality
type NATSExecutor struct {
	URL       string
	Creds     NATSCreds
	credsFile string
	stdin     string
}

// NewNATSExecutor creates a new NATSExecutor instance
func NewNATSExecutor(url string, creds NATSCreds) (*NATSExecutor, error) {
	// Create temporary directory for credentials if it doesn't exist
	tmpDir := filepath.Join(os.TempDir(), "nats-creds")
	if err := os.MkdirAll(tmpDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %v", err)
	}

	// Create temporary file for credentials
	credsFile := filepath.Join(tmpDir, fmt.Sprintf("%s.creds", creds.AccountName))

	// Decode and write base64 credentials to file
	credsData, err := base64.StdEncoding.DecodeString(creds.Creds)
	if err != nil {
		return nil, fmt.Errorf("failed to decode credentials: %v", err)
	}

	if err := os.WriteFile(credsFile, credsData, 0600); err != nil {
		return nil, fmt.Errorf("failed to write credentials file: %v", err)
	}

	logger.Debug("Created NATS credentials file",
		"account", creds.AccountName,
		"file", credsFile,
	)

	return &NATSExecutor{
		URL:       url,
		Creds:     creds,
		credsFile: credsFile,
	}, nil
}

// SetStdin sets the STDIN input for the next command execution
func (e *NATSExecutor) SetStdin(input string) {
	e.stdin = input
}

// ExecuteCommand executes a NATS CLI command with the configured credentials
func (e *NATSExecutor) ExecuteCommand(args ...string) (string, error) {
	baseArgs := []string{"-s", e.URL, "--creds", e.credsFile}
	args = append(baseArgs, args...)

	logger.Debug("Executing NATS command",
		"account", e.Creds.AccountName,
		"command", strings.Join(args, " "),
	)

	cmd := exec.Command("nats", args...)

	// If stdin is set, use it
	if e.stdin != "" {
		cmd.Stdin = strings.NewReader(e.stdin)
		// Clear stdin after use
		defer func() { e.stdin = "" }()
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("NATS command failed",
			"error", err,
			"output", string(output),
			"account", e.Creds.AccountName,
			"command", strings.Join(args, " "),
		)
		return "", fmt.Errorf("NATS command failed: %v, output: %s", err, string(output))
	}

	return string(output), nil
}

// Cleanup removes the temporary credentials file
func (e *NATSExecutor) Cleanup() error {
	if e.credsFile != "" {
		logger.Debug("Cleaning up NATS credentials file",
			"account", e.Creds.AccountName,
			"file", e.credsFile,
		)
		return os.Remove(e.credsFile)
	}
	return nil
}
