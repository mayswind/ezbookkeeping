package ollama

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/llm/data"
	"github.com/mayswind/ezbookkeeping/pkg/llm/provider"
	"github.com/mayswind/ezbookkeeping/pkg/llm/provider/common"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

const ollamaChatCompletionsPath = "api/chat"

// OllamaLargeLanguageModelAdapter defines the structure of Ollama large language model adapter
type OllamaLargeLanguageModelAdapter struct {
	common.HttpLargeLanguageModelAdapter
	OllamaServerURL string
	OllamaModelID   string
}

// OllamaMessageRole defines the role of Ollama chat message
type OllamaMessageRole string

const (
	OllamaMessageRoleSystem OllamaMessageRole = "system"
	OllamaMessageRoleUser   OllamaMessageRole = "user"
)

// OllamaChatRequest defines the structure of Ollama chat request
type OllamaChatRequest struct {
	Model    string                      `json:"model"`
	Stream   bool                        `json:"stream"`
	Messages []*OllamaChatRequestMessage `json:"messages"`
	Format   string                      `json:"format,omitempty"`
}

// OllamaChatRequestMessage defines the structure of Ollama chat request message
type OllamaChatRequestMessage struct {
	Role    OllamaMessageRole `json:"role"`
	Content string            `json:"content"`
	Images  []string          `json:"images,omitempty"`
}

// OllamaChatResponse defines the structure of Ollama chat response
type OllamaChatResponse struct {
	Message *OllamaChatResponseMessage `json:"message"`
}

// OllamaChatResponseMessage defines the structure of Ollama chat response message
type OllamaChatResponseMessage struct {
	Content *string `json:"content"`
}

// BuildTextualRequest returns the http request by Ollama large language model adapter
func (p *OllamaLargeLanguageModelAdapter) BuildTextualRequest(c core.Context, uid int64, request *data.LargeLanguageModelRequest, responseType data.LargeLanguageModelResponseFormat) (*http.Request, error) {
	requestBody, err := p.buildJsonRequestBody(c, uid, request, responseType)

	if err != nil {
		return nil, err
	}

	httpRequest, err := http.NewRequest("POST", p.getOllamaRequestUrl(), bytes.NewReader(requestBody))

	if err != nil {
		return nil, err
	}

	httpRequest.Header.Set("Content-Type", "application/json")

	return httpRequest, nil
}

// ParseTextualResponse returns the textual response by Ollama large language model adapter
func (p *OllamaLargeLanguageModelAdapter) ParseTextualResponse(c core.Context, uid int64, body []byte, responseType data.LargeLanguageModelResponseFormat) (*data.LargeLanguageModelTextualResponse, error) {
	chatResponse := &OllamaChatResponse{}
	err := json.Unmarshal(body, &chatResponse)

	if err != nil {
		log.Errorf(c, "[ollama_large_language_model_adapter.ParseTextualResponse] failed to parse chat response for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	if chatResponse == nil || chatResponse.Message == nil || chatResponse.Message.Content == nil {
		log.Errorf(c, "[ollama_large_language_model_adapter.ParseTextualResponse] chat response is invalid for user \"uid:%d\"", uid)
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	textualResponse := &data.LargeLanguageModelTextualResponse{
		Content: *chatResponse.Message.Content,
	}

	return textualResponse, nil
}

func (p *OllamaLargeLanguageModelAdapter) buildJsonRequestBody(c core.Context, uid int64, request *data.LargeLanguageModelRequest, responseType data.LargeLanguageModelResponseFormat) ([]byte, error) {
	if p.OllamaModelID == "" {
		return nil, errs.ErrInvalidLLMModelId
	}

	chatRequest := &OllamaChatRequest{
		Model:    p.OllamaModelID,
		Stream:   request.Stream,
		Messages: make([]*OllamaChatRequestMessage, 0, 2),
	}

	if request.SystemPrompt != "" {
		chatRequest.Messages = append(chatRequest.Messages, &OllamaChatRequestMessage{
			Role:    OllamaMessageRoleSystem,
			Content: request.SystemPrompt,
		})
	}

	if len(request.UserPrompt) > 0 {
		if request.UserPromptType == data.LARGE_LANGUAGE_MODEL_REQUEST_PROMPT_TYPE_IMAGE_URL {
			imageBase64Data := base64.StdEncoding.EncodeToString(request.UserPrompt)
			chatRequest.Messages = append(chatRequest.Messages, &OllamaChatRequestMessage{
				Role:   OllamaMessageRoleUser,
				Images: []string{imageBase64Data},
			})
		} else {
			chatRequest.Messages = append(chatRequest.Messages, &OllamaChatRequestMessage{
				Role:    OllamaMessageRoleUser,
				Content: string(request.UserPrompt),
			})
		}
	}

	if responseType == data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON {
		chatRequest.Format = "json"
	}

	requestBodyBytes, err := json.Marshal(chatRequest)

	if err != nil {
		log.Errorf(c, "[ollama_large_language_model_adapter.buildJsonRequestBody] failed to marshal request body for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	log.Debugf(c, "[ollama_large_language_model_adapter.buildJsonRequestBody] request body is %s", requestBodyBytes)
	return requestBodyBytes, nil
}

func (p *OllamaLargeLanguageModelAdapter) getOllamaRequestUrl() string {
	url := p.OllamaServerURL

	if url[len(url)-1] != '/' {
		url += "/"
	}

	url += ollamaChatCompletionsPath
	return url
}

// NewOllamaLargeLanguageModelProvider creates a new Ollama large language model provider instance
func NewOllamaLargeLanguageModelProvider(llmConfig *settings.LLMConfig) provider.LargeLanguageModelProvider {
	return common.NewCommonHttpLargeLanguageModelProvider(llmConfig, &OllamaLargeLanguageModelAdapter{
		OllamaServerURL: llmConfig.OllamaServerURL,
		OllamaModelID:   llmConfig.OllamaModelID,
	})
}
