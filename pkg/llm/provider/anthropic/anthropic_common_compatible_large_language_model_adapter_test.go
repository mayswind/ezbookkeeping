package anthropic

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/llm/data"
)

func TestCommonAnthropicMessagesAPILargeLanguageModelAdapter_buildJsonRequestBody_TextualUserPrompt(t *testing.T) {
	adapter := &CommonAnthropicMessagesAPILargeLanguageModelAdapter{
		apiProvider: &AnthropicOfficialMessagesAPIProvider{
			AnthropicModelID:   "test",
			AnthropicMaxTokens: 128,
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

	assert.Equal(t, "{\"model\":\"test\",\"max_tokens\":128,\"stream\":false,\"system\":\"You are a helpful assistant.\",\"messages\":[{\"role\":\"user\",\"content\":\"Hello, how are you?\"}],\"thinking\":{\"type\":\"disabled\"}}", string(bodyBytes))
}

func TestCommonAnthropicMessagesAPILargeLanguageModelAdapter_buildJsonRequestBody_ImageUserPrompt(t *testing.T) {
	adapter := &CommonAnthropicMessagesAPILargeLanguageModelAdapter{
		apiProvider: &AnthropicOfficialMessagesAPIProvider{
			AnthropicModelID:   "test",
			AnthropicMaxTokens: 128,
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

	assert.Equal(t, "{\"model\":\"test\",\"max_tokens\":128,\"stream\":false,\"system\":\"What's in this image?\",\"messages\":[{\"role\":\"user\",\"content\":[{\"source\":{\"data\":\"ZmFrZWRhdGE=\",\"media_type\":\"image/png\",\"type\":\"base64\"},\"type\":\"image\"}]}],\"thinking\":{\"type\":\"disabled\"}}", string(bodyBytes))
}

func TestCommonAnthropicMessagesAPILargeLanguageModelAdapter_ParseTextualResponse_ValidJsonResponse(t *testing.T) {
	adapter := &CommonAnthropicMessagesAPILargeLanguageModelAdapter{
		apiProvider: &AnthropicOfficialMessagesAPIProvider{},
	}

	response := `{
		"id": "test-123",
		"role": "assistant",
		"type": "message",
		"model": "test",
		"usage": {
			"input_tokens": 13,
			"output_tokens": 7
		},
		"content": [
			{
				"type": "text",
				"text": "This is a test response"
			}
		]
	}`

	result, err := adapter.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.Nil(t, err)
	assert.Equal(t, "This is a test response", result.Content)
}

func TestCommonAnthropicMessagesAPILargeLanguageModelAdapter_ParseTextualResponse_EmptyContentText(t *testing.T) {
	adapter := &CommonAnthropicMessagesAPILargeLanguageModelAdapter{
		apiProvider: &AnthropicOfficialMessagesAPIProvider{},
	}

	response := `{
		"id": "test-123",
		"role": "assistant",
		"content": [
			{
				"type": "text",
				"text": ""
			}
		]
	}`

	result, err := adapter.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.Nil(t, err)
	assert.Equal(t, "", result.Content)
}

func TestCommonAnthropicMessagesAPILargeLanguageModelAdapter_ParseTextualResponse_EmptyContent(t *testing.T) {
	adapter := &CommonAnthropicMessagesAPILargeLanguageModelAdapter{
		apiProvider: &AnthropicOfficialMessagesAPIProvider{},
	}

	response := `{
		"id": "test-123",
		"role": "assistant",
		"content": []
	}`

	_, err := adapter.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.EqualError(t, err, "failed to request third party api")
}

func TestCommonAnthropicMessagesAPILargeLanguageModelAdapter_ParseTextualResponse_NoContentText(t *testing.T) {
	adapter := &CommonAnthropicMessagesAPILargeLanguageModelAdapter{
		apiProvider: &AnthropicOfficialMessagesAPIProvider{},
	}

	response := `{
		"id": "msg_123",
		"role": "assistant",
		"content": [
			{
				"type": "text"
			}
		]
	}`

	_, err := adapter.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.EqualError(t, err, "failed to request third party api")
}

func TestCommonAnthropicMessagesAPILargeLanguageModelAdapter_ParseTextualResponse_InvalidJson(t *testing.T) {
	adapter := &CommonAnthropicMessagesAPILargeLanguageModelAdapter{
		apiProvider: &AnthropicOfficialMessagesAPIProvider{},
	}

	response := "error"

	_, err := adapter.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.EqualError(t, err, "failed to request third party api")
}
