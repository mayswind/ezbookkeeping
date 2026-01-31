package lmstudio

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

const lmStudioChatPath = "api/v1/chat"

// LMStudioLargeLanguageModelAdapter defines the structure of LM Studio large language model adapter
type LMStudioLargeLanguageModelAdapter struct {
	common.HttpLargeLanguageModelAdapter
	LMStudioServerURL string
	LMStudioToken     string
	LMStudioModelID   string
}

// LMStudioChatRequest defines the structure of LM Studio chat request
type LMStudioChatRequest struct {
	Model        string                      `json:"model"`
	SystemPrompt string                      `json:"system_prompt,omitempty"`
	Input        []*LMStudioChatRequestInput `json:"input"`
}

// LMStudioChatRequestInput defines the structure of LM Studio chat request message
type LMStudioChatRequestInput struct {
	Type    string `json:"type"`
	Content string `json:"content,omitempty"`
	DataUrl string `json:"data_url,omitempty"`
}

// LMStudioChatResponse defines the structure of LM Studio chat response
type LMStudioChatResponse struct {
	Output []*LMStudioChatResponseOutput `json:"output"`
}

// LMStudioChatResponseOutput defines the structure of LM Studio chat response message
type LMStudioChatResponseOutput struct {
	Content *string `json:"content"`
}

// BuildTextualRequest returns the http request by LM Studio large language model adapter
func (p *LMStudioLargeLanguageModelAdapter) BuildTextualRequest(c core.Context, uid int64, request *data.LargeLanguageModelRequest, responseType data.LargeLanguageModelResponseFormat) (*http.Request, error) {
	requestBody, err := p.buildJsonRequestBody(c, uid, request, responseType)

	if err != nil {
		return nil, err
	}

	httpRequest, err := http.NewRequest("POST", p.getLMStudioRequestUrl(), bytes.NewReader(requestBody))

	if err != nil {
		return nil, err
	}

	if p.LMStudioToken != "" {
		httpRequest.Header.Set("Authorization", "Bearer "+p.LMStudioToken)
	}

	httpRequest.Header.Set("Content-Type", "application/json")

	return httpRequest, nil
}

// ParseTextualResponse returns the textual response by LM Studio large language model adapter
func (p *LMStudioLargeLanguageModelAdapter) ParseTextualResponse(c core.Context, uid int64, body []byte, responseType data.LargeLanguageModelResponseFormat) (*data.LargeLanguageModelTextualResponse, error) {
	chatResponse := &LMStudioChatResponse{}
	err := json.Unmarshal(body, &chatResponse)

	if err != nil {
		log.Errorf(c, "[lm_studio_large_language_model_adapter.ParseTextualResponse] failed to parse chat response for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	if chatResponse == nil || len(chatResponse.Output) < 1 || chatResponse.Output[0].Content == nil {
		log.Errorf(c, "[lm_studio_large_language_model_adapter.ParseTextualResponse] chat response is invalid for user \"uid:%d\"", uid)
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	textualResponse := &data.LargeLanguageModelTextualResponse{
		Content: *chatResponse.Output[0].Content,
	}

	return textualResponse, nil
}

func (p *LMStudioLargeLanguageModelAdapter) buildJsonRequestBody(c core.Context, uid int64, request *data.LargeLanguageModelRequest, responseType data.LargeLanguageModelResponseFormat) ([]byte, error) {
	if p.LMStudioModelID == "" {
		return nil, errs.ErrInvalidLLMModelId
	}

	chatRequest := &LMStudioChatRequest{
		Model: p.LMStudioModelID,
		Input: make([]*LMStudioChatRequestInput, 0, 1),
	}

	if request.SystemPrompt != "" {
		chatRequest.SystemPrompt = request.SystemPrompt
	}

	if len(request.UserPrompt) > 0 {
		if request.UserPromptType == data.LARGE_LANGUAGE_MODEL_REQUEST_PROMPT_TYPE_IMAGE_URL {
			imageBase64Data := "data:" + request.UserPromptContentType + ";base64," + base64.StdEncoding.EncodeToString(request.UserPrompt)
			chatRequest.Input = append(chatRequest.Input, &LMStudioChatRequestInput{
				Type:    "image",
				DataUrl: imageBase64Data,
			})
		} else {
			chatRequest.Input = append(chatRequest.Input, &LMStudioChatRequestInput{
				Type:    "text",
				Content: string(request.UserPrompt),
			})
		}
	}

	requestBodyBytes, err := json.Marshal(chatRequest)

	if err != nil {
		log.Errorf(c, "[lm_studio_large_language_model_adapter.buildJsonRequestBody] failed to marshal request body for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	log.Debugf(c, "[lm_studio_large_language_model_adapter.buildJsonRequestBody] request body is %s", requestBodyBytes)
	return requestBodyBytes, nil
}

func (p *LMStudioLargeLanguageModelAdapter) getLMStudioRequestUrl() string {
	url := p.LMStudioServerURL

	if url[len(url)-1] != '/' {
		url += "/"
	}

	url += lmStudioChatPath
	return url
}

// NewLMStudioLargeLanguageModelProvider creates a new LM Studio large language model provider instance
func NewLMStudioLargeLanguageModelProvider(llmConfig *settings.LLMConfig, enableResponseLog bool) provider.LargeLanguageModelProvider {
	return common.NewCommonHttpLargeLanguageModelProvider(llmConfig, enableResponseLog, &LMStudioLargeLanguageModelAdapter{
		LMStudioServerURL: llmConfig.LMStudioServerURL,
		LMStudioToken:     llmConfig.LMStudioToken,
		LMStudioModelID:   llmConfig.LMStudioModelID,
	})
}
