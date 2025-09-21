package llm

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
)

func TestOpenAICommonChatCompletionsHttpLargeLanguageModelProvider_buildJsonRequestBody_TextualUserPrompt(t *testing.T) {
	provider := &OpenAICommonChatCompletionsHttpLargeLanguageModelProvider{
		provider: &OpenAILargeLanguageModelProvider{
			OpenAIModelID: "test",
		},
	}

	request := &LargeLanguageModelRequest{
		SystemPrompt: "You are a helpful assistant.",
		UserPrompt:   []byte("Hello, how are you?"),
	}

	bodyBytes, err := provider.buildJsonRequestBody(core.NewNullContext(), 0, request, LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.Nil(t, err)

	var body map[string]interface{}
	err = json.Unmarshal(bodyBytes, &body)
	assert.Nil(t, err)

	assert.Equal(t, "{\"messages\":[{\"content\":\"You are a helpful assistant.\",\"role\":\"system\"},{\"content\":\"Hello, how are you?\",\"role\":\"user\"}],\"model\":\"test\",\"response_format\":{\"type\":\"json_object\"},\"stream\":false}", string(bodyBytes))
}

func TestOpenAICommonChatCompletionsHttpLargeLanguageModelProvider_buildJsonRequestBody_ImageUserPrompt(t *testing.T) {
	provider := &OpenAICommonChatCompletionsHttpLargeLanguageModelProvider{
		provider: &OpenAILargeLanguageModelProvider{
			OpenAIModelID: "test",
		},
	}

	request := &LargeLanguageModelRequest{
		SystemPrompt:   "What's in this image?",
		UserPrompt:     []byte("fakedata"),
		UserPromptType: LARGE_LANGUAGE_MODEL_REQUEST_PROMPT_TYPE_IMAGE_URL,
	}

	bodyBytes, err := provider.buildJsonRequestBody(core.NewNullContext(), 0, request, LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.Nil(t, err)

	var body map[string]interface{}
	err = json.Unmarshal(bodyBytes, &body)
	assert.Nil(t, err)

	assert.Equal(t, "{\"messages\":[{\"content\":\"What's in this image?\",\"role\":\"system\"},{\"content\":[{\"image_url\":{\"url\":\"data:image/png;base64,ZmFrZWRhdGE=\"},\"type\":\"image_url\"}],\"role\":\"user\"}],\"model\":\"test\",\"response_format\":{\"type\":\"json_object\"},\"stream\":false}", string(bodyBytes))
}

func TestOpenAICommonChatCompletionsHttpLargeLanguageModelProvider_ParseTextualResponse_ValidJsonResponse(t *testing.T) {
	provider := &OpenAICommonChatCompletionsHttpLargeLanguageModelProvider{
		provider: &OpenAILargeLanguageModelProvider{},
	}

	response := `{
		"id": "test-123",
		"object": "chat.completion",
		"created": 1234567890,
		"model": "test",
		"usage": {
			"prompt_tokens": 13,
			"completion_tokens": 7,
			"total_tokens": 20
		},
		"choices": [
			{
				"finish_reason": "stop",
				"index": 0,
				"message": {
					"role": "assistant",
					"content": "This is a test response"
				}
			}
		]
	}`

	result, err := provider.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.Nil(t, err)
	assert.Equal(t, "This is a test response", result.Content)
}

func TestOpenAICommonChatCompletionsHttpLargeLanguageModelProvider_ParseTextualResponse_EmptyResponse(t *testing.T) {
	provider := &OpenAICommonChatCompletionsHttpLargeLanguageModelProvider{
		provider: &OpenAILargeLanguageModelProvider{},
	}

	response := `{
		"id": "test-123",
		"object": "chat.completion",
		"choices": [
			{
				"finish_reason": "stop",
				"index": 0,
				"message": {
					"role": "assistant",
					"content": ""
				}
			}
		]
	}`

	result, err := provider.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.Nil(t, err)
	assert.Equal(t, "", result.Content)
}

func TestOpenAICommonChatCompletionsHttpLargeLanguageModelProvider_ParseTextualResponse_EmptyChoices(t *testing.T) {
	provider := &OpenAICommonChatCompletionsHttpLargeLanguageModelProvider{
		provider: &OpenAILargeLanguageModelProvider{},
	}

	response := `{
		"id": "test-123",
		"object": "chat.completion",
		"choices": []
	}`

	_, err := provider.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.EqualError(t, err, "failed to request third party api")
}

func TestOpenAICommonChatCompletionsHttpLargeLanguageModelProvider_ParseTextualResponse_NoChoiceContent(t *testing.T) {
	provider := &OpenAICommonChatCompletionsHttpLargeLanguageModelProvider{
		provider: &OpenAILargeLanguageModelProvider{},
	}

	response := `{
		"id": "chatcmpl-123",
		"object": "chat.completion",
		"choices": [
			{
				"finish_reason": "stop",
				"index": 0,
				"message": {
					"role": "assistant"
				}
			}
		]
	}`

	_, err := provider.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.EqualError(t, err, "failed to request third party api")
}

func TestOpenAICommonChatCompletionsHttpLargeLanguageModelProvider_ParseTextualResponse_InvalidJson(t *testing.T) {
	provider := &OpenAICommonChatCompletionsHttpLargeLanguageModelProvider{
		provider: &OpenAILargeLanguageModelProvider{},
	}

	response := "error"

	_, err := provider.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.EqualError(t, err, "failed to request third party api")
}
