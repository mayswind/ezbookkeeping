package openai

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenAICompatibleChatCompletionsAPIProvider_GetFinalRequestUrl(t *testing.T) {
	apiProvider := &OpenAICompatibleChatCompletionsAPIProvider{
		OpenAICompatibleBaseURL: "https://api.example.com/v1/",
	}
	url := apiProvider.getFinalChatCompletionsRequestUrl()
	assert.Equal(t, "https://api.example.com/v1/chat/completions", url)

	apiProvider = &OpenAICompatibleChatCompletionsAPIProvider{
		OpenAICompatibleBaseURL: "https://api.example.com/v1",
	}
	url = apiProvider.getFinalChatCompletionsRequestUrl()
	assert.Equal(t, "https://api.example.com/v1/chat/completions", url)

	apiProvider = &OpenAICompatibleChatCompletionsAPIProvider{
		OpenAICompatibleBaseURL: "https://example.com/api",
	}
	url = apiProvider.getFinalChatCompletionsRequestUrl()
	assert.Equal(t, "https://example.com/api/chat/completions", url)
}
