package tools

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	mcpnats "github.com/sinadarbouy/mcp-nats"
	"github.com/sinadarbouy/mcp-nats/internal/logger"
	"github.com/sinadarbouy/mcp-nats/test/utils/containers"
	"github.com/sinadarbouy/mcp-nats/tools/common"
	"github.com/stretchr/testify/suite"
)

// AccountTestSuite defines the test suite for account-related functionality
type AccountTestSuite struct {
	suite.Suite
	ctx           context.Context
	natsContainer *containers.NatsContainer
	natsURL       string
	accountTools  *AccountTools
}

func TestAccountSuite(t *testing.T) {
	suite.Run(t, new(AccountTestSuite))
}

func (s *AccountTestSuite) SetupSuite() {
	// Initialize logger for tests
	logger.Initialize(logger.Config{
		Level:      logger.LevelDebug,
		JSONFormat: false,
	})

	// Setup context
	s.ctx = context.Background()

	// Start NATS container
	s.natsContainer = containers.NewNatsContainer(s.ctx, s.T())

	// Get NATS URL
	natsHost := s.natsContainer.Host
	natsPort := s.natsContainer.Port
	s.natsURL = fmt.Sprintf("nats://%s:%s", natsHost, natsPort)

	// Create test credentials for the SYS account
	// For test containers, we can use empty credentials since auth is not required
	testCreds := map[string]common.NATSCreds{
		"SYS": {
			AccountName: "SYS",
			Creds:       base64.StdEncoding.EncodeToString([]byte("")), // Empty creds for test container
		},
	}

	// Add NATS URL and credentials to context
	s.ctx = mcpnats.WithNatsURL(s.ctx, s.natsURL)
	s.ctx = mcpnats.WithNatsCreds(s.ctx, testCreds)

	// Create NATSServerTools instance
	natsTools := &NATSServerTools{
		executors: make(map[string]*common.NATSExecutor),
	}

	// Create AccountTools instance
	s.accountTools = NewAccountTools(natsTools)
}

func (s *AccountTestSuite) TestAccountInfoHandler() {
	// Get the account info handler
	handler := s.accountTools.accountInfoHandler()

	// Create a mock request with account_name
	request := mcp.CallToolRequest{
		Params: struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments,omitempty"`
			Meta      *struct {
				ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
			} `json:"_meta,omitempty"`
		}{
			Name: "account_info",
			Arguments: map[string]interface{}{
				"account_name": "SYS", // Using SYS account which should exist by default
			},
		},
	}

	// Call the handler - this should now work with the proper context
	result, err := handler(s.ctx, request)

	// Assertions for successful account info request
	s.Assert().NoError(err, "Account info handler should not return an error")
	s.Assert().NotNil(result, "Result should not be nil")
	s.Assert().NotEmpty(result.Content, "Result content should not be empty")

	// Verify that we got some account info back
	if len(result.Content) > 0 {
		if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
			s.Assert().NotEmpty(textContent.Text, "Account info text should not be empty")
			s.T().Logf("Account info response: %s", textContent.Text)
		}
	}
}

func (s *AccountTestSuite) TestAccountReportConnectionsHandler() {
	// Get the account report connections handler
	handler := s.accountTools.accountReportConnectionsHandler()

	// Create a mock request with account_name
	request := mcp.CallToolRequest{
		Params: struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments,omitempty"`
			Meta      *struct {
				ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
			} `json:"_meta,omitempty"`
		}{
			Name: "account_report_connections",
			Arguments: map[string]interface{}{
				"account_name": "SYS",
			},
		},
	}

	// Call the handler
	result, err := handler(s.ctx, request)

	// Assertions
	s.Assert().NoError(err, "Account report connections handler should not return an error")
	s.Assert().NotNil(result, "Result should not be nil")
	s.Assert().NotEmpty(result.Content, "Result content should not be empty")

	// Verify that we got some report back
	if len(result.Content) > 0 {
		if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
			s.Assert().NotEmpty(textContent.Text, "Account report connections text should not be empty")
			s.T().Logf("Account report connections response: %s", textContent.Text)
		}
	}
}

func (s *AccountTestSuite) TestAccountReportConnectionsHandlerWithOptions() {
	// Get the account report connections handler
	handler := s.accountTools.accountReportConnectionsHandler()

	// Create a mock request with account_name and options
	request := mcp.CallToolRequest{
		Params: struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments,omitempty"`
			Meta      *struct {
				ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
			} `json:"_meta,omitempty"`
		}{
			Name: "account_report_connections",
			Arguments: map[string]interface{}{
				"account_name": "SYS",
				"sort":         "subs",
				"top":          float64(10),
			},
		},
	}

	// Call the handler
	result, err := handler(s.ctx, request)

	// Assertions
	s.Assert().NoError(err, "Account report connections handler with options should not return an error")
	s.Assert().NotNil(result, "Result should not be nil")
	s.Assert().NotEmpty(result.Content, "Result content should not be empty")

	// Verify that we got some report back
	if len(result.Content) > 0 {
		if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
			s.Assert().NotEmpty(textContent.Text, "Account report connections text should not be empty")
			s.T().Logf("Account report connections with options response: %s", textContent.Text)
		}
	}
}

func (s *AccountTestSuite) TestAccountReportStatisticsHandler() {
	// Get the account report statistics handler
	handler := s.accountTools.accountReportStatisticsHandler()

	// Create a mock request with account_name
	request := mcp.CallToolRequest{
		Params: struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments,omitempty"`
			Meta      *struct {
				ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
			} `json:"_meta,omitempty"`
		}{
			Name: "account_report_statistics",
			Arguments: map[string]interface{}{
				"account_name": "SYS",
			},
		},
	}

	// Call the handler
	result, err := handler(s.ctx, request)

	// Assertions
	s.Assert().NoError(err, "Account report statistics handler should not return an error")
	s.Assert().NotNil(result, "Result should not be nil")
	s.Assert().NotEmpty(result.Content, "Result content should not be empty")

	// Verify that we got some statistics back
	if len(result.Content) > 0 {
		if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
			s.Assert().NotEmpty(textContent.Text, "Account report statistics text should not be empty")
			s.T().Logf("Account report statistics response: %s", textContent.Text)
		}
	}
}

func (s *AccountTestSuite) TestAccountBackupHandler() {
	// Create a temporary directory for backup target
	tempDir, err := os.MkdirTemp("", "nats_backup_test_*")
	s.Require().NoError(err, "Failed to create temporary directory")
	defer func() { _ = os.RemoveAll(tempDir) }() // Clean up after test

	// Get the account backup handler
	handler := s.accountTools.accountBackupHandler()

	// Create a mock request with account_name and target
	request := mcp.CallToolRequest{
		Params: struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments,omitempty"`
			Meta      *struct {
				ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
			} `json:"_meta,omitempty"`
		}{
			Name: "account_backup",
			Arguments: map[string]interface{}{
				"account_name": "SYS",
				"target":       tempDir,
				"force":        true, // Use force to avoid prompts in tests
			},
		},
	}

	// Call the handler
	result, err := handler(s.ctx, request)

	// For backup, we expect either success or a controlled error (no streams to backup)
	// The important thing is that the handler executes without crashing
	if err == nil {
		// If successful, verify we got a valid result and content back
		s.Assert().NotNil(result, "Result should not be nil when successful")
		s.Assert().NotEmpty(result.Content, "Result content should not be empty")
		if len(result.Content) > 0 {
			if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
				s.Assert().NotEmpty(textContent.Text, "Account backup text should not be empty")
				s.T().Logf("Account backup response: %s", textContent.Text)
			}
		}
	} else {
		// If there's an error (e.g., no streams found), it should be a controlled error
		// The result will be nil in this case, which is expected
		s.Assert().Nil(result, "Result should be nil when there's an error")
		s.T().Logf("Account backup returned expected error (no streams found): %v", err)
	}
}

func (s *AccountTestSuite) TestAccountBackupHandlerWithOptions() {
	// Create a temporary directory for backup target
	tempDir, err := os.MkdirTemp("", "nats_backup_test_with_options_*")
	s.Require().NoError(err, "Failed to create temporary directory")
	defer func() { _ = os.RemoveAll(tempDir) }() // Clean up after test

	// Get the account backup handler
	handler := s.accountTools.accountBackupHandler()

	// Create a mock request with account_name, target, and various options
	request := mcp.CallToolRequest{
		Params: struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments,omitempty"`
			Meta      *struct {
				ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
			} `json:"_meta,omitempty"`
		}{
			Name: "account_backup",
			Arguments: map[string]interface{}{
				"account_name": "SYS",
				"target":       tempDir,
				"check":        true,
				"consumers":    true,
				"force":        true,
			},
		},
	}

	// Call the handler
	result, err := handler(s.ctx, request)

	// For backup with options, we expect either success or a controlled error (no streams to backup)
	// The important thing is that the handler executes without crashing
	if err == nil {
		// If successful, verify we got a valid result and content back
		s.Assert().NotNil(result, "Result should not be nil when successful")
		s.Assert().NotEmpty(result.Content, "Result content should not be empty")
		if len(result.Content) > 0 {
			if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
				s.Assert().NotEmpty(textContent.Text, "Account backup text should not be empty")
				s.T().Logf("Account backup with options response: %s", textContent.Text)
			}
		}
	} else {
		// If there's an error (e.g., no streams found), it should be a controlled error
		// The result will be nil in this case, which is expected
		s.Assert().Nil(result, "Result should be nil when there's an error")
		s.T().Logf("Account backup with options returned expected error (no streams found): %v", err)
	}
}

func (s *AccountTestSuite) TestAccountRestoreHandler() {
	// Create a temporary directory structure that simulates a backup directory
	tempDir, err := os.MkdirTemp("", "nats_restore_test_*")
	s.Require().NoError(err, "Failed to create temporary directory")
	defer func() { _ = os.RemoveAll(tempDir) }() // Clean up after test

	// Create a mock backup directory structure
	backupDir := filepath.Join(tempDir, "backup")
	err = os.MkdirAll(backupDir, 0755)
	s.Require().NoError(err, "Failed to create backup directory")

	// Get the account restore handler
	handler := s.accountTools.accountRestoreHandler()

	// Create a mock request with account_name and directory
	request := mcp.CallToolRequest{
		Params: struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments,omitempty"`
			Meta      *struct {
				ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
			} `json:"_meta,omitempty"`
		}{
			Name: "account_restore",
			Arguments: map[string]interface{}{
				"account_name": "SYS",
				"directory":    backupDir,
			},
		},
	}

	// Call the handler - Note: This might fail if there's no actual backup data,
	// but we're testing that the handler can be called without panicking
	result, err := handler(s.ctx, request)

	// For restore, we expect either success or a controlled error (not a panic)
	// The important thing is that the handler executes without crashing
	if err == nil {
		// If successful, verify we got a valid result and content back
		s.Assert().NotNil(result, "Result should not be nil when successful")
		s.Assert().NotEmpty(result.Content, "Result content should not be empty")
		if len(result.Content) > 0 {
			if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
				s.Assert().NotEmpty(textContent.Text, "Account restore text should not be empty")
				s.T().Logf("Account restore response: %s", textContent.Text)
			}
		}
	} else {
		// If there's an error, it should be a controlled error, not a panic
		// The result will be nil in this case, which is expected
		s.Assert().Nil(result, "Result should be nil when there's an error")
		s.T().Logf("Account restore returned expected error (no backup data): %v", err)
	}
}

func (s *AccountTestSuite) TestAccountRestoreHandlerWithOptions() {
	// Create a temporary directory structure that simulates a backup directory
	tempDir, err := os.MkdirTemp("", "nats_restore_test_with_options_*")
	s.Require().NoError(err, "Failed to create temporary directory")
	defer func() { _ = os.RemoveAll(tempDir) }() // Clean up after test

	// Create a mock backup directory structure
	backupDir := filepath.Join(tempDir, "backup")
	err = os.MkdirAll(backupDir, 0755)
	s.Require().NoError(err, "Failed to create backup directory")

	// Get the account restore handler
	handler := s.accountTools.accountRestoreHandler()

	// Create a mock request with account_name, directory, and options
	request := mcp.CallToolRequest{
		Params: struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments,omitempty"`
			Meta      *struct {
				ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
			} `json:"_meta,omitempty"`
		}{
			Name: "account_restore",
			Arguments: map[string]interface{}{
				"account_name": "SYS",
				"directory":    backupDir,
				"cluster":      "test-cluster",
				"tags":         []interface{}{"tag1", "tag2"},
			},
		},
	}

	// Call the handler
	result, err := handler(s.ctx, request)

	// For restore with options, we expect either success or a controlled error
	if err == nil {
		// If successful, verify we got a valid result and content back
		s.Assert().NotNil(result, "Result should not be nil when successful")
		s.Assert().NotEmpty(result.Content, "Result content should not be empty")
		if len(result.Content) > 0 {
			if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
				s.Assert().NotEmpty(textContent.Text, "Account restore text should not be empty")
				s.T().Logf("Account restore with options response: %s", textContent.Text)
			}
		}
	} else {
		// If there's an error, it should be a controlled error, not a panic
		// The result will be nil in this case, which is expected
		s.Assert().Nil(result, "Result should be nil when there's an error")
		s.T().Logf("Account restore with options returned expected error (no backup data): %v", err)
	}
}

func (s *AccountTestSuite) TestAccountTLSHandler() {
	// Get the account TLS handler
	handler := s.accountTools.accountTLSHandler()

	// Create a mock request with account_name
	request := mcp.CallToolRequest{
		Params: struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments,omitempty"`
			Meta      *struct {
				ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
			} `json:"_meta,omitempty"`
		}{
			Name: "account_tls",
			Arguments: map[string]interface{}{
				"account_name": "SYS",
			},
		},
	}

	// Call the handler
	result, err := handler(s.ctx, request)

	// Assertions - TLS might not be configured in test environment, so we accept both success and error
	if err == nil {
		// If successful, verify we got a valid result and content back
		s.Assert().NotNil(result, "Result should not be nil when successful")
		s.Assert().NotEmpty(result.Content, "Result content should not be empty")
		if len(result.Content) > 0 {
			if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
				s.Assert().NotEmpty(textContent.Text, "Account TLS text should not be empty")
				s.T().Logf("Account TLS response: %s", textContent.Text)
			}
		}
	} else {
		// If there's an error (e.g., no TLS configured), it should be a controlled error
		// The result will be nil in this case, which is expected
		s.Assert().Nil(result, "Result should be nil when there's an error")
		s.T().Logf("Account TLS returned expected error (no TLS configured): %v", err)
	}
}

func (s *AccountTestSuite) TestAccountTLSHandlerWithOptions() {
	// Get the account TLS handler
	handler := s.accountTools.accountTLSHandler()

	// Create a mock request with account_name and options
	request := mcp.CallToolRequest{
		Params: struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments,omitempty"`
			Meta      *struct {
				ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
			} `json:"_meta,omitempty"`
		}{
			Name: "account_tls",
			Arguments: map[string]interface{}{
				"account_name": "SYS",
				"expire_warn":  "1w",
				"ocsp":         true,
				"pem":          true,
			},
		},
	}

	// Call the handler
	result, err := handler(s.ctx, request)

	// Assertions - TLS might not be configured in test environment, so we accept both success and error
	if err == nil {
		// If successful, verify we got a valid result and content back
		s.Assert().NotNil(result, "Result should not be nil when successful")
		s.Assert().NotEmpty(result.Content, "Result content should not be empty")
		if len(result.Content) > 0 {
			if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
				s.Assert().NotEmpty(textContent.Text, "Account TLS text should not be empty")
				s.T().Logf("Account TLS with options response: %s", textContent.Text)
			}
		}
	} else {
		// If there's an error (e.g., no TLS configured), it should be a controlled error
		// The result will be nil in this case, which is expected
		s.Assert().Nil(result, "Result should be nil when there's an error")
		s.T().Logf("Account TLS with options returned expected error (no TLS configured): %v", err)
	}
}

func (s *AccountTestSuite) TearDownSuite() {
	if s.natsContainer != nil {
		_ = s.natsContainer.Container.Terminate(s.ctx)
	}
}
