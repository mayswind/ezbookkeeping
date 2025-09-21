package llm

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

const ollamaChatCompletionsPath = "api/chat"

// OllamaLargeLanguageModelProvider defines the structure of Ollama large language model provider
type OllamaLargeLanguageModelProvider struct {
	CommonHttpLargeLanguageModelProvider
	OllamaServerURL string
	OllamaModelID   string
}

// BuildTextualRequest returns the http request by Ollama provider
func (p *OllamaLargeLanguageModelProvider) BuildTextualRequest(c core.Context, uid int64, request *LargeLanguageModelRequest, responseType LargeLanguageModelResponseFormat) (*http.Request, error) {
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

// ParseTextualResponse returns the textual response by Ollama provider
func (p *OllamaLargeLanguageModelProvider) ParseTextualResponse(c core.Context, uid int64, body []byte, responseType LargeLanguageModelResponseFormat) (*LargeLanguageModelTextualResponse, error) {
	responseBody := make(map[string]any)
	err := json.Unmarshal(body, &responseBody)

	if err != nil {
		log.Errorf(c, "[ollama_large_language_model_provider.ParseTextualResponse] failed to parse response for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	message, ok := responseBody["message"].(map[string]any)

	if !ok {
		log.Errorf(c, "[ollama_large_language_model_provider.ParseTextualResponse] no message found in response for user \"uid:%d\"", uid)
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	content, ok := message["content"].(string)

	if !ok {
		log.Errorf(c, "[ollama_large_language_model_provider.ParseTextualResponse] no content found in message for user \"uid:%d\"", uid)
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

// GetModelID returns the model id of Ollama provider
func (p *OllamaLargeLanguageModelProvider) GetModelID() string {
	return p.OllamaModelID
}

func (p *OllamaLargeLanguageModelProvider) buildJsonRequestBody(c core.Context, uid int64, request *LargeLanguageModelRequest, responseType LargeLanguageModelResponseFormat) ([]byte, error) {
	if p.OllamaModelID == "" {
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
		imageBase64Data := base64.StdEncoding.EncodeToString(request.UserPrompt)
		if request.UserPromptType == LARGE_LANGUAGE_MODEL_REQUEST_PROMPT_TYPE_IMAGE_URL {
			requestMessages = append(requestMessages, map[string]any{
				"role":    "user",
				"content": "",
				"images":  []string{imageBase64Data},
			})
		} else {
			requestMessages = append(requestMessages, map[string]string{
				"role":    "user",
				"content": string(request.UserPrompt),
			})
		}
	}

	requestBody := make(map[string]any)
	requestBody["model"] = p.OllamaModelID
	requestBody["stream"] = request.Stream
	requestBody["messages"] = requestMessages

	if responseType == LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON {
		requestBody["format"] = "json"
	}

	requestBodyBytes, err := json.Marshal(requestBody)

	if err != nil {
		log.Errorf(c, "[ollama_large_language_model_provider.buildJsonRequestBody] failed to marshal request body for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	log.Debugf(c, "[ollama_large_language_model_provider.buildJsonRequestBody] request body is %s", requestBodyBytes)
	return requestBodyBytes, nil
}

func (p *OllamaLargeLanguageModelProvider) getOllamaRequestUrl() string {
	url := p.OllamaServerURL

	if url[len(url)-1] != '/' {
		url += "/"
	}

	url += ollamaChatCompletionsPath
	return url
}

// NewOllamaLargeLanguageModelProvider creates a new Ollama large language model provider instance
func NewOllamaLargeLanguageModelProvider(llmConfig *settings.LLMConfig) LargeLanguageModelProvider {
	return newCommonHttpLargeLanguageModelProvider(&OllamaLargeLanguageModelProvider{
		OllamaServerURL: llmConfig.OllamaServerURL,
		OllamaModelID:   llmConfig.OllamaModelID,
	})
}
