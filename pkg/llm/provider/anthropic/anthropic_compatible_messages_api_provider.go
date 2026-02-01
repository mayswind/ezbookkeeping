package anthropic

import (
	"net/http"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/llm/provider"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

const anthropicCompatibleMessagesPath = "messages"

// AnthropicCompatibleMessagesAPIProvider defines the structure of Anthropic compatible messages API provider
type AnthropicCompatibleMessagesAPIProvider struct {
	AnthropicMessagesAPIProvider
	AnthropicCompatibleBaseURL    string
	AnthropicCompatibleAPIVersion string
	AnthropicCompatibleAPIKey     string
	AnthropicCompatibleModelID    string
	AnthropicCompatibleMaxTokens  uint32
}

// BuildMessagesHttpRequest returns the messages http request by Anthropic compatible messages API provider
func (p *AnthropicCompatibleMessagesAPIProvider) BuildMessagesHttpRequest(c core.Context, uid int64) (*http.Request, error) {
	req, err := http.NewRequest("POST", p.getFinalMessagesRequestUrl(), nil)

	if err != nil {
		return nil, err
	}

	if p.AnthropicCompatibleAPIVersion != "" {
		req.Header.Set("anthropic-version", p.AnthropicCompatibleAPIVersion)
	}

	if p.AnthropicCompatibleAPIKey != "" {
		req.Header.Set("X-Api-Key", p.AnthropicCompatibleAPIKey)
	}

	return req, nil
}

// GetModelID returns the model id of Anthropic compatible messages API provider
func (p *AnthropicCompatibleMessagesAPIProvider) GetModelID() string {
	return p.AnthropicCompatibleModelID
}

// GetMaxTokens returns the max tokens to generate of Anthropic compatible messages API provider
func (p *AnthropicCompatibleMessagesAPIProvider) GetMaxTokens() uint32 {
	return p.AnthropicCompatibleMaxTokens
}

func (p *AnthropicCompatibleMessagesAPIProvider) getFinalMessagesRequestUrl() string {
	url := p.AnthropicCompatibleBaseURL

	if url[len(url)-1] != '/' {
		url += "/"
	}

	url += anthropicCompatibleMessagesPath
	return url
}

// NewAnthropicCompatibleLargeLanguageModelProvider creates a new Anthropic compatible large language model provider instance
func NewAnthropicCompatibleLargeLanguageModelProvider(llmConfig *settings.LLMConfig, enableResponseLog bool) provider.LargeLanguageModelProvider {
	return newCommonAnthropicMessagesAPILargeLanguageModelAdapter(llmConfig, enableResponseLog, &AnthropicCompatibleMessagesAPIProvider{
		AnthropicCompatibleBaseURL:    llmConfig.AnthropicCompatibleBaseURL,
		AnthropicCompatibleAPIVersion: llmConfig.AnthropicCompatibleAPIVersion,
		AnthropicCompatibleAPIKey:     llmConfig.AnthropicCompatibleAPIKey,
		AnthropicCompatibleModelID:    llmConfig.AnthropicCompatibleModelID,
		AnthropicCompatibleMaxTokens:  llmConfig.AnthropicCompatibleMaxTokens,
	})
}
