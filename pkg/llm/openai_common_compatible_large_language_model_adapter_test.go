package llm

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
)

func TestCommonOpenAIChatCompletionsAPILargeLanguageModelAdapter_buildJsonRequestBody_TextualUserPrompt(t *testing.T) {
	adapter := &CommonOpenAIChatCompletionsAPILargeLanguageModelAdapter{
		apiProvider: &OpenAIOfficialChatCompletionsAPIProvider{
			OpenAIModelID: "test",
		},
	}

	request := &LargeLanguageModelRequest{
		SystemPrompt: "You are a helpful assistant.",
		UserPrompt:   []byte("Hello, how are you?"),
	}

	bodyBytes, err := adapter.buildJsonRequestBody(core.NewNullContext(), 0, request, LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.Nil(t, err)

	var body map[string]interface{}
	err = json.Unmarshal(bodyBytes, &body)
	assert.Nil(t, err)

	assert.Equal(t, "{\"messages\":[{\"content\":\"You are a helpful assistant.\",\"role\":\"system\"},{\"content\":\"Hello, how are you?\",\"role\":\"user\"}],\"model\":\"test\",\"response_format\":{\"type\":\"json_object\"},\"stream\":false}", string(bodyBytes))
}

func TestCommonOpenAIChatCompletionsAPILargeLanguageModelAdapter_buildJsonRequestBody_ImageUserPrompt(t *testing.T) {
	adapter := &CommonOpenAIChatCompletionsAPILargeLanguageModelAdapter{
		apiProvider: &OpenAIOfficialChatCompletionsAPIProvider{
			OpenAIModelID: "test",
		},
	}

	request := &LargeLanguageModelRequest{
		SystemPrompt:   "What's in this image?",
		UserPrompt:     []byte("fakedata"),
		UserPromptType: LARGE_LANGUAGE_MODEL_REQUEST_PROMPT_TYPE_IMAGE_URL,
	}

	bodyBytes, err := adapter.buildJsonRequestBody(core.NewNullContext(), 0, request, LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.Nil(t, err)

	var body map[string]interface{}
	err = json.Unmarshal(bodyBytes, &body)
	assert.Nil(t, err)

	assert.Equal(t, "{\"messages\":[{\"content\":\"What's in this image?\",\"role\":\"system\"},{\"content\":[{\"image_url\":{\"url\":\"data:image/png;base64,ZmFrZWRhdGE=\"},\"type\":\"image_url\"}],\"role\":\"user\"}],\"model\":\"test\",\"response_format\":{\"type\":\"json_object\"},\"stream\":false}", string(bodyBytes))
}

func TestCommonOpenAIChatCompletionsAPILargeLanguageModelAdapter_ParseTextualResponse_ValidJsonResponse(t *testing.T) {
	adapter := &CommonOpenAIChatCompletionsAPILargeLanguageModelAdapter{
		apiProvider: &OpenAIOfficialChatCompletionsAPIProvider{},
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

	result, err := adapter.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.Nil(t, err)
	assert.Equal(t, "This is a test response", result.Content)
}

func TestCommonOpenAIChatCompletionsAPILargeLanguageModelAdapter_ParseTextualResponse_EmptyResponse(t *testing.T) {
	adapter := &CommonOpenAIChatCompletionsAPILargeLanguageModelAdapter{
		apiProvider: &OpenAIOfficialChatCompletionsAPIProvider{},
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

	result, err := adapter.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.Nil(t, err)
	assert.Equal(t, "", result.Content)
}

func TestCommonOpenAIChatCompletionsAPILargeLanguageModelAdapter_ParseTextualResponse_EmptyChoices(t *testing.T) {
	adapter := &CommonOpenAIChatCompletionsAPILargeLanguageModelAdapter{
		apiProvider: &OpenAIOfficialChatCompletionsAPIProvider{},
	}

	response := `{
		"id": "test-123",
		"object": "chat.completion",
		"choices": []
	}`

	_, err := adapter.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.EqualError(t, err, "failed to request third party api")
}

func TestCommonOpenAIChatCompletionsAPILargeLanguageModelAdapter_ParseTextualResponse_NoChoiceContent(t *testing.T) {
	adapter := &CommonOpenAIChatCompletionsAPILargeLanguageModelAdapter{
		apiProvider: &OpenAIOfficialChatCompletionsAPIProvider{},
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

	_, err := adapter.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.EqualError(t, err, "failed to request third party api")
}

func TestCommonOpenAIChatCompletionsAPILargeLanguageModelAdapter_ParseTextualResponse_InvalidJson(t *testing.T) {
	adapter := &CommonOpenAIChatCompletionsAPILargeLanguageModelAdapter{
		apiProvider: &OpenAIOfficialChatCompletionsAPIProvider{},
	}

	response := "error"

	_, err := adapter.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.EqualError(t, err, "failed to request third party api")
}
