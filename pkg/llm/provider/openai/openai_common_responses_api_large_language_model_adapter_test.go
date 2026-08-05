package openai

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/llm/data"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

type openAIResponsesTestResponse struct {
	Message string `json:"message"`
}

func TestCommonOpenAIResponsesAPILargeLanguageModelAdapter_buildJsonRequestBody_TextualUserPrompt(t *testing.T) {
	adapter := &CommonOpenAIResponsesAPILargeLanguageModelAdapter{
		apiProvider: &OpenAIOfficialResponsesAPIProvider{
			OpenAIModelID: "test",
		},
	}

	request := &data.LargeLanguageModelRequest{
		SystemPrompt: "You are a helpful assistant.",
		UserPrompt:   []byte("Hello, how are you?"),
	}

	bodyBytes, err := adapter.buildJsonRequestBody(core.NewNullContext(), 0, request, data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.Nil(t, err)

	var body map[string]any
	err = json.Unmarshal(bodyBytes, &body)
	assert.Nil(t, err)

	assert.Equal(t, `{"model":"test","stream":false,"input":[{"role":"system","content":"You are a helpful assistant."},{"role":"user","content":"Hello, how are you?"}],"text":{"format":{"type":"json_object"}}}`, string(bodyBytes))
}

func TestCommonOpenAIResponsesAPILargeLanguageModelAdapter_buildJsonRequestBody_ImageUserPrompt(t *testing.T) {
	adapter := &CommonOpenAIResponsesAPILargeLanguageModelAdapter{
		apiProvider: &OpenAIOfficialResponsesAPIProvider{
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

	var body map[string]any
	err = json.Unmarshal(bodyBytes, &body)
	assert.Nil(t, err)

	assert.Equal(t, `{"model":"test","stream":false,"input":[{"role":"system","content":"What's in this image?"},{"role":"user","content":[{"type":"input_image","image_url":"data:image/png;base64,ZmFrZWRhdGE="}]}],"text":{"format":{"type":"json_object"}}}`, string(bodyBytes))
}

func TestCommonOpenAIResponsesAPILargeLanguageModelAdapter_buildJsonRequestBody_ThinkingHighReasoningEffort(t *testing.T) {
	adapter := &CommonOpenAIResponsesAPILargeLanguageModelAdapter{
		apiProvider: &OpenAIOfficialResponsesAPIProvider{
			OpenAIModelID: "test",
		},
		ThinkingLevel: settings.LLMThinkingHigh,
	}

	request := &data.LargeLanguageModelRequest{
		UserPrompt: []byte("Hello, how are you?"),
	}

	bodyBytes, err := adapter.buildJsonRequestBody(core.NewNullContext(), 0, request, data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.Nil(t, err)

	assert.Equal(t, `{"model":"test","stream":false,"input":[{"role":"user","content":"Hello, how are you?"}],"reasoning":{"effort":"high"},"text":{"format":{"type":"json_object"}}}`, string(bodyBytes))
}

func TestCommonOpenAIResponsesAPILargeLanguageModelAdapter_buildJsonRequestBody_JsonSchema(t *testing.T) {
	adapter := &CommonOpenAIResponsesAPILargeLanguageModelAdapter{
		apiProvider: &OpenAIOfficialResponsesAPIProvider{
			OpenAIModelID: "test",
		},
	}

	request := &data.LargeLanguageModelRequest{
		UserPrompt:             []byte("Hello"),
		ResponseJsonObjectType: reflect.TypeOf(openAIResponsesTestResponse{}),
	}

	bodyBytes, err := adapter.buildJsonRequestBody(core.NewNullContext(), 0, request, data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.Nil(t, err)

	var body map[string]any
	err = json.Unmarshal(bodyBytes, &body)
	assert.Nil(t, err)

	textConfig := body["text"].(map[string]any)
	format := textConfig["format"].(map[string]any)
	assert.Equal(t, "json_schema", format["type"])
	assert.Equal(t, "response", format["name"])
	assert.Equal(t, true, format["strict"])
	assert.NotNil(t, format["schema"])
}

func TestCommonOpenAIResponsesAPILargeLanguageModelAdapter_ParseTextualResponse_ValidJsonResponse(t *testing.T) {
	adapter := &CommonOpenAIResponsesAPILargeLanguageModelAdapter{}
	response := `{
		"output": [
			{
				"type": "reasoning",
				"summary": []
			},
			{
				"type": "message",
				"content": [
					{
						"type": "output_text",
						"text": "This is a test response"
					}
				]
			}
		]
	}`

	result, err := adapter.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.Nil(t, err)
	assert.Equal(t, "This is a test response", result.Content)
}

func TestCommonOpenAIResponsesAPILargeLanguageModelAdapter_ParseTextualResponse_EmptyResponse(t *testing.T) {
	adapter := &CommonOpenAIResponsesAPILargeLanguageModelAdapter{}
	response := `{
		"output": [
			{
				"type": "message",
				"content": [
					{
						"type": "output_text",
						"text": ""
					}
				]
			}
		]
	}`

	result, err := adapter.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.Nil(t, err)
	assert.Equal(t, "", result.Content)
}

func TestCommonOpenAIResponsesAPILargeLanguageModelAdapter_ParseTextualResponse_NoOutputText(t *testing.T) {
	adapter := &CommonOpenAIResponsesAPILargeLanguageModelAdapter{
		apiProvider: &OpenAIOfficialResponsesAPIProvider{},
	}

	response := `{
		"output": [
			{
				"type": "message",
				"content": [
					{
						"type": "refusal",
						"refusal": "I cannot help with that."
					}
				]
			}
		]
	}`

	_, err := adapter.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.EqualError(t, err, "failed to request third party api")
}

func TestCommonOpenAIResponsesAPILargeLanguageModelAdapter_ParseTextualResponse_InvalidJson(t *testing.T) {
	adapter := &CommonOpenAIResponsesAPILargeLanguageModelAdapter{
		apiProvider: &OpenAIOfficialResponsesAPIProvider{},
	}

	response := "error"

	_, err := adapter.ParseTextualResponse(core.NewNullContext(), 0, []byte(response), data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
	assert.EqualError(t, err, "failed to request third party api")
}
