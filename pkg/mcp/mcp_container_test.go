package mcp

import (
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

func initializeMCPContainerTest(t *testing.T) {
	t.Helper()
	previousContainer := Container
	t.Cleanup(func() {
		Container = previousContainer
	})

	err := InitializeMCPHandlers(&settings.Config{})

	if err != nil {
		t.Fatal(err)
	}
}

func TestMCPContainerGetMCPTools(t *testing.T) {
	initializeMCPContainerTest(t)
	var allTools int64

	for _, tool := range Container.GetMCPTools() {
		bit := int64(1) << (tool.ToolType - 1)
		assert.Equal(t, int64(0), allTools&bit)
		allTools |= bit
	}

	assert.Equal(t, core.MCPAllTools, allTools)
}

func TestMCPContainerHandleTool_DisabledTool(t *testing.T) {
	initializeMCPContainerTest(t)

	context := &core.WebContext{
		Context: &gin.Context{},
	}

	for _, tool := range Container.GetMCPTools() {
		request := &MCPCallToolRequest{
			Name: tool.Name,
		}

		context.SetTokenRecord(&models.TokenRecord{
			Context: `{"enabledTools":"0"}`,
		})
		actualValue, err := Container.HandleTool(context, request, nil, nil, nil)
		assert.Nil(t, actualValue)
		assert.Equal(t, errs.ErrMCPToolNotEnabled, err)

		bit := int64(1) << (tool.ToolType - 1)
		context.SetTokenRecord(&models.TokenRecord{
			Context: `{"enabledTools":"` + strconv.FormatInt(core.MCPAllTools & ^bit, 10) + `"}`,
		})
		actualValue, err = Container.HandleTool(context, request, nil, nil, nil)
		assert.Nil(t, actualValue)
		assert.Equal(t, errs.ErrMCPToolNotEnabled, err)
	}
}

func TestMCPContainerHandleTool_InvalidContext(t *testing.T) {
	initializeMCPContainerTest(t)

	context := &core.WebContext{
		Context: &gin.Context{},
	}

	for _, tool := range Container.GetMCPTools() {
		request := &MCPCallToolRequest{
			Name: tool.Name,
		}

		context.SetTokenRecord(&models.TokenRecord{
			Context: "invalid",
		})
		actualValue, err := Container.HandleTool(context, request, nil, nil, nil)
		assert.Nil(t, actualValue)
		assert.Equal(t, errs.ErrInvalidToken, err)

		context.SetTokenRecord(&models.TokenRecord{
			Context: `{"enabledTools":"-1"}`,
		})
		actualValue, err = Container.HandleTool(context, request, nil, nil, nil)
		assert.Nil(t, actualValue)
		assert.Equal(t, errs.ErrMCPToolNotEnabled, err)
	}
}
