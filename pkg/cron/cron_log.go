package cron

import (
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/log"
)

// GocronLoggerAdapter represents the logger adapter for gocron
type GocronLoggerAdapter struct {
}

// Debug logs debug log
func (logger GocronLoggerAdapter) Debug(msg string, args ...any) {
	log.Debugf(core.NewNullContext(), logger.getFinalLog(msg, args...))
}

// Info logs info log
func (logger GocronLoggerAdapter) Info(msg string, args ...any) {
	log.Infof(core.NewNullContext(), logger.getFinalLog(msg, args...))
}

// Warn logs warn log
func (logger GocronLoggerAdapter) Warn(msg string, args ...any) {
	log.Warnf(core.NewNullContext(), logger.getFinalLog(msg, args...))
}

// Error logs error log
func (logger GocronLoggerAdapter) Error(msg string, args ...any) {
	log.Errorf(core.NewNullContext(), logger.getFinalLog(msg, args...))
}

func (logger GocronLoggerAdapter) getFinalLog(msg string, args ...any) string {
	var ret strings.Builder
	ret.WriteString(msg)

	for i := 0; i < len(args); i++ {
		if ret.Len() > 0 {
			ret.WriteRune(' ')
		}

		ret.WriteString(fmt.Sprint(args[i]))
	}

	return ret.String()
}

// NewGocronLoggerAdapter returns a new GocronLoggerAdapter instance
func NewGocronLoggerAdapter() gocron.Logger {
	return GocronLoggerAdapter{}
}
