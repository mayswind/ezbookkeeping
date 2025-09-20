package llm

import (
	"net/http"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// OpenAILargeLanguageModelProvider defines the structure of OpenAI large language model provider
type OpenAILargeLanguageModelProvider struct {
	OpenAIChatCompletionsLargeLanguageModelProvider
	OpenAIAPIKey                   string
	ReceiptImageRecognitionModelID string
}

const openAIChatCompletionsUrl = "https://api.openai.com/v1/chat/completions"

// BuildChatCompletionsHttpRequest returns the chat completions http request by OpenAI provider
func (p *OpenAILargeLanguageModelProvider) BuildChatCompletionsHttpRequest(c core.Context, uid int64) (*http.Request, error) {
	req, err := http.NewRequest("POST", openAIChatCompletionsUrl, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+p.OpenAIAPIKey)

	return req, nil
}

// GetReceiptImageRecognitionModelID returns the receipt image recognition model id of OpenAI provider
func (p *OpenAILargeLanguageModelProvider) GetReceiptImageRecognitionModelID() string {
	return p.ReceiptImageRecognitionModelID
}

// NewOpenAILargeLanguageModelProvider creates a new OpenAI large language model provider instance
func NewOpenAILargeLanguageModelProvider(config *settings.Config) LargeLanguageModelProvider {
	return newOpenAICommonChatCompletionsHttpLargeLanguageModelProvider(&OpenAILargeLanguageModelProvider{
		OpenAIAPIKey:                   config.OpenAIAPIKey,
		ReceiptImageRecognitionModelID: config.OpenAIReceiptImageRecognitionModelID,
	})
}
