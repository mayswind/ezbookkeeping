package extapi

import (
	"net/http"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	extservices "github.com/mayswind/ezbookkeeping/pkg/ext/services"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// lazyZipWriter starts the download (headers and status) only when the first bytes arrive, so a failure before that can
// still be answered with a normal error instead of a broken file.
type lazyZipWriter struct {
	c       *core.WebContext
	started bool
}

func (w *lazyZipWriter) Write(data []byte) (int, error) {
	if !w.started {
		w.started = true
		w.c.Header("Content-Type", "application/zip")
		w.c.Header("Content-Disposition", `attachment; filename="business-data-`+time.Now().UTC().Format("20060102")+`.zip"`)
		w.c.Header("Cache-Control", "no-store")
		w.c.Header("X-Content-Type-Options", "nosniff")
		w.c.Status(http.StatusOK)
	}

	return w.c.Writer.Write(data)
}

// ExportBusinessHandler streams the caller's own business records as a ZIP of CSV files. It is always about the caller's
// own business and only they can ask for it: the route is on the identity list, so a business header cannot redirect it.
func (h *Handlers) ExportBusinessHandler(c *core.WebContext) {
	uid := c.GetActualUid()
	out := &lazyZipWriter{c: c}

	if err := extservices.Exports.WriteZip(c, uid, out); err != nil {
		log.Errorf(c, "[ext.export] failed to export the business of user \"uid:%d\", because %s", uid, err.Error())

		if !out.started {
			utils.PrintJsonErrorResult(c, errs.Or(err, errs.ErrOperationFailed))
		}

		return
	}

	log.Infof(c, "[ext.export] user \"uid:%d\" exported the business records", uid)
}
