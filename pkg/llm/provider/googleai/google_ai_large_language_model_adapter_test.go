package googleai

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/llm/data"
)

func TestGoogleAILargeLanguageModelAdapter_buildJsonRequestBody_TextualUserPrompt(t *testing.T) {
	adapter := &GoogleAILargeLanguageModelAdapter{
		GoogleAIModelID: "test",
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

	assert.Equal(t, "{\"contents\":[{\"parts\":[{\"text\":\"You are a helpful assistant.\"},{\"text\":\"Hello, how are you?\"}]}]}", string(bodyBytes))
}

func TestGoogleAILargeLanguageModelAdapter_buildJsonRequestBody_ImageUserPrompt(t *testing.T) {
	adapter := &GoogleAILargeLanguageModelAdapter{
		GoogleAIModelID: "test",
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

	assert.Equal(t, "{\"contents\":[{\"parts\":[{\"text\":\"What's in this image?\"},{\"inlineData\":{\"mimeType\":\"image/png\",\"data\":\"ZmFrZWRhdGE=\"}}]}]}", string(bodyBytes))
}

func TestGoogleAILargeLanguageModelAdapter_ParseTextualResponse_ValidJsonResponse(t *testing.T) {
	adapter := &GoogleAILargeLanguageModelAdapter{
		GoogleAIModelID: "test",
	}

	response := `{
		"responseId": "test-123",
		"modelVersion": "test",
		"usageMetadata": {
			"promptTokenCount": 13,
			"candidatesTokenCount": 7,
			"totalTokenCount": 20
		},
		"candidates": [
			{
				"content": {
					"parts": [
						{
							"text": "This is a test response"
						}
					]
				},
				"finish_reason": "stop",
				"index": 0
			}
		]
	}`

	result, err := adapter.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.Nil(t, err)
	assert.Equal(t, "This is a test response", result.Content)
}

func TestGoogleAILargeLanguageModelAdapter_ParseTextualResponse_EmptyResponse(t *testing.T) {
	adapter := &GoogleAILargeLanguageModelAdapter{
		GoogleAIModelID: "test",
	}

	response := `{
		"responseId": "test-123",
		"modelVersion": "test",
		"usageMetadata": {
			"promptTokenCount": 13,
			"candidatesTokenCount": 7,
			"totalTokenCount": 20
		},
		"candidates": [
			{
				"content": {
					"parts": [
						{
							"text": ""
						}
					]
				},
				"finish_reason": "stop",
				"index": 0
			}
		]
	}`

	result, err := adapter.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.Nil(t, err)
	assert.Equal(t, "", result.Content)
}

func TestGoogleAILargeLanguageModelAdapter_ParseTextualResponse_EmptyCandidates(t *testing.T) {
	adapter := &GoogleAILargeLanguageModelAdapter{
		GoogleAIModelID: "test",
	}

	response := `{
		"responseId": "test-123",
		"modelVersion": "test",
		"usageMetadata": {
			"promptTokenCount": 13,
			"candidatesTokenCount": 7,
			"totalTokenCount": 20
		},
		"candidates": []
	}`

	_, err := adapter.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.EqualError(t, err, "failed to request third party api")
}

func TestGoogleAILargeLanguageModelAdapter_ParseTextualResponse_NoPartText(t *testing.T) {
	adapter := &GoogleAILargeLanguageModelAdapter{
		GoogleAIModelID: "test",
	}

	response := `{
		"responseId": "test-123",
		"modelVersion": "test",
		"usageMetadata": {
			"promptTokenCount": 13,
			"candidatesTokenCount": 7,
			"totalTokenCount": 20
		},
		"candidates": [
			{
				"content": {
					"parts": [
						{
						}
					]
				},
				"finish_reason": "stop",
				"index": 0
			}
		]
	}`

	_, err := adapter.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.EqualError(t, err, "failed to request third party api")
}

func TestGoogleAILargeLanguageModelAdapter_ParseTextualResponse_InvalidJson(t *testing.T) {
	adapter := &GoogleAILargeLanguageModelAdapter{
		GoogleAIModelID: "test",
	}

	response := "error"

	_, err := adapter.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.EqualError(t, err, "failed to request third party api")
}
