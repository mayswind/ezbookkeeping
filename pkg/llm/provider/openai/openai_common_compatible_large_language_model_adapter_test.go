package openai

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/llm/data"
)

func TestCommonOpenAIChatCompletionsAPILargeLanguageModelAdapter_buildJsonRequestBody_TextualUserPrompt(t *testing.T) {
	adapter := &CommonOpenAIChatCompletionsAPILargeLanguageModelAdapter{
		apiProvider: &OpenAIOfficialChatCompletionsAPIProvider{
			OpenAIModelID: "test",
		},
	}

	request := &data.LargeLanguageModelRequest{
		SystemPrompt: "You are a helpful assistant.",
		UserPrompt:   []byte("Hello, how are you?"),
	}

	bodyBytes, err := adapter.buildJsonRequestBody(core.NewNullContext(), 0, request, data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.Nil(t, err)

	var body map[string]interface{}
	err = json.Unmarshal(bodyBytes, &body)
	assert.Nil(t, err)

	assert.Equal(t, "{\"model\":\"test\",\"stream\":false,\"messages\":[{\"role\":\"system\",\"content\":\"You are a helpful assistant.\"},{\"role\":\"user\",\"content\":\"Hello, how are you?\"}],\"response_format\":{\"type\":\"json_object\"}}", string(bodyBytes))
}

func TestCommonOpenAIChatCompletionsAPILargeLanguageModelAdapter_buildJsonRequestBody_ImageUserPrompt(t *testing.T) {
	adapter := &CommonOpenAIChatCompletionsAPILargeLanguageModelAdapter{
		apiProvider: &OpenAIOfficialChatCompletionsAPIProvider{
			OpenAIModelID: "test",
		},
	}

	request := &data.LargeLanguageModelRequest{
		SystemPrompt:          "What's in this image?",
		UserPrompt:            []byte("fakedata"),
		UserPromptType:        data.LARGE_LANGUAGE_MODEL_REQUEST_PROMPT_TYPE_IMAGE_URL,
		UserPromptContentType: "image/png",
	}

	bodyBytes, err := adapter.buildJsonRequestBody(core.NewNullContext(), 0, request, data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.Nil(t, err)

	var body map[string]interface{}
	err = json.Unmarshal(bodyBytes, &body)
	assert.Nil(t, err)

	assert.Equal(t, "{\"model\":\"test\",\"stream\":false,\"messages\":[{\"role\":\"system\",\"content\":\"What's in this image?\"},{\"role\":\"user\",\"content\":[{\"type\":\"image_url\",\"image_url\":{\"url\":\"data:image/png;base64,ZmFrZWRhdGE=\"}}]}],\"response_format\":{\"type\":\"json_object\"}}", string(bodyBytes))
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

	result, err := adapter.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
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

	result, err := adapter.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
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

	_, err := adapter.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
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

	_, err := adapter.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.EqualError(t, err, "failed to request third party api")
}

func TestCommonOpenAIChatCompletionsAPILargeLanguageModelAdapter_ParseTextualResponse_InvalidJson(t *testing.T) {
	adapter := &CommonOpenAIChatCompletionsAPILargeLanguageModelAdapter{
		apiProvider: &OpenAIOfficialChatCompletionsAPIProvider{},
	}

	response := "error"

	_, err := adapter.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.EqualError(t, err, "failed to request third party api")
}
