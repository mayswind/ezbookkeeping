package openai

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenAICompatibleResponsesAPIProvider_GetFinalRequestUrl(t *testing.T) {
	apiProvider := &OpenAICompatibleResponsesAPIProvider{
		OpenAICompatibleBaseURL: "https://api.example.com/v1/",
	}
	url := apiProvider.getFinalResponsesRequestUrl()
	assert.Equal(t, "https://api.example.com/v1/responses", url)

	apiProvider = &OpenAICompatibleResponsesAPIProvider{
		OpenAICompatibleBaseURL: "https://api.example.com/v1",
	}
	url = apiProvider.getFinalResponsesRequestUrl()
	assert.Equal(t, "https://api.example.com/v1/responses", url)

	apiProvider = &OpenAICompatibleResponsesAPIProvider{
		OpenAICompatibleBaseURL: "https://example.com/api",
	}
	url = apiProvider.getFinalResponsesRequestUrl()
	assert.Equal(t, "https://example.com/api/responses", url)
}
