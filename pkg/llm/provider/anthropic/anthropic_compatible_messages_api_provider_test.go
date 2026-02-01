package anthropic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnthropicCompatibleMessagesAPIProvider_GetFinalRequestUrl(t *testing.T) {
	apiProvider := &AnthropicCompatibleMessagesAPIProvider{
		AnthropicCompatibleBaseURL: "https://api.example.com/v1/",
	}
	url := apiProvider.getFinalMessagesRequestUrl()
	assert.Equal(t, "https://api.example.com/v1/messages", url)

	apiProvider = &AnthropicCompatibleMessagesAPIProvider{
		AnthropicCompatibleBaseURL: "https://api.example.com/v1",
	}
	url = apiProvider.getFinalMessagesRequestUrl()
	assert.Equal(t, "https://api.example.com/v1/messages", url)

	apiProvider = &AnthropicCompatibleMessagesAPIProvider{
		AnthropicCompatibleBaseURL: "https://example.com/api",
	}
	url = apiProvider.getFinalMessagesRequestUrl()
	assert.Equal(t, "https://example.com/api/messages", url)
}
