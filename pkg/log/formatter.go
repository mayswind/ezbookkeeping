package log

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/mayswind/lab/pkg/utils"
)

type LogFormatter struct {
	Prefix       string
	DisableLevel bool
}

func (f *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer

	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	b.WriteString(utils.FormatToLongDateTime(time.Now()))
	b.WriteString(" ")

	if f.Prefix != "" {
		b.WriteString(f.Prefix)
		b.WriteString(" ")
	}

	if !f.DisableLevel {
		b.WriteString("[")
		b.WriteString(strings.ToUpper(entry.Level.String()))
		b.WriteString("] ")
	}

	b.WriteString(entry.Message)

	if requestId, exists := entry.Data[LOG_FIELD_REQUEST_ID]; exists {
		b.WriteString(fmt.Sprintf(", r=%s", requestId))
	}

	b.WriteString("\n")

	if extra, exists := entry.Data[LOG_FIELD_EXTRA]; exists {
		b.WriteString(extra.(string))
	}

	return b.Bytes(), nil
}
