package ollama

import (
	"encoding/json"
	"testing"

	"github.com/mayswind/ezbookkeeping/pkg/llm/data"
	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
)

func TestOllamaLargeLanguageModelAdapter_buildJsonRequestBody_TextualUserPrompt(t *testing.T) {
	adapter := &OllamaLargeLanguageModelAdapter{
		OllamaModelID: "test",
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

	assert.Equal(t, "{\"model\":\"test\",\"stream\":false,\"messages\":[{\"role\":\"system\",\"content\":\"You are a helpful assistant.\"},{\"role\":\"user\",\"content\":\"Hello, how are you?\"}],\"format\":\"json\"}", string(bodyBytes))
}

func TestOllamaLargeLanguageModelAdapter_buildJsonRequestBody_ImageUserPrompt(t *testing.T) {
	adapter := &OllamaLargeLanguageModelAdapter{
		OllamaModelID: "test",
	}

	request := &data.LargeLanguageModelRequest{
		SystemPrompt:   "What's in this image?",
		UserPrompt:     []byte("fakedata"),
		UserPromptType: data.LARGE_LANGUAGE_MODEL_REQUEST_PROMPT_TYPE_IMAGE_URL,
	}

	bodyBytes, err := adapter.buildJsonRequestBody(core.NewNullContext(), 0, request, data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.Nil(t, err)

	var body map[string]interface{}
	err = json.Unmarshal(bodyBytes, &body)
	assert.Nil(t, err)

	assert.Equal(t, "{\"model\":\"test\",\"stream\":false,\"messages\":[{\"role\":\"system\",\"content\":\"What's in this image?\"},{\"role\":\"user\",\"content\":\"\",\"images\":[\"ZmFrZWRhdGE=\"]}],\"format\":\"json\"}", string(bodyBytes))
}

func TestOllamaLargeLanguageModelAdapter_ParseTextualResponse_ValidJsonResponse(t *testing.T) {
	adapter := &OllamaLargeLanguageModelAdapter{}

	response := `{
		"model": "test",
		"created_at": "2025-09-01T01:02:03.456789Z",
		"message": {
			"role": "assistant",
			"content": "This is a test response"
		}
	}`

	result, err := adapter.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.Nil(t, err)
	assert.Equal(t, "This is a test response", result.Content)
}

func TestOllamaLargeLanguageModelAdapter_ParseTextualResponse_EmptyResponse(t *testing.T) {
	adapter := &OllamaLargeLanguageModelAdapter{}

	response := `{
		"model": "test",
		"created_at": "2025-09-01T01:02:03.456789Z",
		"message": {
			"role": "assistant",
			"content": ""
		}
	}`

	result, err := adapter.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.Nil(t, err)
	assert.Equal(t, "", result.Content)
}

func TestOllamaLargeLanguageModelAdapter_ParseTextualResponse_EmptyMessage(t *testing.T) {
	adapter := &OllamaLargeLanguageModelAdapter{}

	response := `{
		"model": "test",
		"created_at": "2025-09-01T01:02:03.456789Z",
		"message": {}
	}`

	_, err := adapter.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.EqualError(t, err, "failed to request third party api")
}

func TestOllamaLargeLanguageModelAdapter_ParseTextualResponse_NoContentFieldInMessage(t *testing.T) {
	adapter := &OllamaLargeLanguageModelAdapter{}

	response := `{
		"model": "test",
		"created_at": "2025-09-01T01:02:03.456789Z",
		"message": {
			"role": "assistant"
		}
	}`

	_, err := adapter.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.EqualError(t, err, "failed to request third party api")
}

func TestOllamaLargeLanguageModelAdapter_ParseTextualResponse_InvalidJson(t *testing.T) {
	adapter := &OllamaLargeLanguageModelAdapter{}

	response := "error"

	_, err := adapter.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.EqualError(t, err, "failed to request third party api")
}

func TestOllamaLargeLanguageModelAdapter_GetOllamaRequestUrl(t *testing.T) {
	adapter := &OllamaLargeLanguageModelAdapter{
		OllamaServerURL: "http://localhost:11434/",
	}
	url := adapter.getOllamaRequestUrl()
	assert.Equal(t, "http://localhost:11434/api/chat", url)

	adapter = &OllamaLargeLanguageModelAdapter{
		OllamaServerURL: "http://localhost:11434",
	}
	url = adapter.getOllamaRequestUrl()
	assert.Equal(t, "http://localhost:11434/api/chat", url)

	adapter = &OllamaLargeLanguageModelAdapter{
		OllamaServerURL: "http://example.com/ollama/",
	}
	url = adapter.getOllamaRequestUrl()
	assert.Equal(t, "http://example.com/ollama/api/chat", url)
}
