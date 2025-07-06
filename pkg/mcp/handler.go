package mcp

import (
	"reflect"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
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
	// Name returns the name of the MCP tool
	Name() string

	// Description returns the description of the MCP tool
	Description() string

	// InputType returns the input type for the MCP tool request
	InputType() reflect.Type

	// OutputType returns the output type for the MCP tool response
	OutputType() reflect.Type

	// Handle processes the MCP call tool request and returns the response
	Handle(*core.WebContext, *MCPCallToolRequest, *models.User, *settings.Config, MCPAvailableServices) (any, []*T, error)
}
