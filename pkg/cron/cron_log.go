package cron

import (
	"github.com/go-co-op/gocron/v2"

	"github.com/mayswind/ezbookkeeping/pkg/log"
)

// GocronLoggerAdapter represents the logger adapter for gocron
type GocronLoggerAdapter struct {
}

// Debug logs debug log
func (logger GocronLoggerAdapter) Debug(msg string, args ...any) {
	log.Debugf(msg, args...)
}

// Info logs info log
func (logger GocronLoggerAdapter) Info(msg string, args ...any) {
	log.Infof(msg, args...)
}

// Warn logs warn log
func (logger GocronLoggerAdapter) Warn(msg string, args ...any) {
	log.Warnf(msg, args...)
}

// Error logs error log
func (logger GocronLoggerAdapter) Error(msg string, args ...any) {
	log.Errorf(msg, args...)
}

// NewGocronLoggerAdapter returns a new GocronLoggerAdapter instance
func NewGocronLoggerAdapter() gocron.Logger {
	return GocronLoggerAdapter{}
}
