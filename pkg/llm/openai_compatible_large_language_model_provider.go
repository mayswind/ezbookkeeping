package llm

import (
	"net/http"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

const openAICompatibleChatCompletionsPath = "chat/completions"

// OpenAICompatibleLargeLanguageModelProvider defines the structure of OpenAI compatible large language model provider
type OpenAICompatibleLargeLanguageModelProvider struct {
	OpenAIChatCompletionsLargeLanguageModelProvider
	OpenAICompatibleBaseURL        string
	OpenAICompatibleAPIKey         string
	ReceiptImageRecognitionModelID string
}

// BuildChatCompletionsHttpRequest returns the chat completions http request by OpenAI compatible provider
func (p *OpenAICompatibleLargeLanguageModelProvider) BuildChatCompletionsHttpRequest(c core.Context, uid int64) (*http.Request, error) {
	req, err := http.NewRequest("POST", p.getFinalChatCompletionsRequestUrl(), nil)

	if err != nil {
		return nil, err
	}

	if p.OpenAICompatibleAPIKey != "" {
		req.Header.Set("Authorization", "Bearer "+p.OpenAICompatibleAPIKey)
	}

	return req, nil
}

// GetReceiptImageRecognitionModelID returns the receipt image recognition model id of OpenAI compatible provider
func (p *OpenAICompatibleLargeLanguageModelProvider) GetReceiptImageRecognitionModelID() string {
	return p.ReceiptImageRecognitionModelID
}

func (p *OpenAICompatibleLargeLanguageModelProvider) getFinalChatCompletionsRequestUrl() string {
	url := p.OpenAICompatibleBaseURL

	if url[len(url)-1] != '/' {
		url += "/"
	}

	url += openAICompatibleChatCompletionsPath
	return url
}

// NewOpenAICompatibleLargeLanguageModelProvider creates a new OpenAI compatible large language model provider instance
func NewOpenAICompatibleLargeLanguageModelProvider(config *settings.Config) LargeLanguageModelProvider {
	return newOpenAICommonChatCompletionsHttpLargeLanguageModelProvider(&OpenAICompatibleLargeLanguageModelProvider{
		OpenAICompatibleBaseURL:        config.OpenAICompatibleBaseURL,
		OpenAICompatibleAPIKey:         config.OpenAICompatibleAPIKey,
		ReceiptImageRecognitionModelID: config.OpenAICompatibleReceiptImageRecognitionModelID,
	})
}
