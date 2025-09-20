package llm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenAICompatibleLargeLanguageModelProvider_GetFinalRequestUrl(t *testing.T) {
	provider := &OpenAICompatibleLargeLanguageModelProvider{
		OpenAICompatibleBaseURL: "https://api.example.com/v1/",
	}
	url := provider.getFinalChatCompletionsRequestUrl()
	assert.Equal(t, "https://api.example.com/v1/chat/completions", url)

	provider = &OpenAICompatibleLargeLanguageModelProvider{
		OpenAICompatibleBaseURL: "https://api.example.com/v1",
	}
	url = provider.getFinalChatCompletionsRequestUrl()
	assert.Equal(t, "https://api.example.com/v1/chat/completions", url)

	provider = &OpenAICompatibleLargeLanguageModelProvider{
		OpenAICompatibleBaseURL: "https://example.com/api",
	}
	url = provider.getFinalChatCompletionsRequestUrl()
	assert.Equal(t, "https://example.com/api/chat/completions", url)
}
