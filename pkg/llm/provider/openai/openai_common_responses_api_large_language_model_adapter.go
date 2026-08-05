package openai

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
	"github.com/mayswind/ezbookkeeping/pkg/llm/data"
	"github.com/mayswind/ezbookkeeping/pkg/llm/provider"
	"github.com/mayswind/ezbookkeeping/pkg/llm/provider/common"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// OpenAIResponsesAPIProvider defines the structure of OpenAI responses API provider
type OpenAIResponsesAPIProvider interface {
	// BuildResponsesHttpRequest returns the responses http request
	BuildResponsesHttpRequest(c core.Context, uid int64) (*http.Request, error)

	// GetModelID returns the model id if supported, otherwise returns empty string
	GetModelID() string
}

// CommonOpenAIResponsesAPILargeLanguageModelAdapter defines the structure of OpenAI common compatible large language model adapter based on responses api
type CommonOpenAIResponsesAPILargeLanguageModelAdapter struct {
	common.HttpLargeLanguageModelAdapter
	apiProvider   OpenAIResponsesAPIProvider
	ThinkingLevel settings.LLMThinkingLevel
}

// OpenAIResponsesRequestResponseFormatType defines the type of OpenAI responses request response format
type OpenAIResponsesRequestResponseFormatType string

// OpenAI Responses Request Response Format Types
const (
	OpenAIResponsesRequestResponseFormatTypeJsonObject OpenAIResponsesRequestResponseFormatType = "json_object"
	OpenAIResponsesRequestResponseFormatTypeJsonSchema OpenAIResponsesRequestResponseFormatType = "json_schema"
)

// OpenAIResponsesRequestReasoning defines the reasoning structure of OpenAI responses request
type OpenAIResponsesRequestReasoning struct {
	Effort string `json:"effort"`
}

// OpenAIResponsesRequest defines the structure of OpenAI responses request
type OpenAIResponsesRequest struct {
	Model     string                            `json:"model"`
	Stream    bool                              `json:"stream"`
	Input     []any                             `json:"input"`
	Reasoning *OpenAIResponsesRequestReasoning  `json:"reasoning,omitempty"`
	Text      *OpenAIResponsesRequestTextConfig `json:"text,omitempty"`
}

// OpenAIResponsesRequestInputMessage defines the structure of an OpenAI responses input message
type OpenAIResponsesRequestInputMessage[T string | []*OpenAIResponsesRequestInputContent] struct {
	Role    OpenAIMessageRole `json:"role"`
	Content T                 `json:"content"`
}

// OpenAIResponsesRequestInputContent defines the structure of OpenAI responses input content
type OpenAIResponsesRequestInputContent struct {
	Type     string `json:"type"`
	ImageURL string `json:"image_url"`
}

// OpenAIResponsesRequestTextConfig defines the text configuration of an OpenAI responses request
type OpenAIResponsesRequestTextConfig struct {
	Format *OpenAIResponsesRequestResponseFormat `json:"format"`
}

// OpenAIResponsesRequestResponseFormat defines the response format of an OpenAI responses request
type OpenAIResponsesRequestResponseFormat struct {
	Type   OpenAIResponsesRequestResponseFormatType `json:"type"`
	Name   string                                   `json:"name,omitempty"`
	Strict bool                                     `json:"strict,omitempty"`
	Schema *jsonschema.Schema                       `json:"schema,omitempty"`
}

// OpenAIResponsesResponse defines the structure of an OpenAI responses response
type OpenAIResponsesResponse struct {
	Output []*OpenAIResponsesResponseOutput `json:"output"`
}

// OpenAIResponsesResponseOutput defines the structure of an OpenAI responses output item
type OpenAIResponsesResponseOutput struct {
	Type    string                                  `json:"type"`
	Content []*OpenAIResponsesResponseOutputContent `json:"content"`
}

// OpenAIResponsesResponseOutputContent defines the structure of OpenAI responses output content
type OpenAIResponsesResponseOutputContent struct {
	Type string  `json:"type"`
	Text *string `json:"text"`
}

// OpenAI Responses Request Reasoning Efforts Mapping
var openAIResponsesRequestReasoningEffortsMapping = map[settings.LLMThinkingLevel]string{
	settings.LLMThinkingDisabled: "none",
	settings.LLMThinkingLow:      "low",
	settings.LLMThinkingMedium:   "medium",
	settings.LLMThinkingEnabled:  "medium",
	settings.LLMThinkingHigh:     "high",
	settings.LLMThinkingXHigh:    "xhigh",
}

// BuildTextualRequest returns the http request by OpenAI common responses adapter
func (p *CommonOpenAIResponsesAPILargeLanguageModelAdapter) BuildTextualRequest(c core.Context, uid int64, request *data.LargeLanguageModelRequest, responseType data.LargeLanguageModelResponseFormat) (*http.Request, error) {
	requestBody, err := p.buildJsonRequestBody(c, uid, request, responseType)

	if err != nil {
		return nil, err
	}

	httpRequest, err := p.apiProvider.BuildResponsesHttpRequest(c, uid)

	if err != nil {
		return nil, err
	}

	httpRequest.Body = io.NopCloser(bytes.NewReader(requestBody))
	httpRequest.Header.Set("Content-Type", "application/json")

	return httpRequest, nil
}

// ParseTextualResponse returns the textual response by OpenAI common responses adapter
func (p *CommonOpenAIResponsesAPILargeLanguageModelAdapter) ParseTextualResponse(c core.Context, uid int64, body []byte, responseType data.LargeLanguageModelResponseFormat) (*data.LargeLanguageModelTextualResponse, error) {
	responsesResponse := &OpenAIResponsesResponse{}
	err := json.Unmarshal(body, responsesResponse)

	if err != nil {
		log.Errorf(c, "[openai_common_responses_api_large_language_model_adapter.ParseTextualResponse] failed to parse responses response for user uid:%d, because %s", uid, err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	var content strings.Builder
	hasOutputText := false

	for _, output := range responsesResponse.Output {
		if output == nil || output.Type != "message" {
			continue
		}

		for _, outputContent := range output.Content {
			if outputContent == nil || outputContent.Type != "output_text" || outputContent.Text == nil {
				continue
			}

			content.WriteString(*outputContent.Text)
			hasOutputText = true
		}
	}

	if !hasOutputText {
		log.Errorf(c, "[openai_common_responses_api_large_language_model_adapter.ParseTextualResponse] responses response is invalid for user uid:%d", uid)
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	return &data.LargeLanguageModelTextualResponse{
		Content: content.String(),
	}, nil
}

func (p *CommonOpenAIResponsesAPILargeLanguageModelAdapter) buildJsonRequestBody(c core.Context, uid int64, request *data.LargeLanguageModelRequest, responseType data.LargeLanguageModelResponseFormat) ([]byte, error) {
	if p.apiProvider.GetModelID() == "" {
		return nil, errs.ErrInvalidLLMModelId
	}

	responsesRequest := &OpenAIResponsesRequest{
		Model:  p.apiProvider.GetModelID(),
		Stream: request.Stream,
		Input:  make([]any, 0, 2),
	}

	if thinkingLevel, exists := openAIResponsesRequestReasoningEffortsMapping[p.ThinkingLevel]; exists {
		responsesRequest.Reasoning = &OpenAIResponsesRequestReasoning{
			Effort: thinkingLevel,
		}
	}

	if request.SystemPrompt != "" {
		responsesRequest.Input = append(responsesRequest.Input, &OpenAIResponsesRequestInputMessage[string]{
			Role:    OpenAIMessageRoleSystem,
			Content: request.SystemPrompt,
		})
	}

	if len(request.UserPrompt) > 0 {
		if request.UserPromptType == data.LARGE_LANGUAGE_MODEL_REQUEST_PROMPT_TYPE_IMAGE_URL {
			imageBase64Data := "data:" + request.UserPromptContentType + ";base64," + base64.StdEncoding.EncodeToString(request.UserPrompt)
			responsesRequest.Input = append(responsesRequest.Input, &OpenAIResponsesRequestInputMessage[[]*OpenAIResponsesRequestInputContent]{
				Role: OpenAIMessageRoleUser,
				Content: []*OpenAIResponsesRequestInputContent{
					{
						Type:     "input_image",
						ImageURL: imageBase64Data,
					},
				},
			})
		} else {
			responsesRequest.Input = append(responsesRequest.Input, &OpenAIResponsesRequestInputMessage[string]{
				Role:    OpenAIMessageRoleUser,
				Content: string(request.UserPrompt),
			})
		}
	}

	if responseType == data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON {
		if request.ResponseJsonObjectType != nil {
			schemaGenerator := jsonschema.Reflector{
				Anonymous:      true,
				DoNotReference: true,
				ExpandedStruct: true,
			}

			schema := schemaGenerator.ReflectFromType(request.ResponseJsonObjectType)
			schema.Version = ""

			responsesRequest.Text = &OpenAIResponsesRequestTextConfig{
				Format: &OpenAIResponsesRequestResponseFormat{
					Type:   OpenAIResponsesRequestResponseFormatTypeJsonSchema,
					Name:   "response",
					Strict: true,
					Schema: schema,
				},
			}
		} else {
			responsesRequest.Text = &OpenAIResponsesRequestTextConfig{
				Format: &OpenAIResponsesRequestResponseFormat{
					Type: OpenAIResponsesRequestResponseFormatTypeJsonObject,
				},
			}
		}
	}

	requestBodyBytes, err := json.Marshal(responsesRequest)

	if err != nil {
		log.Errorf(c, "[openai_common_responses_api_large_language_model_adapter.buildJsonRequestBody] failed to marshal request body for user uid:%d, because %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	log.Debugf(c, "[openai_common_responses_api_large_language_model_adapter.buildJsonRequestBody] request body is %s", requestBodyBytes)
	return requestBodyBytes, nil
}

func newCommonOpenAIResponsesAPILargeLanguageModelAdapter(llmConfig *settings.LLMConfig, enableResponseLog bool, apiProvider OpenAIResponsesAPIProvider) provider.LargeLanguageModelProvider {
	return common.NewCommonHttpLargeLanguageModelProvider(llmConfig, enableResponseLog, &CommonOpenAIResponsesAPILargeLanguageModelAdapter{
		apiProvider:   apiProvider,
		ThinkingLevel: llmConfig.EnableThinking,
	})
}
