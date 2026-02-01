package anthropic

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/llm/data"
	"github.com/mayswind/ezbookkeeping/pkg/llm/provider"
	"github.com/mayswind/ezbookkeeping/pkg/llm/provider/common"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// AnthropicMessagesAPIProvider defines the structure of Anthropic messages API provider
type AnthropicMessagesAPIProvider interface {
	// BuildMessagesHttpRequest returns the messages http request
	BuildMessagesHttpRequest(c core.Context, uid int64) (*http.Request, error)

	// GetModelID returns the model id
	GetModelID() string

	// GetMaxTokens returns the max tokens to generate
	GetMaxTokens() uint32
}

// CommonAnthropicMessagesAPILargeLanguageModelAdapter defines the structure of Anthropic common compatible large language model adapter based on messages api
type CommonAnthropicMessagesAPILargeLanguageModelAdapter struct {
	common.HttpLargeLanguageModelAdapter
	apiProvider AnthropicMessagesAPIProvider
}

// AnthropicMessageRole defines the role of Anthropic message
type AnthropicMessageRole string

// Anthropic Message Roles
const (
	AnthropicMessageRoleUser AnthropicMessageRole = "user"
)

type AnthropicThinkingType string

// Anthropic Thinking Types
const (
	AnthropicThinkingTypeDisabled AnthropicThinkingType = "disabled"
)

// AnthropicMessagesRequest defines the structure of Anthropic messages request
type AnthropicMessagesRequest struct {
	Model     string                                       `json:"model"`
	MaxTokens uint32                                       `json:"max_tokens"`
	Stream    bool                                         `json:"stream"`
	System    string                                       `json:"system,omitempty"`
	Messages  []any                                        `json:"messages"`
	Thinking  *AnthropicMessagesRequestThinkingConfigParam `json:"thinking,omitempty"`
}

// AnthropicMessagesRequestMessage defines the structure of Anthropic messages request message
type AnthropicMessagesRequestMessage[T string | []*AnthropicMessagesRequestImageBlockParam] struct {
	Role    AnthropicMessageRole `json:"role"`
	Content T                    `json:"content"`
}

// AnthropicMessagesRequestImageBlockParam defines the structure of Anthropic messages request image content block param
type AnthropicMessagesRequestImageBlockParam struct {
	Source *AnthropicMessagesRequestBase64ImageSource `json:"source"`
	Type   string                                     `json:"type"`
}

// AnthropicMessagesRequestBase64ImageSource defines the structure of Anthropic messages request base64 image source
type AnthropicMessagesRequestBase64ImageSource struct {
	Data      string `json:"data"`
	MediaType string `json:"media_type"`
	Type      string `json:"type"`
}

// AnthropicMessagesRequestThinkingConfigParam defines the structure of Anthropic messages request thinking config param
type AnthropicMessagesRequestThinkingConfigParam struct {
	Type AnthropicThinkingType `json:"type"`
}

// AnthropicMessagesResponse defines the structure of Anthropic messages response
type AnthropicMessagesResponse struct {
	Content []*AnthropicMessagesResponseContentBlock `json:"content"`
}

// AnthropicMessagesResponseContentBlock defines the structure of Anthropic messages response content block
type AnthropicMessagesResponseContentBlock struct {
	Text *string `json:"text"`
}

// BuildTextualRequest returns the http request by Anthropic common compatible adapter
func (p *CommonAnthropicMessagesAPILargeLanguageModelAdapter) BuildTextualRequest(c core.Context, uid int64, request *data.LargeLanguageModelRequest, responseType data.LargeLanguageModelResponseFormat) (*http.Request, error) {
	requestBody, err := p.buildJsonRequestBody(c, uid, request, responseType)

	if err != nil {
		return nil, err
	}

	httpRequest, err := p.apiProvider.BuildMessagesHttpRequest(c, uid)

	if err != nil {
		return nil, err
	}

	httpRequest.Body = io.NopCloser(bytes.NewReader(requestBody))
	httpRequest.Header.Set("Content-Type", "application/json")

	return httpRequest, nil
}

// ParseTextualResponse returns the textual response by Anthropic common compatible adapter
func (p *CommonAnthropicMessagesAPILargeLanguageModelAdapter) ParseTextualResponse(c core.Context, uid int64, body []byte, responseType data.LargeLanguageModelResponseFormat) (*data.LargeLanguageModelTextualResponse, error) {
	messagesResponse := &AnthropicMessagesResponse{}
	err := json.Unmarshal(body, &messagesResponse)

	if err != nil {
		log.Errorf(c, "[anthropic_common_compatible_large_language_model_adapter.ParseTextualResponse] failed to parse messages response for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	if messagesResponse == nil || messagesResponse.Content == nil || len(messagesResponse.Content) < 1 || messagesResponse.Content[0].Text == nil {
		log.Errorf(c, "[anthropic_common_compatible_large_language_model_adapter.ParseTextualResponse] messages response is invalid for user \"uid:%d\"", uid)
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	textualResponse := &data.LargeLanguageModelTextualResponse{
		Content: *messagesResponse.Content[0].Text,
	}

	return textualResponse, nil
}

func (p *CommonAnthropicMessagesAPILargeLanguageModelAdapter) buildJsonRequestBody(c core.Context, uid int64, request *data.LargeLanguageModelRequest, responseType data.LargeLanguageModelResponseFormat) ([]byte, error) {
	if p.apiProvider.GetModelID() == "" {
		return nil, errs.ErrInvalidLLMModelId
	}

	messagesRequest := &AnthropicMessagesRequest{
		Model:     p.apiProvider.GetModelID(),
		MaxTokens: p.apiProvider.GetMaxTokens(),
		Stream:    request.Stream,
		Messages:  make([]any, 0, 1),
		Thinking: &AnthropicMessagesRequestThinkingConfigParam{
			Type: AnthropicThinkingTypeDisabled,
		},
	}

	if request.SystemPrompt != "" {
		messagesRequest.System = request.SystemPrompt
	}

	if len(request.UserPrompt) > 0 {
		if request.UserPromptType == data.LARGE_LANGUAGE_MODEL_REQUEST_PROMPT_TYPE_IMAGE_URL {
			imageBase64Data := base64.StdEncoding.EncodeToString(request.UserPrompt)
			messagesRequest.Messages = append(messagesRequest.Messages, &AnthropicMessagesRequestMessage[[]*AnthropicMessagesRequestImageBlockParam]{
				Role: AnthropicMessageRoleUser,
				Content: []*AnthropicMessagesRequestImageBlockParam{
					{
						Type: "image",
						Source: &AnthropicMessagesRequestBase64ImageSource{
							Data:      imageBase64Data,
							MediaType: request.UserPromptContentType,
							Type:      "base64",
						},
					},
				},
			})
		} else {
			messagesRequest.Messages = append(messagesRequest.Messages, &AnthropicMessagesRequestMessage[string]{
				Role:    AnthropicMessageRoleUser,
				Content: string(request.UserPrompt),
			})
		}
	}

	requestBodyBytes, err := json.Marshal(messagesRequest)

	if err != nil {
		log.Errorf(c, "[anthropic_common_compatible_large_language_model_adapter.buildJsonRequestBody] failed to marshal request body for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	log.Debugf(c, "[anthropic_common_compatible_large_language_model_adapter.buildJsonRequestBody] request body is %s", requestBodyBytes)
	return requestBodyBytes, nil
}

func newCommonAnthropicMessagesAPILargeLanguageModelAdapter(llmConfig *settings.LLMConfig, enableResponseLog bool, apiProvider AnthropicMessagesAPIProvider) provider.LargeLanguageModelProvider {
	return common.NewCommonHttpLargeLanguageModelProvider(llmConfig, enableResponseLog, &CommonAnthropicMessagesAPILargeLanguageModelAdapter{
		apiProvider: apiProvider,
	})
}
