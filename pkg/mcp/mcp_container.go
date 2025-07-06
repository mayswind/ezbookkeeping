package mcp

import (
	"github.com/invopop/jsonschema"
	orderedmap "github.com/wk8/go-ordered-map/v2"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// MCPContainer contains the all mcp handlers
type MCPContainer struct {
	mcpTextContentTools      *orderedmap.OrderedMap[string, MCPToolHandler[MCPTextContent]]
	mcpImageContentTools     *orderedmap.OrderedMap[string, MCPToolHandler[MCPImageContent]]
	mcpAudioContentTools     *orderedmap.OrderedMap[string, MCPToolHandler[MCPAudioContent]]
	mcpResourceLinkTools     *orderedmap.OrderedMap[string, MCPToolHandler[MCPResourceLink]]
	mcpEmbeddedResourceTools *orderedmap.OrderedMap[string, MCPToolHandler[MCPEmbeddedResource]]
	mcpTools                 []*MCPTool
}

// Initialize a mcp handler container singleton instance
var (
	Container = &MCPContainer{}
)

// GetMCPTools returns the registered MCP tools
func (c *MCPContainer) GetMCPTools() []*MCPTool {
	if len(c.mcpTools) == 0 {
		return nil
	}

	return c.mcpTools
}

// HandleTool returns the result of the MCP tool handler based on the tool name
func (c *MCPContainer) HandleTool(ctx *core.WebContext, callToolReq *MCPCallToolRequest, currentConfig *settings.Config, services MCPAvailableServices) (any, *errs.Error) {
	if handler, exists := c.mcpTextContentTools.Get(callToolReq.Name); exists {
		return handleTool(ctx, handler, currentConfig, services, callToolReq)
	}

	if handler, exists := c.mcpImageContentTools.Get(callToolReq.Name); exists {
		return handleTool(ctx, handler, currentConfig, services, callToolReq)
	}

	if handler, exists := c.mcpAudioContentTools.Get(callToolReq.Name); exists {
		return handleTool(ctx, handler, currentConfig, services, callToolReq)
	}

	if handler, exists := c.mcpResourceLinkTools.Get(callToolReq.Name); exists {
		return handleTool(ctx, handler, currentConfig, services, callToolReq)
	}

	if handler, exists := c.mcpEmbeddedResourceTools.Get(callToolReq.Name); exists {
		return handleTool(ctx, handler, currentConfig, services, callToolReq)
	}

	return nil, errs.ErrApiNotFound
}

// InitializeMCPHandlers initializes the all mcp handlers according to the config
func InitializeMCPHandlers(config *settings.Config) error {
	container := &MCPContainer{
		mcpTextContentTools:      orderedmap.New[string, MCPToolHandler[MCPTextContent]](),
		mcpImageContentTools:     orderedmap.New[string, MCPToolHandler[MCPImageContent]](),
		mcpAudioContentTools:     orderedmap.New[string, MCPToolHandler[MCPAudioContent]](),
		mcpResourceLinkTools:     orderedmap.New[string, MCPToolHandler[MCPResourceLink]](),
		mcpEmbeddedResourceTools: orderedmap.New[string, MCPToolHandler[MCPEmbeddedResource]](),
		mcpTools:                 make([]*MCPTool, 0),
	}

	registerMCPTextContentToolHandler(container, MCPQueryAllAccountsToolHandler)
	registerMCPTextContentToolHandler(container, MCPQueryAllTransactionCategoriesToolHandler)
	registerMCPTextContentToolHandler(container, MCPQueryAllTransactionTagsToolHandler)
	registerMCPTextContentToolHandler(container, MCPQueryLatestExchangeRatesToolHandler)

	Container = container
	return nil
}

func registerMCPTextContentToolHandler(c *MCPContainer, handler MCPToolHandler[MCPTextContent]) {
	registerMCPToolHandler(c, c.mcpTextContentTools, handler)
}

func registerMCPImageContentToolHandler(c *MCPContainer, handler MCPToolHandler[MCPImageContent]) {
	registerMCPToolHandler(c, c.mcpImageContentTools, handler)
}

func registerMCPAudioContentToolHandler(c *MCPContainer, handler MCPToolHandler[MCPAudioContent]) {
	registerMCPToolHandler(c, c.mcpAudioContentTools, handler)
}

func registerMCPResourceLinkToolHandler(c *MCPContainer, handler MCPToolHandler[MCPResourceLink]) {
	registerMCPToolHandler(c, c.mcpResourceLinkTools, handler)
}

func registerMCPEmbeddedResourceToolHandler(c *MCPContainer, handler MCPToolHandler[MCPEmbeddedResource]) {
	registerMCPToolHandler(c, c.mcpEmbeddedResourceTools, handler)
}

func registerMCPToolHandler[T MCPTextContent | MCPImageContent | MCPAudioContent | MCPResourceLink | MCPEmbeddedResource](c *MCPContainer, mcpToolHandlerMap *orderedmap.OrderedMap[string, MCPToolHandler[T]], handler MCPToolHandler[T]) {
	if _, exists := mcpToolHandlerMap.Get(handler.Name()); exists {
		return
	}

	mcpToolHandlerMap.Set(handler.Name(), handler)
	c.mcpTools = append(c.mcpTools, createNewMCPToolInfo(handler.Name(), handler))
}

func handleTool[T MCPTextContent | MCPImageContent | MCPAudioContent | MCPResourceLink | MCPEmbeddedResource](ctx *core.WebContext, handler MCPToolHandler[T], currentConfig *settings.Config, services MCPAvailableServices, callToolReq *MCPCallToolRequest) (any, *errs.Error) {
	structuredResponse, result, err := handler.Handle(ctx, callToolReq, currentConfig, services)

	if err != nil {
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	callToolResp := MCPCallToolResponse[T]{
		Content: result,
		IsError: false,
	}

	if ctx.GetHeader(MCPProtocolVersionHeaderName) > string(ToolResultStructuredContentMinVersion) {
		callToolResp.StructuredContent = structuredResponse
	}

	return callToolResp, nil
}

func createNewMCPToolInfo[T MCPTextContent | MCPImageContent | MCPAudioContent | MCPResourceLink | MCPEmbeddedResource](name string, handler MCPToolHandler[T]) *MCPTool {
	mcpTool := &MCPTool{
		Name:        name,
		Description: handler.Description(),
	}

	schemeGenerator := jsonschema.Reflector{
		Anonymous:      true,
		DoNotReference: true,
		ExpandedStruct: true,
	}

	if handler.InputType() != nil {
		schema := schemeGenerator.ReflectFromType(handler.InputType())
		schema.Version = ""
		mcpTool.InputSchema = schema
	} else {
		mcpTool.InputSchema = &jsonschema.Schema{
			Type: "object",
		}
	}

	if handler.OutputType() != nil {
		schema := schemeGenerator.ReflectFromType(handler.OutputType())
		schema.Version = ""
		mcpTool.OutputSchema = schema
	} else {
		mcpTool.OutputSchema = &jsonschema.Schema{
			Type: "object",
		}
	}

	return mcpTool
}
