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

// NATSAuthStrategy defines the interface for different authentication strategies
type NATSAuthStrategy interface {
	BuildArgs(baseURL string) []string
	GetAccountName() string
	Cleanup() error
}

// AnonymousAuthStrategy implements anonymous authentication
type AnonymousAuthStrategy struct {
	accountName string
}

func NewAnonymousAuthStrategy() *AnonymousAuthStrategy {
	return &AnonymousAuthStrategy{
		accountName: "anonymous",
	}
}

func (a *AnonymousAuthStrategy) BuildArgs(baseURL string) []string {
	return []string{"-s", baseURL}
}

func (a *AnonymousAuthStrategy) GetAccountName() string {
	return a.accountName
}

func (a *AnonymousAuthStrategy) Cleanup() error {
	return nil // No cleanup needed for anonymous auth
}

// UserPassAuthStrategy implements user/password authentication
type UserPassAuthStrategy struct {
	user        string
	password    string
	accountName string
}

// NewUserPassAuthStrategy creates a new UserPassAuthStrategy instance
func NewUserPassAuthStrategy(user, password string) *UserPassAuthStrategy {
	return &UserPassAuthStrategy{
		user:        user,
		password:    password,
		accountName: fmt.Sprintf("userpass_%s", user),
	}
}

// BuildArgs builds the arguments for the NATS CLI command
func (u *UserPassAuthStrategy) BuildArgs(baseURL string) []string {
	return []string{"-s", baseURL, "--user", u.user, "--password", u.password}
}

// GetAccountName returns the account name for this authentication strategy
func (u *UserPassAuthStrategy) GetAccountName() string {
	return u.accountName
}

// Cleanup cleans up any resources used by this authentication strategy
func (u *UserPassAuthStrategy) Cleanup() error {
	return nil // No cleanup needed for user/pass auth
}

// CredentialsAuthStrategy implements credentials-based authentication
type CredentialsAuthStrategy struct {
	creds       NATSCreds
	credsFile   string
	accountName string
}

// NewCredentialsAuthStrategy creates a new CredentialsAuthStrategy instance
func NewCredentialsAuthStrategy(creds NATSCreds) (*CredentialsAuthStrategy, error) {
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

	return &CredentialsAuthStrategy{
		creds:       creds,
		credsFile:   credsFile,
		accountName: creds.AccountName,
	}, nil
}

// BuildArgs builds the arguments for the NATS CLI command
func (c *CredentialsAuthStrategy) BuildArgs(baseURL string) []string {
	return []string{"-s", baseURL, "--creds", c.credsFile}
}

// GetAccountName returns the account name for this authentication strategy
func (c *CredentialsAuthStrategy) GetAccountName() string {
	return c.accountName
}

// Cleanup cleans up any resources used by this authentication strategy
func (c *CredentialsAuthStrategy) Cleanup() error {
	if c.credsFile != "" {
		logger.Debug("Cleaning up NATS credentials file",
			"account", c.accountName,
			"file", c.credsFile,
		)
		return os.Remove(c.credsFile)
	}
	return nil
}

// NATSExecutor provides common NATS command execution functionality
type NATSExecutor struct {
	URL      string
	Strategy NATSAuthStrategy
	stdin    string
}

// NewNATSExecutor creates a new NATSExecutor instance with credentials
func NewNATSExecutor(url string, creds NATSCreds) (*NATSExecutor, error) {
	strategy, err := NewCredentialsAuthStrategy(creds)
	if err != nil {
		return nil, err
	}

	return &NATSExecutor{
		URL:      url,
		Strategy: strategy,
	}, nil
}

// NewAnonymousNATSExecutor creates a new NATSExecutor instance for anonymous connections
func NewAnonymousNATSExecutor(url string) *NATSExecutor {
	return &NATSExecutor{
		URL:      url,
		Strategy: NewAnonymousAuthStrategy(),
	}
}

// NewUserPassNATSExecutor creates a new NATSExecutor instance with user/password authentication
func NewUserPassNATSExecutor(url, user, password string) *NATSExecutor {
	return &NATSExecutor{
		URL:      url,
		Strategy: NewUserPassAuthStrategy(user, password),
	}
}

// SetStdin sets the STDIN input for the next command execution
func (e *NATSExecutor) SetStdin(input string) {
	e.stdin = input
}

// ExecuteCommand executes a NATS CLI command with the configured authentication
func (e *NATSExecutor) ExecuteCommand(args ...string) (string, error) {
	baseArgs := e.Strategy.BuildArgs(e.URL)
	args = append(baseArgs, args...)

	logger.Debug("Executing NATS command",
		"account", e.Strategy.GetAccountName(),
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
			"account", e.Strategy.GetAccountName(),
			"command", strings.Join(args, " "),
		)
		return "", fmt.Errorf("NATS command failed: %v, output: %s", err, string(output))
	}

	return string(output), nil
}

// Cleanup removes the temporary credentials file
func (e *NATSExecutor) Cleanup() error {
	return e.Strategy.Cleanup()
}

// GetAccountName returns the account name for this executor
func (e *NATSExecutor) GetAccountName() string {
	return e.Strategy.GetAccountName()
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

// GetUserPassFromEnv gets NATS user/password from environment variables
func GetUserPassFromEnv() (string, string) {
	user := os.Getenv("NATS_USER")
	password := os.Getenv("NATS_PASSWORD")
	return user, password
}

// GetAuthStrategy determines the current authentication strategy
func GetAuthStrategy() string {
	// Check for no-authentication flag
	if os.Getenv("NATS_NO_AUTHENTICATION") == "true" {
		return "anonymous"
	}

	// Check for user/password authentication
	user, password := GetUserPassFromEnv()
	if user != "" && password != "" {
		return "userpass"
	}

	// Default to credentials-based authentication
	return "credentials"
}

// IsAccountNameRequired determines if account_name is required based on auth strategy
func IsAccountNameRequired() bool {
	strategy := GetAuthStrategy()
	// Only credentials-based authentication requires account_name
	return strategy == "credentials"
}

// DetermineAccountName determines the account name to use based on authentication strategy and request parameters
func DetermineAccountName(requestArgs map[string]interface{}) (string, error) {
	if IsAccountNameRequired() {
		// For credentials-based auth, account_name is required from request
		accountNameVal, ok := requestArgs["account_name"].(string)
		if !ok {
			return "", fmt.Errorf("missing account_name")
		}
		return accountNameVal, nil
	} else {
		// For non-credentials auth, use strategy's account name
		strategy := GetAuthStrategy()
		switch strategy {
		case "anonymous":
			return "anonymous", nil
		case "userpass":
			user, _ := GetUserPassFromEnv()
			return fmt.Sprintf("userpass_%s", user), nil
		default:
			return "default", nil
		}
	}
}
