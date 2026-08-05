package openai

import (
	"net/http"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/llm/provider"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// OpenAIOfficialResponsesAPIProvider defines the structure of OpenAI official responses API provider
type OpenAIOfficialResponsesAPIProvider struct {
	OpenAIResponsesAPIProvider
	OpenAIAPIKey  string
	OpenAIModelID string
}

const openAIResponsesUrl = "https://api.openai.com/v1/responses"

// BuildResponsesHttpRequest returns the responses http request by OpenAI official responses API provider
func (p *OpenAIOfficialResponsesAPIProvider) BuildResponsesHttpRequest(c core.Context, uid int64) (*http.Request, error) {
	req, err := http.NewRequest("POST", openAIResponsesUrl, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+p.OpenAIAPIKey)

	return req, nil
}

// GetModelID returns the model id of OpenAI official responses API provider
func (p *OpenAIOfficialResponsesAPIProvider) GetModelID() string {
	return p.OpenAIModelID
}

// NewOpenAILargeLanguageModelProvider creates a new OpenAI large language model provider instance
func NewOpenAILargeLanguageModelProvider(llmConfig *settings.LLMConfig, enableResponseLog bool) provider.LargeLanguageModelProvider {
	return newCommonOpenAIResponsesAPILargeLanguageModelAdapter(llmConfig, enableResponseLog, &OpenAIOfficialResponsesAPIProvider{
		OpenAIAPIKey:  llmConfig.OpenAIAPIKey,
		OpenAIModelID: llmConfig.OpenAIModelID,
	})
}
