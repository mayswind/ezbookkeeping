package llm

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// LargeLanguageModelProviderContainer contains the current large language model provider
type LargeLanguageModelProviderContainer struct {
	current LargeLanguageModelProvider
}

// Initialize a large language model provider container singleton instance
var (
	Container = &LargeLanguageModelProviderContainer{}
)

// InitializeLargeLanguageModelProvider initializes the current large language model provider according to the config
func InitializeLargeLanguageModelProvider(config *settings.Config) error {
	if config.LLMProvider == settings.OpenAILLMProvider {
		Container.current = NewOpenAILargeLanguageModelProvider(config)
		return nil
	} else if config.LLMProvider == settings.OpenAICompatibleLLMProvider {
		Container.current = NewOpenAICompatibleLargeLanguageModelProvider(config)
		return nil
	} else if config.LLMProvider == settings.OpenRouterLLMProvider {
		Container.current = NewOpenRouterLargeLanguageModelProvider(config)
		return nil
	} else if config.LLMProvider == settings.OllamaLLMProvider {
		Container.current = NewOllamaLargeLanguageModelProvider(config)
		return nil
	}

	return errs.ErrInvalidLLMProvider
}

// GetJsonResponseByReceiptImageRecognitionModel returns the json response from the current large language model provider by receipt image recognition model
func (l *LargeLanguageModelProviderContainer) GetJsonResponseByReceiptImageRecognitionModel(c core.Context, uid int64, currentConfig *settings.Config, request *LargeLanguageModelRequest) (*LargeLanguageModelTextualResponse, error) {
	if Container.current == nil {
		return nil, errs.ErrInvalidLLMProvider
	}

	return l.current.GetJsonResponseByReceiptImageRecognitionModel(c, uid, currentConfig, request)
}
