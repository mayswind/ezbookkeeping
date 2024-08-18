package api

import (
	"bytes"
	"image/png"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

const (
	qrCodeDefaultWidth  int = 320
	qrCodeDefaultHeight int = 320
)

// QrCodesApi represents qrcode generator api
type QrCodesApi struct {
	ApiUsingConfig
}

// Initialize a qrcode generator api singleton instance
var (
	QrCodes = &QrCodesApi{
		ApiUsingConfig: ApiUsingConfig{
			container: settings.Container,
		},
	}
)

// MobileUrlQrCodeHandler returns a mobile url qr code image
func (a *QrCodesApi) MobileUrlQrCodeHandler(c *core.WebContext) ([]byte, string, *errs.Error) {
	fullUrl := a.CurrentConfig().RootUrl + "mobile"
	data, err := a.generateUrlQrCode(c, fullUrl)

	if err != nil {
		return nil, "", errs.ErrOperationFailed
	}

	return data, "image/png", nil
}

func (a *QrCodesApi) generateUrlQrCode(c *core.WebContext, url string) ([]byte, *errs.Error) {
	qrCodeImg, _ := qr.Encode(url, qr.M, qr.Auto)
	qrCodeImg, _ = barcode.Scale(qrCodeImg, qrCodeDefaultWidth, qrCodeDefaultHeight)
	imgData := &bytes.Buffer{}

	if err := png.Encode(imgData, qrCodeImg); err != nil {
		return nil, errs.ErrOperationFailed
	}

	return imgData.Bytes(), nil
}
