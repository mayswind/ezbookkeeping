package mcp

import (
	"encoding/json"
	"reflect"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// MCPAllQueryTransactionTagsResponse represents the response structure for querying transaction tags
type MCPAllQueryTransactionTagsResponse struct {
	Tags []string `json:"tags" jsonschema_description:"List of transaction tags"`
}

type mcpQueryAllTransactionTagsToolHandler struct{}

var MCPQueryAllTransactionTagsToolHandler = &mcpQueryAllTransactionTagsToolHandler{}

// Name returns the name of the MCP tool
func (h *mcpQueryAllTransactionTagsToolHandler) Name() string {
	return "query_all_transaction_tags"
}

// Description returns the description of the MCP tool
func (h *mcpQueryAllTransactionTagsToolHandler) Description() string {
	return "Query transaction tags for the current user in ezBookkeeping."
}

// InputType returns the input type for the MCP tool request
func (h *mcpQueryAllTransactionTagsToolHandler) InputType() reflect.Type {
	return nil
}

// OutputType returns the output type for the MCP tool response
func (h *mcpQueryAllTransactionTagsToolHandler) OutputType() reflect.Type {
	return reflect.TypeOf(&MCPAllQueryTransactionTagsResponse{})
}

// Handle processes the MCP call tool request and returns the response
func (h *mcpQueryAllTransactionTagsToolHandler) Handle(c *core.WebContext, callToolReq *MCPCallToolRequest, user *models.User, currentConfig *settings.Config, services MCPAvailableServices) (any, []*MCPTextContent, error) {
	uid := user.Uid
	tags, err := services.GetTransactionTagService().GetAllTagsByUid(c, uid)

	if err != nil {
		log.Errorf(c, "[query_all_transaction_tags.Handle] failed to get tags for user \"uid:%d\", because %s", uid, err.Error())
		return nil, nil, err
	}

	tagNames := make([]string, 0, len(tags))

	for i := 0; i < len(tags); i++ {
		if tags[i].Hidden {
			continue
		}

		tagNames = append(tagNames, tags[i].Name)
	}

	response := MCPAllQueryTransactionTagsResponse{
		Tags: tagNames,
	}

	content, err := json.Marshal(response)

	if err != nil {
		return nil, nil, err
	}

	return response, []*MCPTextContent{
		NewMCPTextContent(string(content)),
	}, nil
}
