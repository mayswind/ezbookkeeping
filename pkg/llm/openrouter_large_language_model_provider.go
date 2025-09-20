package llm

import (
	"net/http"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// OpenRouterLargeLanguageModelProvider defines the structure of OpenRouter large language model provider
type OpenRouterLargeLanguageModelProvider struct {
	OpenAIChatCompletionsLargeLanguageModelProvider
	OpenRouterAPIKey               string
	ReceiptImageRecognitionModelID string
}

const openRouterChatCompletionsUrl = "https://openrouter.ai/api/v1/chat/completions"

// BuildChatCompletionsHttpRequest returns the chat completions http request by OpenRouter provider
func (p *OpenRouterLargeLanguageModelProvider) BuildChatCompletionsHttpRequest(c core.Context, uid int64) (*http.Request, error) {
	req, err := http.NewRequest("POST", openRouterChatCompletionsUrl, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+p.OpenRouterAPIKey)
	req.Header.Set("HTTP-Referer", "https://ezbookkeeping.mayswind.net/")
	req.Header.Set("X-Title", "ezBookkeeping")

	return req, nil
}

// GetReceiptImageRecognitionModelID returns the receipt image recognition model id of OpenRouter provider
func (p *OpenRouterLargeLanguageModelProvider) GetReceiptImageRecognitionModelID() string {
	return p.ReceiptImageRecognitionModelID
}

// NewOpenRouterLargeLanguageModelProvider creates a new OpenRouter large language model provider instance
func NewOpenRouterLargeLanguageModelProvider(config *settings.Config) LargeLanguageModelProvider {
	return newOpenAICommonChatCompletionsHttpLargeLanguageModelProvider(&OpenRouterLargeLanguageModelProvider{
		OpenRouterAPIKey:               config.OpenRouterAPIKey,
		ReceiptImageRecognitionModelID: config.OpenRouterReceiptImageRecognitionModelID,
	})
}
