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
}

// Initialize a qrcode generator api singleton instance
var (
	QrCodes = &QrCodesApi{}
)

// MobileUrlQrCodeHandler returns a mobile url qr code image
func (a *QrCodesApi) MobileUrlQrCodeHandler(c *core.Context) ([]byte, string, *errs.Error) {
	fullUrl := settings.Container.Current.RootUrl + "mobile"
	data, err := a.generateUrlQrCode(c, fullUrl)

	if err != nil {
		return nil, "", errs.ErrOperationFailed
	}

	return data, "", nil
}

func (a *QrCodesApi) generateUrlQrCode(c *core.Context, url string) ([]byte, *errs.Error) {
	qrCodeImg, _ := qr.Encode(url, qr.M, qr.Auto)
	qrCodeImg, _ = barcode.Scale(qrCodeImg, qrCodeDefaultWidth, qrCodeDefaultHeight)
	imgData := &bytes.Buffer{}

	if err := png.Encode(imgData, qrCodeImg); err != nil {
		return nil, errs.ErrOperationFailed
	}

	return imgData.Bytes(), nil
}
