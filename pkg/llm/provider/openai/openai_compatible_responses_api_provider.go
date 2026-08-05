package openai

import (
	"net/http"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/llm/provider"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

const openAICompatibleResponsesPath = "responses"

// OpenAICompatibleResponsesAPIProvider defines the structure of OpenAI compatible responses API provider
type OpenAICompatibleResponsesAPIProvider struct {
	OpenAIResponsesAPIProvider
	OpenAICompatibleBaseURL string
	OpenAICompatibleAPIKey  string
	OpenAICompatibleModelID string
}

// BuildResponsesHttpRequest returns the responses http request by OpenAI compatible responses API provider
func (p *OpenAICompatibleResponsesAPIProvider) BuildResponsesHttpRequest(c core.Context, uid int64) (*http.Request, error) {
	req, err := http.NewRequest("POST", p.getFinalResponsesRequestUrl(), nil)

	if err != nil {
		return nil, err
	}

	if p.OpenAICompatibleAPIKey != "" {
		req.Header.Set("Authorization", "Bearer "+p.OpenAICompatibleAPIKey)
	}

	return req, nil
}

// GetModelID returns the model id of OpenAI compatible responses API provider
func (p *OpenAICompatibleResponsesAPIProvider) GetModelID() string {
	return p.OpenAICompatibleModelID
}

func (p *OpenAICompatibleResponsesAPIProvider) getFinalResponsesRequestUrl() string {
	url := p.OpenAICompatibleBaseURL

	if url[len(url)-1] != '/' {
		url += "/"
	}

	url += openAICompatibleResponsesPath
	return url
}

// NewOpenAIResponsesCompatibleLargeLanguageModelProvider creates a new OpenAI responses compatible large language model provider instance
func NewOpenAIResponsesCompatibleLargeLanguageModelProvider(llmConfig *settings.LLMConfig, enableResponseLog bool) provider.LargeLanguageModelProvider {
	return newCommonOpenAIResponsesAPILargeLanguageModelAdapter(llmConfig, enableResponseLog, &OpenAICompatibleResponsesAPIProvider{
		OpenAICompatibleBaseURL: llmConfig.OpenAICompatibleBaseURL,
		OpenAICompatibleAPIKey:  llmConfig.OpenAICompatibleAPIKey,
		OpenAICompatibleModelID: llmConfig.OpenAICompatibleModelID,
	})
}
