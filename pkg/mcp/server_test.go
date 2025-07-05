package mcp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

func TestMCPInitialize(t *testing.T) {
	server := NewMCPServer()
	
	// Create a test request
	initReq := JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "initialize",
		Params: MCPInitializeRequest{
			ProtocolVersion: "2024-11-05",
			Capabilities: MCPCapabilities{
				Resources: &MCPResourceCapabilities{},
				Tools:     &MCPToolCapabilities{},
			},
			ClientInfo: MCPClientInfo{
				Name:    "test-client",
				Version: "1.0.0",
			},
		},
		ID: 1,
	}

	jsonData, _ := json.Marshal(initReq)

	// Create test HTTP request
	req, _ := http.NewRequest("POST", "/mcp", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Create test recorder
	recorder := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(recorder)
	ginCtx.Request = req

	// Mock user ID
	ginCtx.Set("uid", int64(1))

	webCtx := core.WrapWebContext(ginCtx)

	// Mock user for testing
	user := &models.User{
		Uid: 1,
	}

	response, err := server.handleInitialize(webCtx, &initReq, user)
	
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "2.0", response.JSONRPC)
	assert.Equal(t, 1, response.ID)
	
	// Check the response result
	result, ok := response.Result.(MCPInitializeResponse)
	assert.True(t, ok)
	assert.Equal(t, "2024-11-05", result.ProtocolVersion)
	assert.Equal(t, "ezbookkeeping-mcp", result.ServerInfo.Name)
}

func TestMCPListResources(t *testing.T) {
	server := NewMCPServer()
	
	// Create a test request
	listReq := JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "resources/list",
		ID:      2,
	}

	jsonData, _ := json.Marshal(listReq)

	// Create test HTTP request
	req, _ := http.NewRequest("POST", "/mcp", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Create test recorder
	recorder := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(recorder)
	ginCtx.Request = req

	webCtx := core.WrapWebContext(ginCtx)

	// Mock user for testing
	user := &models.User{
		Uid: 1,
	}

	response, err := server.handleListResources(webCtx, &listReq, user)
	
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "2.0", response.JSONRPC)
	assert.Equal(t, 2, response.ID)
	
	// Check the response result
	result, ok := response.Result.(MCPListResourcesResponse)
	assert.True(t, ok)
	assert.Greater(t, len(result.Resources), 0)
	
	// Check that we have expected resources
	resourceURIs := make(map[string]bool)
	for _, resource := range result.Resources {
		resourceURIs[resource.URI] = true
	}
	
	assert.True(t, resourceURIs["transactions://recent"])
	assert.True(t, resourceURIs["accounts://list"])
	assert.True(t, resourceURIs["categories://list"])
	assert.True(t, resourceURIs["tags://list"])
}

func TestMCPListTools(t *testing.T) {
	server := NewMCPServer()
	
	// Create a test request
	listReq := JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "tools/list",
		ID:      3,
	}

	jsonData, _ := json.Marshal(listReq)

	// Create test HTTP request
	req, _ := http.NewRequest("POST", "/mcp", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Create test recorder
	recorder := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(recorder)
	ginCtx.Request = req

	webCtx := core.WrapWebContext(ginCtx)

	// Mock user for testing
	user := &models.User{
		Uid: 1,
	}

	response, err := server.handleListTools(webCtx, &listReq, user)
	
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "2.0", response.JSONRPC)
	assert.Equal(t, 3, response.ID)
	
	// Check the response result
	result, ok := response.Result.(MCPListToolsResponse)
	assert.True(t, ok)
	assert.Greater(t, len(result.Tools), 0)
	
	// Check that we have expected tools
	toolNames := make(map[string]bool)
	for _, tool := range result.Tools {
		toolNames[tool.Name] = true
	}
	
	assert.True(t, toolNames["query_transactions"])
	assert.True(t, toolNames["get_transaction_statistics"])
	assert.True(t, toolNames["get_account_balance"])
}