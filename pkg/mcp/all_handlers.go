package mcp

var mcpTextContentTools = map[string]MCPToolHandler[MCPTextContent]{
	"query_latest_exchange_rates": MCPQueryLatestExchangeRatesRequestToolHandler,
}

var mcpImageContentTools = map[string]MCPToolHandler[MCPImageContent]{}
var mcpAudioContentTools = map[string]MCPToolHandler[MCPAudioContent]{}
var mcpResourceLinkTools = map[string]MCPToolHandler[MCPResourceLink]{}
var mcpEmbeddedResourceTools = map[string]MCPToolHandler[MCPEmbeddedResource]{}

var AllMCPToolInfos = GetAllMCPToolInfos()
