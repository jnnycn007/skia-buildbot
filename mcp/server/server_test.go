package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.skia.org/infra/mcp/common"
	"go.skia.org/infra/mcp/common/mocks"
)

func TestCreateServer_Success(t *testing.T) {
	flags := &mcpFlags{
		ServiceName: string(HelloWorld),
	}

	server, err := createMcpSSEServer(flags)
	assert.Nil(t, err)
	assert.NotNil(t, server)
}

func TestCreateServer_Invalid(t *testing.T) {
	flags := &mcpFlags{
		ServiceName: "random",
	}

	server, err := createMcpSSEServer(flags)
	assert.NotNil(t, err)
	assert.Nil(t, server)
}

func TestCreateMcpSSEServer_ArgumentTypeSwitch(t *testing.T) {
	originalServiceRegistry := make(map[string]serviceFactory)
	for k, v := range serviceRegistry {
		originalServiceRegistry[k] = v
	}
	defer func() {
		serviceRegistry = originalServiceRegistry
	}()

	testCases := []struct {
		name             string
		argType          common.ToolArgumentType
		expectError      bool
		expectedErrorMsg string
	}{
		{name: "StringArgument", argType: common.StringArgument, expectError: false},
		{name: "BooleanArgument", argType: common.BooleanArgument, expectError: false},
		{name: "NumberArgument", argType: common.NumberArgument, expectError: false},
		{name: "ObjectArgument", argType: common.ObjectArgument, expectError: false},
		{name: "ArrayArgument", argType: common.ArrayArgument, expectError: false},
		{name: "InvalidArgumentType", argType: common.ToolArgumentType(99), expectError: true, expectedErrorMsg: "Invalid argument type 99"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := &mocks.MockArgumentService{ArgTypeToTest: tc.argType}
			testServiceName := "testargumentservice_" + tc.name

			serviceRegistry[testServiceName] = func() common.McpService { return mockService }
			defer delete(serviceRegistry, testServiceName)

			flags := &mcpFlags{
				ServiceName: testServiceName,
			}

			server, err := createMcpSSEServer(flags)

			if tc.expectError {
				require.Error(t, err)
				assert.Nil(t, server)
				if tc.expectedErrorMsg != "" {
					assert.Contains(t, err.Error(), tc.expectedErrorMsg)
				}
			} else {
				require.NoError(t, err)
				assert.NotNil(t, server)
			}
		})
	}
}

func TestCreateMcpSSEServer_ToolWithNoArguments(t *testing.T) {
	originalServiceRegistry := make(map[string]serviceFactory)
	for k, v := range serviceRegistry {
		originalServiceRegistry[k] = v
	}
	defer func() {
		serviceRegistry = originalServiceRegistry
	}()

	mockService := &mocks.MockArgumentService{
		CustomTools: []common.Tool{
			{Name: "noArgTool", Description: "Tool with no arguments", Arguments: nil, Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) { return nil, nil }},
		},
	}
	testServiceName := "testnoargservice"
	serviceRegistry[testServiceName] = func() common.McpService { return mockService }
	defer delete(serviceRegistry, testServiceName)

	flags := &mcpFlags{ServiceName: testServiceName}
	server, err := createMcpSSEServer(flags)
	require.NoError(t, err)
	assert.NotNil(t, server)
}

func TestCreateMcpSSEServer_ServiceWithNoTools(t *testing.T) {
	originalServiceRegistry := make(map[string]serviceFactory)
	for k, v := range serviceRegistry {
		originalServiceRegistry[k] = v
	}
	defer func() {
		serviceRegistry = originalServiceRegistry
	}()

	mockService := &mocks.MockArgumentService{CustomTools: []common.Tool{}} // No tools
	testServiceName := "testnotoolservice"
	serviceRegistry[testServiceName] = func() common.McpService { return mockService }
	defer delete(serviceRegistry, testServiceName)

	flags := &mcpFlags{ServiceName: testServiceName}
	server, err := createMcpSSEServer(flags)
	require.NoError(t, err)
	assert.NotNil(t, server)
}

func TestCreateMcpSSEServer_ServiceInitError(t *testing.T) {
	// This test doesn't need to manipulate the global serviceRegistry if the service is already registered or if we register a new one.
	// For consistency with other tests, we can follow the same pattern.
	expectedErr := fmt.Errorf("init failed")
	mockService := &mocks.MockArgumentService{InitError: expectedErr}
	testServiceName := "testiniterrorservice" // This service name won't be in the default registry

	serviceRegistry[testServiceName] = func() common.McpService { return mockService }
	defer delete(serviceRegistry, testServiceName) // Clean up after test

	flags := &mcpFlags{ServiceName: testServiceName}
	server, err := createMcpSSEServer(flags)
	require.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, server)
}
