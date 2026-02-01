package lmstudio

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/llm/data"
)

func TestLMStudioLargeLanguageModelAdapter_buildJsonRequestBody_TextualUserPrompt(t *testing.T) {
	adapter := &LMStudioLargeLanguageModelAdapter{
		LMStudioModelID: "test",
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

	assert.Equal(t, "{\"model\":\"test\",\"stream\":false,\"system_prompt\":\"You are a helpful assistant.\",\"input\":[{\"type\":\"text\",\"content\":\"Hello, how are you?\"}]}", string(bodyBytes))
}

func TestLMStudioLargeLanguageModelAdapter_buildJsonRequestBody_ImageUserPrompt(t *testing.T) {
	adapter := &LMStudioLargeLanguageModelAdapter{
		LMStudioModelID: "test",
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

	assert.Equal(t, "{\"model\":\"test\",\"stream\":false,\"system_prompt\":\"What's in this image?\",\"input\":[{\"type\":\"image\",\"data_url\":\"data:image/png;base64,ZmFrZWRhdGE=\"}]}", string(bodyBytes))
}

func TestLMStudioLargeLanguageModelAdapter_ParseTextualResponse_ValidJsonResponse(t *testing.T) {
	adapter := &LMStudioLargeLanguageModelAdapter{}

	response := `{
		"model_instance_id": "test",
		"output": [
			{
				"type": "message",
				"content": "This is a test response"
			}
		]
	}`

	result, err := adapter.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.Nil(t, err)
	assert.Equal(t, "This is a test response", result.Content)
}

func TestLMStudioLargeLanguageModelAdapter_ParseTextualResponse_EmptyOutputContent(t *testing.T) {
	adapter := &LMStudioLargeLanguageModelAdapter{}

	response := `{
		"model_instance_id": "test",
		"output": [
			{
				"type": "message",
				"content": ""
			}
		]
	}`

	result, err := adapter.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.Nil(t, err)
	assert.Equal(t, "", result.Content)
}

func TestLMStudioLargeLanguageModelAdapter_ParseTextualResponse_EmptyOutput(t *testing.T) {
	adapter := &LMStudioLargeLanguageModelAdapter{}

	response := `{
		"model_instance_id": "test",
		"output": []
	}`

	_, err := adapter.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.EqualError(t, err, "failed to request third party api")
}

func TestLMStudioLargeLanguageModelAdapter_ParseTextualResponse_NoContentFieldInOutput(t *testing.T) {
	adapter := &LMStudioLargeLanguageModelAdapter{}

	response := `{
		"model_instance_id": "test",
		"output": [
			{
				"type": "message"
			}
		]
	}`

	_, err := adapter.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.EqualError(t, err, "failed to request third party api")
}

func TestLMStudioLargeLanguageModelAdapter_ParseTextualResponse_InvalidJson(t *testing.T) {
	adapter := &LMStudioLargeLanguageModelAdapter{}

	response := "error"

	_, err := adapter.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.EqualError(t, err, "failed to request third party api")
}

func TestLMStudioLargeLanguageModelAdapter_GetOllamaRequestUrl(t *testing.T) {
	adapter := &LMStudioLargeLanguageModelAdapter{
		LMStudioServerURL: "http://localhost:1234/",
	}
	url := adapter.getLMStudioRequestUrl()
	assert.Equal(t, "http://localhost:1234/api/v1/chat", url)

	adapter = &LMStudioLargeLanguageModelAdapter{
		LMStudioServerURL: "http://localhost:1234",
	}
	url = adapter.getLMStudioRequestUrl()
	assert.Equal(t, "http://localhost:1234/api/v1/chat", url)

	adapter = &LMStudioLargeLanguageModelAdapter{
		LMStudioServerURL: "http://example.com/lmstudio/",
	}
	url = adapter.getLMStudioRequestUrl()
	assert.Equal(t, "http://example.com/lmstudio/api/v1/chat", url)
}
