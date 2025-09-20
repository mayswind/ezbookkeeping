package llm

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// LargeLanguageModelProvider defines the structure of large language model provider
type LargeLanguageModelProvider interface {
	// GetJsonResponseByReceiptImageRecognitionModel returns the json response from the large language model provider by receipt image recognition model
	GetJsonResponseByReceiptImageRecognitionModel(c core.Context, uid int64, currentConfig *settings.Config, request *LargeLanguageModelRequest) (*LargeLanguageModelTextualResponse, error)
}
