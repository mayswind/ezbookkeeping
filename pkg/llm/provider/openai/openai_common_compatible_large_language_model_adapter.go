package openai

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"

	"github.com/invopop/jsonschema"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/llm/data"
	"github.com/mayswind/ezbookkeeping/pkg/llm/provider"
	"github.com/mayswind/ezbookkeeping/pkg/llm/provider/common"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
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
	common.HttpLargeLanguageModelAdapter
	apiProvider OpenAIChatCompletionsAPIProvider
}

// OpenAIMessageRole defines the role of OpenAI chat completions message
type OpenAIMessageRole string

// OpenAI Message Roles
const (
	OpenAIMessageRoleSystem OpenAIMessageRole = "system"
	OpenAIMessageRoleUser   OpenAIMessageRole = "user"
)

// OpenAIChatCompletionsRequestResponseFormatType defines the type of OpenAI chat completions request response format
type OpenAIChatCompletionsRequestResponseFormatType string

// OpenAI Chat Completions Request Response Format Types
const (
	OpenAIChatCompletionsRequestResponseFormatTypeJsonObject OpenAIChatCompletionsRequestResponseFormatType = "json_object"
	OpenAIChatCompletionsRequestResponseFormatTypeJsonSchema OpenAIChatCompletionsRequestResponseFormatType = "json_schema"
)

// OpenAIChatCompletionsRequest defines the structure of OpenAI chat completions request
type OpenAIChatCompletionsRequest struct {
	Model          string                                      `json:"model"`
	Stream         bool                                        `json:"stream"`
	Messages       []any                                       `json:"messages"`
	ResponseFormat *OpenAIChatCompletionsRequestResponseFormat `json:"response_format,omitempty"`
}

// OpenAIChatCompletionsRequestMessage defines the structure of OpenAI chat completions request message
type OpenAIChatCompletionsRequestMessage[T string | []*OpenAIChatCompletionsRequestImageContent] struct {
	Role    OpenAIMessageRole `json:"role"`
	Content T                 `json:"content"`
}

// OpenAIChatCompletionsRequestImageContent defines the structure of OpenAI chat completions request image content
type OpenAIChatCompletionsRequestImageContent struct {
	Type     string                                `json:"type"`
	ImageURL *OpenAIChatCompletionsRequestImageUrl `json:"image_url"`
}

// OpenAIChatCompletionsRequestResponseFormat defines the structure of OpenAI chat completions request response format
type OpenAIChatCompletionsRequestResponseFormat struct {
	Type       OpenAIChatCompletionsRequestResponseFormatType `json:"type"`
	JsonSchema *jsonschema.Schema                             `json:"json_schema,omitempty"`
}

// OpenAIChatCompletionsRequestImageUrl defines the structure of OpenAI image url
type OpenAIChatCompletionsRequestImageUrl struct {
	Url string `json:"url"`
}

// OpenAIChatCompletionsResponse defines the structure of OpenAI chat completions response
type OpenAIChatCompletionsResponse struct {
	Choices []*OpenAIChatCompletionsResponseChoice `json:"choices"`
}

// OpenAIChatCompletionsResponseChoice defines the structure of OpenAI chat completions response choice
type OpenAIChatCompletionsResponseChoice struct {
	Message *OpenAIChatCompletionsResponseMessage `json:"message"`
}

// OpenAIChatCompletionsResponseMessage defines the structure of OpenAI chat completions response message
type OpenAIChatCompletionsResponseMessage struct {
	Content *string `json:"content"`
}

// BuildTextualRequest returns the http request by OpenAI common compatible adapter
func (p *CommonOpenAIChatCompletionsAPILargeLanguageModelAdapter) BuildTextualRequest(c core.Context, uid int64, request *data.LargeLanguageModelRequest, responseType data.LargeLanguageModelResponseFormat) (*http.Request, error) {
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
func (p *CommonOpenAIChatCompletionsAPILargeLanguageModelAdapter) ParseTextualResponse(c core.Context, uid int64, body []byte, responseType data.LargeLanguageModelResponseFormat) (*data.LargeLanguageModelTextualResponse, error) {
	chatCompletionsResponse := &OpenAIChatCompletionsResponse{}
	err := json.Unmarshal(body, &chatCompletionsResponse)

	if err != nil {
		log.Errorf(c, "[openai_common_compatible_large_language_model_adapter.ParseTextualResponse] failed to parse chat completions response for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	if chatCompletionsResponse == nil || chatCompletionsResponse.Choices == nil || len(chatCompletionsResponse.Choices) < 1 ||
		chatCompletionsResponse.Choices[0].Message == nil ||
		chatCompletionsResponse.Choices[0].Message.Content == nil {
		log.Errorf(c, "[openai_common_compatible_large_language_model_adapter.ParseTextualResponse] chat completions response is invalid for user \"uid:%d\"", uid)
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	textualResponse := &data.LargeLanguageModelTextualResponse{
		Content: *chatCompletionsResponse.Choices[0].Message.Content,
	}

	return textualResponse, nil
}

func (p *CommonOpenAIChatCompletionsAPILargeLanguageModelAdapter) buildJsonRequestBody(c core.Context, uid int64, request *data.LargeLanguageModelRequest, responseType data.LargeLanguageModelResponseFormat) ([]byte, error) {
	if p.apiProvider.GetModelID() == "" {
		return nil, errs.ErrInvalidLLMModelId
	}

	chatCompletionsRequest := &OpenAIChatCompletionsRequest{
		Model:    p.apiProvider.GetModelID(),
		Stream:   request.Stream,
		Messages: make([]any, 0, 2),
	}

	if request.SystemPrompt != "" {
		chatCompletionsRequest.Messages = append(chatCompletionsRequest.Messages, &OpenAIChatCompletionsRequestMessage[string]{
			Role:    OpenAIMessageRoleSystem,
			Content: request.SystemPrompt,
		})
	}

	if len(request.UserPrompt) > 0 {
		if request.UserPromptType == data.LARGE_LANGUAGE_MODEL_REQUEST_PROMPT_TYPE_IMAGE_URL {
			imageBase64Data := "data:" + request.UserPromptContentType + ";base64," + base64.StdEncoding.EncodeToString(request.UserPrompt)
			chatCompletionsRequest.Messages = append(chatCompletionsRequest.Messages, &OpenAIChatCompletionsRequestMessage[[]*OpenAIChatCompletionsRequestImageContent]{
				Role: OpenAIMessageRoleUser,
				Content: []*OpenAIChatCompletionsRequestImageContent{
					{
						Type: "image_url",
						ImageURL: &OpenAIChatCompletionsRequestImageUrl{
							Url: imageBase64Data,
						},
					},
				},
			})
		} else {
			chatCompletionsRequest.Messages = append(chatCompletionsRequest.Messages, &OpenAIChatCompletionsRequestMessage[string]{
				Role:    OpenAIMessageRoleUser,
				Content: string(request.UserPrompt),
			})
		}
	}

	if responseType == data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON {
		if request.ResponseJsonObjectType != nil {
			schemeGenerator := jsonschema.Reflector{
				Anonymous:      true,
				DoNotReference: true,
				ExpandedStruct: true,
			}

			schema := schemeGenerator.ReflectFromType(request.ResponseJsonObjectType)
			schema.Version = ""

			chatCompletionsRequest.ResponseFormat = &OpenAIChatCompletionsRequestResponseFormat{
				Type:       OpenAIChatCompletionsRequestResponseFormatTypeJsonSchema,
				JsonSchema: schema,
			}
		} else {
			chatCompletionsRequest.ResponseFormat = &OpenAIChatCompletionsRequestResponseFormat{
				Type: OpenAIChatCompletionsRequestResponseFormatTypeJsonObject,
			}
		}
	}

	requestBodyBytes, err := json.Marshal(chatCompletionsRequest)

	if err != nil {
		log.Errorf(c, "[openai_common_compatible_large_language_model_adapter.buildJsonRequestBody] failed to marshal request body for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	log.Debugf(c, "[openai_common_compatible_large_language_model_adapter.buildJsonRequestBody] request body is %s", requestBodyBytes)
	return requestBodyBytes, nil
}

func newCommonOpenAIChatCompletionsAPILargeLanguageModelAdapter(llmConfig *settings.LLMConfig, enableResponseLog bool, apiProvider OpenAIChatCompletionsAPIProvider) provider.LargeLanguageModelProvider {
	return common.NewCommonHttpLargeLanguageModelProvider(llmConfig, enableResponseLog, &CommonOpenAIChatCompletionsAPILargeLanguageModelAdapter{
		apiProvider: apiProvider,
	})
}
