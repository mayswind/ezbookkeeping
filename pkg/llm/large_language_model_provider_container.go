package llm

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// LargeLanguageModelProviderContainer contains the current large language model provider
type LargeLanguageModelProviderContainer struct {
	receiptImageRecognitionCurrentProvider LargeLanguageModelProvider
}

// Initialize a large language model provider container singleton instance
var (
	Container = &LargeLanguageModelProviderContainer{}
)

// InitializeLargeLanguageModelProvider initializes the current large language model provider according to the config
func InitializeLargeLanguageModelProvider(config *settings.Config) error {
	var err error = nil

	if config.ReceiptImageRecognitionLLMConfig != nil {
		Container.receiptImageRecognitionCurrentProvider, err = initializeLargeLanguageModelProvider(config.ReceiptImageRecognitionLLMConfig)

		if err != nil {
			return err
		}
	}

	return nil
}

func initializeLargeLanguageModelProvider(llmConfig *settings.LLMConfig) (LargeLanguageModelProvider, error) {
	if llmConfig.LLMProvider == settings.OpenAILLMProvider {
		return NewOpenAILargeLanguageModelProvider(llmConfig), nil
	} else if llmConfig.LLMProvider == settings.OpenAICompatibleLLMProvider {
		return NewOpenAICompatibleLargeLanguageModelProvider(llmConfig), nil
	} else if llmConfig.LLMProvider == settings.OpenRouterLLMProvider {
		return NewOpenRouterLargeLanguageModelProvider(llmConfig), nil
	} else if llmConfig.LLMProvider == settings.OllamaLLMProvider {
		return NewOllamaLargeLanguageModelProvider(llmConfig), nil
	} else if llmConfig.LLMProvider == "" {
		return nil, nil
	}

	return nil, errs.ErrInvalidLLMProvider
}

// GetJsonResponseByReceiptImageRecognitionModel returns the json response from the current large language model provider by receipt image recognition model
func (l *LargeLanguageModelProviderContainer) GetJsonResponseByReceiptImageRecognitionModel(c core.Context, uid int64, currentConfig *settings.Config, request *LargeLanguageModelRequest) (*LargeLanguageModelTextualResponse, error) {
	if currentConfig.ReceiptImageRecognitionLLMConfig == nil || Container.receiptImageRecognitionCurrentProvider == nil {
		return nil, errs.ErrInvalidLLMProvider
	}

	return l.receiptImageRecognitionCurrentProvider.GetJsonResponse(c, uid, currentConfig.ReceiptImageRecognitionLLMConfig, request)
}
