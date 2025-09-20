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

// HttpLargeLanguageModelProvider defines the structure of http large language model provider
type HttpLargeLanguageModelProvider interface {
	// BuildTextualRequest returns the http request by the provider api definition
	BuildTextualRequest(c core.Context, uid int64, request *LargeLanguageModelRequest, modelId string, responseType LargeLanguageModelResponseFormat) (*http.Request, error)

	// ParseTextualResponse returns the textual response entity by the provider api definition
	ParseTextualResponse(c core.Context, uid int64, body []byte, responseType LargeLanguageModelResponseFormat) (*LargeLanguageModelTextualResponse, error)

	// GetReceiptImageRecognitionModelID returns the receipt image recognition model id if supported, otherwise returns empty string
	GetReceiptImageRecognitionModelID() string
}

// CommonHttpLargeLanguageModelProvider defines the structure of common http large language model provider
type CommonHttpLargeLanguageModelProvider struct {
	LargeLanguageModelProvider
	provider HttpLargeLanguageModelProvider
}

// GetJsonResponseByReceiptImageRecognitionModel returns the json response from the OpenAI common compatible large language model provider
func (p *CommonHttpLargeLanguageModelProvider) GetJsonResponseByReceiptImageRecognitionModel(c core.Context, uid int64, currentConfig *settings.Config, request *LargeLanguageModelRequest) (*LargeLanguageModelTextualResponse, error) {
	return p.getTextualResponse(c, uid, currentConfig, request, p.provider.GetReceiptImageRecognitionModelID(), LARGE_LANGUAGE_MODEL_RESPONSE_FORMAT_JSON)
}

func (p *CommonHttpLargeLanguageModelProvider) getTextualResponse(c core.Context, uid int64, currentConfig *settings.Config, request *LargeLanguageModelRequest, modelId string, responseType LargeLanguageModelResponseFormat) (*LargeLanguageModelTextualResponse, error) {
	if modelId == "" {
		return nil, errs.ErrInvalidLLMModelId
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	utils.SetProxyUrl(transport, currentConfig.LargeLanguageModelAPIProxy)

	if currentConfig.LargeLanguageModelAPISkipTLSVerify {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(currentConfig.LargeLanguageModelAPIRequestTimeout) * time.Millisecond,
	}

	httpRequest, err := p.provider.BuildTextualRequest(c, uid, request, modelId, responseType)

	if err != nil {
		log.Errorf(c, "[http_large_language_model_provider.getTextualResponse] failed to build requests for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	httpRequest.Header.Set("User-Agent", settings.GetUserAgent())

	resp, err := client.Do(httpRequest)

	if err != nil {
		log.Errorf(c, "[http_large_language_model_provider.getTextualResponse] failed to request large language model api for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	log.Debugf(c, "[http_large_language_model_provider.getTextualResponse] response is %s", body)

	if resp.StatusCode != 200 {
		log.Errorf(c, "[http_large_language_model_provider.getTextualResponse] failed to get large language model api response for user \"uid:%d\", because response code is %d", uid, resp.StatusCode)
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	return p.provider.ParseTextualResponse(c, uid, body, responseType)
}

func newCommonHttpLargeLanguageModelProvider(provider HttpLargeLanguageModelProvider) *CommonHttpLargeLanguageModelProvider {
	return &CommonHttpLargeLanguageModelProvider{
		provider: provider,
	}
}
