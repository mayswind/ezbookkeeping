package llm

import "reflect"

type LargeLanguageModelRequestPromptType byte

// Large Language Model Request Prompt Type
const (
	LARGE_LANGUAGE_MODEL_REQUEST_PROMPT_TYPE_TEXT      LargeLanguageModelRequestPromptType = 0
	LARGE_LANGUAGE_MODEL_REQUEST_PROMPT_TYPE_IMAGE_URL LargeLanguageModelRequestPromptType = 1
)

type LargeLanguageModelResponseFormat byte

// Large Language Model Response Format
const (
	LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_TEXT LargeLanguageModelResponseFormat = 0
	LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON LargeLanguageModelResponseFormat = 1
)

// LargeLanguageModelRequest represents a request to a large language model
type LargeLanguageModelRequest struct {
	Stream                 bool
	SystemPrompt           string
	UserPrompt             []byte
	UserPromptType         LargeLanguageModelRequestPromptType
	ResponseJsonObjectType reflect.Type
}

// LargeLanguageModelTextualResponse represents a textual response from a large language model
type LargeLanguageModelTextualResponse struct {
	Content string
}
