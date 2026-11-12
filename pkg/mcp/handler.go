package mcp

import (
	"reflect"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// MCPToolDataOperationType represents how an MCP tool operates on business records
type MCPToolDataOperationType uint8

// MCP tool data operation types (modification includes deletion)
const (
	MCP_TOOL_DATA_OPERATION_READ_ONLY MCPToolDataOperationType = 1
	MCP_TOOL_DATA_OPERATION_ADD       MCPToolDataOperationType = 2
	MCP_TOOL_DATA_OPERATION_MODIFY    MCPToolDataOperationType = 3
)

// MCPAvailableServices holds the services available for MCP tools
type MCPAvailableServices interface {
	GetTransactionService() *services.TransactionService
	GetTransactionCategoryService() *services.TransactionCategoryService
	GetTransactionTagService() *services.TransactionTagService
	GetAccountService() *services.AccountService
	GetUserService() *services.UserService
}

// MCPToolHandler defines the MCP tool handler
type MCPToolHandler[T MCPTextContent | MCPImageContent | MCPAudioContent | MCPResourceLink | MCPEmbeddedResource] interface {
	// ToolType returns the MCP tool type
	ToolType() core.MCPToolType

	// Name returns the name of the MCP tool
	Name() string

	// Description returns the description of the MCP tool
	Description() string

	// DataOperationType returns how the MCP tool operates on business records
	DataOperationType() MCPToolDataOperationType

	// IsIdempotent returns whether repeated calls with the same arguments have no additional effects
	IsIdempotent() bool

	// DependsOnExternalData returns whether the MCP tool may depend on third-party data
	DependsOnExternalData() bool

	// InputType returns the input type for the MCP tool request
	InputType() reflect.Type

	// OutputType returns the output type for the MCP tool response
	OutputType() reflect.Type

	// Handle processes the MCP call tool request and returns the response
	Handle(*core.WebContext, *MCPCallToolRequest, *models.User, *settings.Config, MCPAvailableServices) (any, []*T, error)
}
