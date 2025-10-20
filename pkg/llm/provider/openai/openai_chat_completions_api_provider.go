package openai

import (
	"net/http"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/llm/provider"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// OpenAIOfficialChatCompletionsAPIProvider defines the structure of OpenAI official chat completions API provider
type OpenAIOfficialChatCompletionsAPIProvider struct {
	OpenAIChatCompletionsAPIProvider
	OpenAIAPIKey  string
	OpenAIModelID string
}

const openAIChatCompletionsUrl = "https://api.openai.com/v1/chat/completions"

// BuildChatCompletionsHttpRequest returns the chat completions http request by OpenAI official chat completions API provider
func (p *OpenAIOfficialChatCompletionsAPIProvider) BuildChatCompletionsHttpRequest(c core.Context, uid int64) (*http.Request, error) {
	req, err := http.NewRequest("POST", openAIChatCompletionsUrl, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+p.OpenAIAPIKey)

	return req, nil
}

// GetModelID returns the model id of OpenAI official chat completions API provider
func (p *OpenAIOfficialChatCompletionsAPIProvider) GetModelID() string {
	return p.OpenAIModelID
}

// NewOpenAILargeLanguageModelProvider creates a new OpenAI large language model provider instance
func NewOpenAILargeLanguageModelProvider(llmConfig *settings.LLMConfig) provider.LargeLanguageModelProvider {
	return newCommonOpenAIChatCompletionsAPILargeLanguageModelAdapter(llmConfig, &OpenAIOfficialChatCompletionsAPIProvider{
		OpenAIAPIKey:  llmConfig.OpenAIAPIKey,
		OpenAIModelID: llmConfig.OpenAIModelID,
	})
}
