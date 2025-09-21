package errs

import "net/http"

// Error codes related to large language model features
var (
	ErrLargeLanguageModelProviderNotEnabled = NewNormalError(NormalSubcategoryLargeLanguageModel, 0, http.StatusBadRequest, "llm provider is not enabled")
	ErrNoAIRecognitionImage                 = NewNormalError(NormalSubcategoryLargeLanguageModel, 1, http.StatusBadRequest, "no image for AI recognition")
	ErrAIRecognitionImageIsEmpty            = NewNormalError(NormalSubcategoryLargeLanguageModel, 2, http.StatusBadRequest, "image for AI recognition is empty")
	ErrExceedMaxAIRecognitionImageFileSize  = NewNormalError(NormalSubcategoryLargeLanguageModel, 3, http.StatusBadRequest, "exceed the maximum size of image file for AI recognition")
	ErrNoTransactionInformationInImage      = NewNormalError(NormalSubcategoryLargeLanguageModel, 4, http.StatusBadRequest, "no transaction information detected")
)
