package datastore

import (
	xorm "xorm.io/xorm/log"

	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// XOrmLoggerAdapter represents the logger adapter for xorm
type XOrmLoggerAdapter struct {
	enable   bool
	logLevel settings.Level
}

// Debug logs debug log
func (logger XOrmLoggerAdapter) Debug(v ...interface{}) {
	log.SqlQuery(v...)
}

// Debugf logs debug log with custom format
func (logger XOrmLoggerAdapter) Debugf(format string, v ...interface{}) {
	log.SqlQueryf(format, v...)
}

// Info logs info log
func (logger XOrmLoggerAdapter) Info(v ...interface{}) {
	log.SqlQuery(v...)
}

// Infof logs info log with custom format
func (logger XOrmLoggerAdapter) Infof(format string, v ...interface{}) {
	log.SqlQueryf(format, v...)
}

// Warn logs warn log
func (logger XOrmLoggerAdapter) Warn(v ...interface{}) {
	log.SqlQuery(v...)
}

// Warnf logs warn log with custom format
func (logger XOrmLoggerAdapter) Warnf(format string, v ...interface{}) {
	log.SqlQueryf(format, v...)
}

// Error logs error log
func (logger XOrmLoggerAdapter) Error(v ...interface{}) {
	log.SqlQuery(v...)
}

// Errorf logs error log with custom format
func (logger XOrmLoggerAdapter) Errorf(format string, v ...interface{}) {
	log.SqlQueryf(format, v...)
}

// Level returns the logger level
func (logger XOrmLoggerAdapter) Level() xorm.LogLevel {
	if logger.logLevel == settings.LOGLEVEL_DEBUG {
		return xorm.LOG_DEBUG
	} else if logger.logLevel == settings.LOGLEVEL_INFO {
		return xorm.LOG_INFO
	} else if logger.logLevel == settings.LOGLEVEL_WARN {
		return xorm.LOG_WARNING
	} else if logger.logLevel == settings.LOGLEVEL_ERROR {
		return xorm.LOG_ERR
	}

	return xorm.LOG_INFO
}

// SetLevel sets the logger level
func (logger XOrmLoggerAdapter) SetLevel(l xorm.LogLevel) {
	if l == xorm.LOG_DEBUG {
		logger.logLevel = settings.LOGLEVEL_DEBUG
	} else if l == xorm.LOG_INFO {
		logger.logLevel = settings.LOGLEVEL_INFO
	} else if l == xorm.LOG_WARNING {
		logger.logLevel = settings.LOGLEVEL_WARN
	} else if l == xorm.LOG_ERR {
		logger.logLevel = settings.LOGLEVEL_ERROR
	}

	logger.logLevel = settings.LOGLEVEL_INFO
}

// ShowSQL sets whether write sql to log
func (logger XOrmLoggerAdapter) ShowSQL(show ...bool) {
	logger.enable = len(show) > 0 && show[0]
}

// IsShowSQL returns whether write sql to log
func (logger XOrmLoggerAdapter) IsShowSQL() bool {
	return logger.enable
}

// NewXOrmLoggerAdapter returns a new XOrmLoggerAdapter instance
func NewXOrmLoggerAdapter(showSql bool, logLevel settings.Level) xorm.Logger {
	return XOrmLoggerAdapter{
		enable:   showSql,
		logLevel: logLevel,
	}
}
