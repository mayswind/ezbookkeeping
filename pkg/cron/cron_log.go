package cron

import (
	"github.com/go-co-op/gocron/v2"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/log"
)

// GocronLoggerAdapter represents the logger adapter for gocron
type GocronLoggerAdapter struct {
}

// Debug logs debug log
func (logger GocronLoggerAdapter) Debug(msg string, args ...any) {
	log.Debugf(core.NewNullContext(), msg, args...)
}

// Info logs info log
func (logger GocronLoggerAdapter) Info(msg string, args ...any) {
	log.Infof(core.NewNullContext(), msg, args...)
}

// Warn logs warn log
func (logger GocronLoggerAdapter) Warn(msg string, args ...any) {
	log.Warnf(core.NewNullContext(), msg, args...)
}

// Error logs error log
func (logger GocronLoggerAdapter) Error(msg string, args ...any) {
	log.Errorf(core.NewNullContext(), msg, args...)
}

// NewGocronLoggerAdapter returns a new GocronLoggerAdapter instance
func NewGocronLoggerAdapter() gocron.Logger {
	return GocronLoggerAdapter{}
}
