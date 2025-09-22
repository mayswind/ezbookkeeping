package llm

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/invopop/jsonschema"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
)

// OpenAIChatCompletionsAPIProvider defines the structure of OpenAI chat completions API provider
type OpenAIChatCompletionsAPIProvider interface {
	// BuildChatCompletionsHttpRequest returns the chat completions http request
	BuildChatCompletionsHttpRequest(c core.Context, uid int64) (*http.Request, error)

	// GetModelID returns the model id if supported, otherwise returns empty string
	GetModelID() string
}

// CommonOpenAIChatCompletionsAPILargeLanguageModelAdapter defines the structure of OpenAI common compatible large language model adapter based on chat completions api
type CommonOpenAIChatCompletionsAPILargeLanguageModelAdapter struct {
	HttpLargeLanguageModelAdapter
	apiProvider OpenAIChatCompletionsAPIProvider
}

// BuildTextualRequest returns the http request by OpenAI common compatible adapter
func (p *CommonOpenAIChatCompletionsAPILargeLanguageModelAdapter) BuildTextualRequest(c core.Context, uid int64, request *LargeLanguageModelRequest, responseType LargeLanguageModelResponseFormat) (*http.Request, error) {
	requestBody, err := p.buildJsonRequestBody(c, uid, request, responseType)

	if err != nil {
		return nil, err
	}

	httpRequest, err := p.apiProvider.BuildChatCompletionsHttpRequest(c, uid)

	if err != nil {
		return nil, err
	}

	httpRequest.Body = io.NopCloser(bytes.NewReader(requestBody))
	httpRequest.Header.Set("Content-Type", "application/json")

	return httpRequest, nil
}

// ParseTextualResponse returns the textual response by OpenAI common compatible adapter
func (p *CommonOpenAIChatCompletionsAPILargeLanguageModelAdapter) ParseTextualResponse(c core.Context, uid int64, body []byte, responseType LargeLanguageModelResponseFormat) (*LargeLanguageModelTextualResponse, error) {
	responseBody := make(map[string]any)
	err := json.Unmarshal(body, &responseBody)

	if err != nil {
		log.Errorf(c, "[openai_common_compatible_large_language_model_adapter.ParseTextualResponse] failed to parse response for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	choices, ok := responseBody["choices"].([]any)

	if !ok || len(choices) < 1 {
		log.Errorf(c, "[openai_common_compatible_large_language_model_adapter.ParseTextualResponse] no choices found in response for user \"uid:%d\"", uid)
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	firstChoice, ok := choices[0].(map[string]any)

	if !ok {
		log.Errorf(c, "[openai_common_compatible_large_language_model_adapter.ParseTextualResponse] invalid choice format in response for user \"uid:%d\"", uid)
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	message, ok := firstChoice["message"].(map[string]any)

	if !ok {
		log.Errorf(c, "[openai_common_compatible_large_language_model_adapter.ParseTextualResponse] no message found in choice for user \"uid:%d\"", uid)
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	content, ok := message["content"].(string)

	if !ok {
		log.Errorf(c, "[openai_common_compatible_large_language_model_adapter.ParseTextualResponse] no content found in message for user \"uid:%d\"", uid)
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	if responseType == LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON {
		if strings.HasPrefix(content, "```json") && strings.HasSuffix(content, "```") {
			content = strings.TrimPrefix(content, "```json")
			content = strings.TrimSuffix(content, "```")
		} else if strings.HasPrefix(content, "```") && strings.HasSuffix(content, "```") {
			content = strings.TrimPrefix(content, "```")
			content = strings.TrimSuffix(content, "```")
		}
	}

	textualResponse := &LargeLanguageModelTextualResponse{
		Content: content,
	}

	return textualResponse, nil
}

func (p *CommonOpenAIChatCompletionsAPILargeLanguageModelAdapter) buildJsonRequestBody(c core.Context, uid int64, request *LargeLanguageModelRequest, responseType LargeLanguageModelResponseFormat) ([]byte, error) {
	if p.apiProvider.GetModelID() == "" {
		return nil, errs.ErrInvalidLLMModelId
	}

	requestMessages := make([]any, 0)

	if request.SystemPrompt != "" {
		requestMessages = append(requestMessages, map[string]string{
			"role":    "system",
			"content": request.SystemPrompt,
		})
	}

	if len(request.UserPrompt) > 0 {
		if request.UserPromptType == LARGE_LANGUAGE_MODEL_REQUEST_PROMPT_TYPE_IMAGE_URL {
			imageBase64Data := "data:image/png;base64," + base64.StdEncoding.EncodeToString(request.UserPrompt)
			requestMessages = append(requestMessages, map[string]any{
				"role": "user",
				"content": []any{
					core.O{
						"type": "image_url",
						"image_url": core.O{
							"url": imageBase64Data,
						},
					},
				},
			})
		} else {
			requestMessages = append(requestMessages, map[string]string{
				"role":    "user",
				"content": string(request.UserPrompt),
			})
		}
	}

	requestBody := make(map[string]any)
	requestBody["model"] = p.apiProvider.GetModelID()
	requestBody["stream"] = request.Stream
	requestBody["messages"] = requestMessages

	if responseType == LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON {
		if request.ResponseJsonObjectType != nil {
			schemeGenerator := jsonschema.Reflector{
				Anonymous:      true,
				DoNotReference: true,
				ExpandedStruct: true,
			}

			schema := schemeGenerator.ReflectFromType(request.ResponseJsonObjectType)
			schema.Version = ""

			requestBody["response_format"] = core.O{
				"type":        "json_schema",
				"json_schema": schema,
			}
		} else {
			requestBody["response_format"] = core.O{
				"type": "json_object",
			}
		}
	}

	requestBodyBytes, err := json.Marshal(requestBody)

	if err != nil {
		log.Errorf(c, "[openai_common_compatible_large_language_model_adapter.buildJsonRequestBody] failed to marshal request body for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	log.Debugf(c, "[openai_common_compatible_large_language_model_adapter.buildJsonRequestBody] request body is %s", requestBodyBytes)
	return requestBodyBytes, nil
}

func newCommonOpenAIChatCompletionsAPILargeLanguageModelAdapter(apiProvider OpenAIChatCompletionsAPIProvider) LargeLanguageModelProvider {
	return newCommonHttpLargeLanguageModelProvider(&CommonOpenAIChatCompletionsAPILargeLanguageModelAdapter{
		apiProvider: apiProvider,
	})
}
