package llm

import (
	"net/http"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

const openAICompatibleChatCompletionsPath = "chat/completions"

// OpenAICompatibleChatCompletionsAPIProvider defines the structure of OpenAI compatible chat completions API provider
type OpenAICompatibleChatCompletionsAPIProvider struct {
	OpenAIChatCompletionsAPIProvider
	OpenAICompatibleBaseURL string
	OpenAICompatibleAPIKey  string
	OpenAICompatibleModelID string
}

// BuildChatCompletionsHttpRequest returns the chat completions http request by OpenAI compatible chat completions API provider
func (p *OpenAICompatibleChatCompletionsAPIProvider) BuildChatCompletionsHttpRequest(c core.Context, uid int64) (*http.Request, error) {
	req, err := http.NewRequest("POST", p.getFinalChatCompletionsRequestUrl(), nil)

	if err != nil {
		return nil, err
	}

	if p.OpenAICompatibleAPIKey != "" {
		req.Header.Set("Authorization", "Bearer "+p.OpenAICompatibleAPIKey)
	}

	return req, nil
}

// GetModelID returns the model id of OpenAI compatible chat completions API provider
func (p *OpenAICompatibleChatCompletionsAPIProvider) GetModelID() string {
	return p.OpenAICompatibleModelID
}

func (p *OpenAICompatibleChatCompletionsAPIProvider) getFinalChatCompletionsRequestUrl() string {
	url := p.OpenAICompatibleBaseURL

	if url[len(url)-1] != '/' {
		url += "/"
	}

	url += openAICompatibleChatCompletionsPath
	return url
}

// NewOpenAICompatibleLargeLanguageModelProvider creates a new OpenAI compatible large language model provider instance
func NewOpenAICompatibleLargeLanguageModelProvider(llmConfig *settings.LLMConfig) LargeLanguageModelProvider {
	return newCommonOpenAIChatCompletionsAPILargeLanguageModelAdapter(&OpenAICompatibleChatCompletionsAPIProvider{
		OpenAICompatibleBaseURL: llmConfig.OpenAICompatibleBaseURL,
		OpenAICompatibleAPIKey:  llmConfig.OpenAICompatibleAPIKey,
		OpenAICompatibleModelID: llmConfig.OpenAICompatibleModelID,
	})
}
