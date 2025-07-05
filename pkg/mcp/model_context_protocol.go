package mcp

import (
	"encoding/base64"
	"encoding/json"

	"github.com/invopop/jsonschema"
)

// MCPProtocolVersion defines the type for Model Context Protocol (MCP) version
type MCPProtocolVersion string

// MCP Protocol Versions
const (
	MCPProtocolVersion20250618 MCPProtocolVersion = "2025-06-18"
	MCPProtocolVersion20250326 MCPProtocolVersion = "2025-03-26"
	MCPProtocolVersion20241105 MCPProtocolVersion = "2024-11-05"
)

// LatestSupportedMCPVersion defines the latest supported version of Model Context Protocol (MCP)
const LatestSupportedMCPVersion = MCPProtocolVersion20250618

// SupportedMCPVersion defines a map of supported MCP versions
var SupportedMCPVersion = map[MCPProtocolVersion]bool{
	MCPProtocolVersion20250618: true,
	MCPProtocolVersion20250326: true,
	MCPProtocolVersion20241105: true,
}

// MCPInitializeRequest defines the request structure for initializing the MCP connection
type MCPInitializeRequest struct {
	ProtocolVersion string             `json:"protocolVersion"`
	ClientInfo      *MCPImplementation `json:"clientInfo"`
}

// MCPInitializeResponse defines the response structure for the MCP initialization request
type MCPInitializeResponse struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    *MCPCapabilities   `json:"capabilities"`
	ServerInfo      *MCPImplementation `json:"serverInfo"`
}

// MCPCapabilities defines the capabilities of the MCP server
type MCPCapabilities struct {
	Resources *MCPResourceCapabilities `json:"resources,omitempty"`
	Tools     *MCPToolCapabilities     `json:"tools,omitempty"`
	Prompts   *MCPPromptCapabilities   `json:"prompts,omitempty"`
}

// MCPImplementation defines the client/server information structure sent in the MCP initialization request/response
type MCPImplementation struct {
	Name    string `json:"name"`
	Title   string `json:"title,omitempty"`
	Version string `json:"version"`
}

// MCPResourceCapabilities defines the capabilities related to resources in the MCP
type MCPResourceCapabilities struct {
	Subscribe   bool `json:"subscribe"`
	ListChanged bool `json:"listChanged"`
}

// MCPToolCapabilities defines the capabilities related to tools in the MCP
type MCPToolCapabilities struct {
	ListChanged bool `json:"listChanged"`
}

// MCPPromptCapabilities defines the capabilities related to prompts in the MCP
type MCPPromptCapabilities struct {
	ListChanged bool `json:"listChanged"`
}

// MCPListResourcesResponse defines the response structure for listing resources in the MCP
type MCPListResourcesResponse struct {
	Resources  []*MCPResource `json:"resources"`
	NextCursor string         `json:"nextCursor,omitempty"`
}

// MCPResource defines the structure of a resource in the MCP
type MCPResource struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Size        int    `json:"size,omitempty"`
	MimeType    string `json:"mimeType,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

// MCPReadResourceRequest defines the request structure for reading a resource in the MCP
type MCPReadResourceRequest struct {
	URI string `json:"uri"`
}

// MCPReadResourceResponse defines the response structure for reading a resource in the MCP
type MCPReadResourceResponse[T MCPTextResourceContents | MCPBlobResourceContents] struct {
	Contents []*T `json:"contents"`
}

// MCPTextResourceContents defines the text contents structure of a resource in the MCP
type MCPTextResourceContents struct {
	URI      string `json:"uri"`
	Text     string `json:"text"`
	MimeType string `json:"mimeType,omitempty"`
}

// MCPBlobResourceContents defines the blob contents structure of a resource in the MCP
type MCPBlobResourceContents struct {
	URI      string `json:"uri"`
	Blob     string `json:"blob"` // Base64 encoded content of the resource
	MimeType string `json:"mimeType,omitempty"`
}

// MCPListToolsResponse defines the response structure for listing tools in the MCP
type MCPListToolsResponse struct {
	Tools      []*MCPTool `json:"tools"`
	NextCursor string     `json:"nextCursor,omitempty"`
}

// MCPTool defines the structure of a tool in the MCP
type MCPTool struct {
	Name         string             `json:"name"`
	InputSchema  *jsonschema.Schema `json:"inputSchema"`
	OutputSchema *jsonschema.Schema `json:"outputSchema"`
	Title        string             `json:"title,omitempty"`
	Description  string             `json:"description,omitempty"`
}

// MCPCallToolRequest defines the request structure for listing tools in the MCP
type MCPCallToolRequest struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments,omitempty"`
}

// MCPCallToolResponse defines the response structure for calling a tool in the MCP
type MCPCallToolResponse[T MCPTextContent | MCPImageContent | MCPAudioContent | MCPResourceLink | MCPEmbeddedResource] struct {
	Content []*T `json:"content"`
	IsError bool `json:"isError,omitempty"`
}

// MCPTextContent defines the text content structure used in MCP
type MCPTextContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// MCPImageContent defines the image content structure used in MCP
type MCPImageContent struct {
	Type     string `json:"type"`
	MimeType string `json:"mimeType"`
	Data     string `json:"data"` // Base64 encoded content for binary data
}

// MCPAudioContent defines the audio content structure used in MCP
type MCPAudioContent struct {
	Type     string `json:"type"`
	MimeType string `json:"mimeType"`
	Data     string `json:"data"` // Base64 encoded content for binary data
}

// MCPResourceLink defines the resource link content structure used in MCP
type MCPResourceLink struct {
	URI         string `json:"uri"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Size        int    `json:"size,omitempty"`
	MimeType    string `json:"mimeType,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

// MCPEmbeddedResource defines the embedded resource content structure used in MCP
type MCPEmbeddedResource struct {
	Type     string `json:"type"`
	Resource any    `json:"resource"`
}

// NewMCPTextContent creates a new instance of MCPTextContent with the given text
func NewMCPTextContent(text string) *MCPTextContent {
	return &MCPTextContent{
		Type: "text",
		Text: text,
	}
}

// NewMCPImageContent creates a new instance of MCPImageContent with the given data and MIME type
func NewMCPImageContent(data []byte, mimeType string) *MCPImageContent {
	return &MCPImageContent{
		Type:     "image",
		MimeType: mimeType,
		Data:     base64.StdEncoding.EncodeToString(data),
	}
}

// NewMCPAudioContent creates a new instance of MCPAudioContent with the given data and MIME type
func NewMCPAudioContent(data []byte, mimeType string) *MCPAudioContent {
	return &MCPAudioContent{
		Type:     "audio",
		MimeType: mimeType,
		Data:     base64.StdEncoding.EncodeToString(data),
	}
}

// NewMCPResourceLink creates a new instance of MCPResourceLink with the given parameters
func NewMCPResourceLink(uri string, name string) *MCPResourceLink {
	return &MCPResourceLink{
		URI:  uri,
		Type: "resource_link",
		Name: name,
	}
}

// NewMCPEmbeddedResource creates a new instance of MCPEmbeddedResource with the given resource
func NewMCPEmbeddedResource[T MCPTextResourceContents | MCPBlobResourceContents](resource *T) *MCPEmbeddedResource {
	return &MCPEmbeddedResource{
		Type:     "resource",
		Resource: resource,
	}
}
