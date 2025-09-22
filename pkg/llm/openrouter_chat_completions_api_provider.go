package llm

import (
	"net/http"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// OpenRouterChatCompletionsAPIProvider defines the structure of OpenRouter chat completions API provider
type OpenRouterChatCompletionsAPIProvider struct {
	OpenAIChatCompletionsAPIProvider
	OpenRouterAPIKey  string
	OpenRouterModelID string
}

const openRouterChatCompletionsUrl = "https://openrouter.ai/api/v1/chat/completions"

// BuildChatCompletionsHttpRequest returns the chat completions http request by OpenRouter chat completions API provider
func (p *OpenRouterChatCompletionsAPIProvider) BuildChatCompletionsHttpRequest(c core.Context, uid int64) (*http.Request, error) {
	req, err := http.NewRequest("POST", openRouterChatCompletionsUrl, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+p.OpenRouterAPIKey)
	req.Header.Set("HTTP-Referer", "https://ezbookkeeping.mayswind.net/")
	req.Header.Set("X-Title", "ezBookkeeping")

	return req, nil
}

// GetModelID returns the model id of OpenRouter chat completions API provider
func (p *OpenRouterChatCompletionsAPIProvider) GetModelID() string {
	return p.OpenRouterModelID
}

// NewOpenRouterLargeLanguageModelProvider creates a new OpenRouter large language model provider instance
func NewOpenRouterLargeLanguageModelProvider(llmConfig *settings.LLMConfig) LargeLanguageModelProvider {
	return newCommonOpenAIChatCompletionsAPILargeLanguageModelAdapter(&OpenRouterChatCompletionsAPIProvider{
		OpenRouterAPIKey:  llmConfig.OpenRouterAPIKey,
		OpenRouterModelID: llmConfig.OpenRouterModelID,
	})
}
