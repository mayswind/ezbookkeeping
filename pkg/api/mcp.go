package api

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/mcp"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// MCPApi represents MCP api
type MCPApi struct {
	ApiUsingConfig
	mcpServer *mcp.MCPServer
}

// Initialize an MCP api singleton instance
var (
	MCP = &MCPApi{
		ApiUsingConfig: ApiUsingConfig{
			container: settings.Container,
		},
		mcpServer: mcp.NewMCPServer(),
	}
)

// MCPHandler handles MCP requests
func (a *MCPApi) MCPHandler(c *core.WebContext) (interface{}, *errs.Error) {
	uid := c.GetCurrentUid()
	if uid <= 0 {
		return nil, errs.ErrUserNotFound
	}

	config := a.CurrentConfig()
	if !config.EnableMCP {
		return nil, errs.ErrMCPNotEnabled
	}

	response, err := a.mcpServer.HandleRequest(c)
	if err != nil {
		log.Warnf(c, "Failed to handle MCP request: %s", err.Error())
		return nil, err
	}

	return response, nil
}