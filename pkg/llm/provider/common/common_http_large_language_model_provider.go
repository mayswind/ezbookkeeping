package common

import (
	"io"
	"net/http"
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/llm/data"
	"github.com/mayswind/ezbookkeeping/pkg/llm/provider"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// HttpLargeLanguageModelAdapter defines the structure of http large language model adapter
type HttpLargeLanguageModelAdapter interface {
	// BuildTextualRequest returns the http request by the provider api definition
	BuildTextualRequest(c core.Context, uid int64, request *data.LargeLanguageModelRequest, responseType data.LargeLanguageModelResponseFormat) (*http.Request, error)

	// ParseTextualResponse returns the textual response entity by the provider api definition
	ParseTextualResponse(c core.Context, uid int64, body []byte, responseType data.LargeLanguageModelResponseFormat) (*data.LargeLanguageModelTextualResponse, error)
}

// CommonHttpLargeLanguageModelProvider defines the structure of common http large language model provider
type CommonHttpLargeLanguageModelProvider struct {
	provider.LargeLanguageModelProvider
	adapter    HttpLargeLanguageModelAdapter
	httpClient *http.Client
}

// GetJsonResponse returns the json response from common http large language model provider
func (p *CommonHttpLargeLanguageModelProvider) GetJsonResponse(c core.Context, uid int64, currentLLMConfig *settings.LLMConfig, request *data.LargeLanguageModelRequest) (*data.LargeLanguageModelTextualResponse, error) {
	response, err := p.getTextualResponse(c, uid, currentLLMConfig, request, data.LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)

	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(response.Content, "```json") && strings.HasSuffix(response.Content, "```") {
		response.Content = strings.TrimPrefix(response.Content, "```json")
		response.Content = strings.TrimSuffix(response.Content, "```")
	} else if strings.HasPrefix(response.Content, "```") && strings.HasSuffix(response.Content, "```") {
		response.Content = strings.TrimPrefix(response.Content, "```")
		response.Content = strings.TrimSuffix(response.Content, "```")
	}

	return response, nil
}

func (p *CommonHttpLargeLanguageModelProvider) getTextualResponse(c core.Context, uid int64, currentLLMConfig *settings.LLMConfig, request *data.LargeLanguageModelRequest, responseType data.LargeLanguageModelResponseFormat) (*data.LargeLanguageModelTextualResponse, error) {
	httpRequest, err := p.adapter.BuildTextualRequest(c, uid, request, responseType)

	if err != nil {
		log.Errorf(c, "[common_http_large_language_model_provider.getTextualResponse] failed to build requests for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	resp, err := p.httpClient.Do(httpRequest)

	if err != nil {
		log.Errorf(c, "[common_http_large_language_model_provider.getTextualResponse] failed to request large language model api for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	log.Debugf(c, "[common_http_large_language_model_provider.getTextualResponse] response is %s", body)

	if resp.StatusCode != 200 {
		log.Errorf(c, "[common_http_large_language_model_provider.getTextualResponse] failed to get large language model api response for user \"uid:%d\", because response code is %d", uid, resp.StatusCode)
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	return p.adapter.ParseTextualResponse(c, uid, body, responseType)
}

// NewCommonHttpLargeLanguageModelProvider creates a http adapter based large language model provider instance
func NewCommonHttpLargeLanguageModelProvider(llmConfig *settings.LLMConfig, adapter HttpLargeLanguageModelAdapter) *CommonHttpLargeLanguageModelProvider {
	return &CommonHttpLargeLanguageModelProvider{
		adapter:    adapter,
		httpClient: utils.NewHttpClient(llmConfig.LargeLanguageModelAPIRequestTimeout, llmConfig.LargeLanguageModelAPIProxy, llmConfig.LargeLanguageModelAPISkipTLSVerify, settings.GetUserAgent()),
	}
}
