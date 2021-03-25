package log

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/mayswind/lab/pkg/utils"
)

// LogFormatter represents a log formatter
type LogFormatter struct {
	Prefix       string
	DisableLevel bool
}

// Format writes to log according to the log entry
func (f *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer

	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	b.WriteString(utils.FormatUnixTimeToLongDateTimeInServerTimezone(time.Now().Unix()))
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

	if requestId, exists := entry.Data[logFieldRequestId]; exists {
		b.WriteString(fmt.Sprintf(", r=%s", requestId))
	}

	b.WriteString("\n")

	if extra, exists := entry.Data[logFieldExtra]; exists {
		b.WriteString(extra.(string))
	}

	return b.Bytes(), nil
}
