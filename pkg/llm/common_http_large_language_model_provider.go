package llm

import (
	"crypto/tls"
	"io"
	"net/http"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// HttpLargeLanguageModelAdapter defines the structure of http large language model adapter
type HttpLargeLanguageModelAdapter interface {
	// BuildTextualRequest returns the http request by the provider api definition
	BuildTextualRequest(c core.Context, uid int64, request *LargeLanguageModelRequest, responseType LargeLanguageModelResponseFormat) (*http.Request, error)

	// ParseTextualResponse returns the textual response entity by the provider api definition
	ParseTextualResponse(c core.Context, uid int64, body []byte, responseType LargeLanguageModelResponseFormat) (*LargeLanguageModelTextualResponse, error)
}

// CommonHttpLargeLanguageModelProvider defines the structure of common http large language model provider
type CommonHttpLargeLanguageModelProvider struct {
	LargeLanguageModelProvider
	adapter HttpLargeLanguageModelAdapter
}

// GetJsonResponse returns the json response from the OpenAI common compatible large language model provider
func (p *CommonHttpLargeLanguageModelProvider) GetJsonResponse(c core.Context, uid int64, currentLLMConfig *settings.LLMConfig, request *LargeLanguageModelRequest) (*LargeLanguageModelTextualResponse, error) {
	return p.getTextualResponse(c, uid, currentLLMConfig, request, LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
}

func (p *CommonHttpLargeLanguageModelProvider) getTextualResponse(c core.Context, uid int64, currentLLMConfig *settings.LLMConfig, request *LargeLanguageModelRequest, responseType LargeLanguageModelResponseFormat) (*LargeLanguageModelTextualResponse, error) {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	utils.SetProxyUrl(transport, currentLLMConfig.LargeLanguageModelAPIProxy)

	if currentLLMConfig.LargeLanguageModelAPISkipTLSVerify {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(currentLLMConfig.LargeLanguageModelAPIRequestTimeout) * time.Millisecond,
	}

	httpRequest, err := p.adapter.BuildTextualRequest(c, uid, request, responseType)

	if err != nil {
		log.Errorf(c, "[common_http_large_language_model_provider.getTextualResponse] failed to build requests for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	httpRequest.Header.Set("User-Agent", settings.GetUserAgent())

	resp, err := client.Do(httpRequest)

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

func newCommonHttpLargeLanguageModelProvider(adapter HttpLargeLanguageModelAdapter) *CommonHttpLargeLanguageModelProvider {
	return &CommonHttpLargeLanguageModelProvider{
		adapter: adapter,
	}
}
