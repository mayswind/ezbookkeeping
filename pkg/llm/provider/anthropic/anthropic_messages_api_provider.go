package anthropic

import (
	"net/http"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/llm/provider"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// AnthropicOfficialMessagesAPIProvider defines the structure of Anthropic official messages API provider
type AnthropicOfficialMessagesAPIProvider struct {
	AnthropicMessagesAPIProvider
	AnthropicAPIKey    string
	AnthropicModelID   string
	AnthropicMaxTokens uint32
}

const anthropicMessagesUrl = "https://api.anthropic.com/v1/messages"
const anthropicAPIVersion = "2023-06-01"

// BuildMessagesHttpRequest returns the messages http request by Anthropic official messages API provider
func (p *AnthropicOfficialMessagesAPIProvider) BuildMessagesHttpRequest(c core.Context, uid int64) (*http.Request, error) {
	req, err := http.NewRequest("POST", anthropicMessagesUrl, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("anthropic-version", anthropicAPIVersion)
	req.Header.Set("X-Api-Key", p.AnthropicAPIKey)

	return req, nil
}

// GetModelID returns the model id of Anthropic official messages API provider
func (p *AnthropicOfficialMessagesAPIProvider) GetModelID() string {
	return p.AnthropicModelID
}

// GetMaxTokens returns the max tokens to generate of Anthropic official messages API provider
func (p *AnthropicOfficialMessagesAPIProvider) GetMaxTokens() uint32 {
	return p.AnthropicMaxTokens
}

// NewAnthropicLargeLanguageModelProvider creates a new Anthropic large language model provider instance
func NewAnthropicLargeLanguageModelProvider(llmConfig *settings.LLMConfig, enableResponseLog bool) provider.LargeLanguageModelProvider {
	return newCommonAnthropicMessagesAPILargeLanguageModelAdapter(llmConfig, enableResponseLog, &AnthropicOfficialMessagesAPIProvider{
		AnthropicAPIKey:    llmConfig.AnthropicAPIKey,
		AnthropicModelID:   llmConfig.AnthropicModelID,
		AnthropicMaxTokens: llmConfig.AnthropicMaxTokens,
	})
}
