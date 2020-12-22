package datastore

import (
	xorm "xorm.io/xorm/log"

	"github.com/mayswind/lab/pkg/log"
	"github.com/mayswind/lab/pkg/settings"
)

type XOrmLoggerAdapter struct {
	enable   bool
	logLevel settings.Level
}

func (logger XOrmLoggerAdapter) Debug(v ...interface{}) {
	log.SqlQuery(v...)
}

func (logger XOrmLoggerAdapter) Debugf(format string, v ...interface{}) {
	log.SqlQueryf(format, v...)
}

func (logger XOrmLoggerAdapter) Info(v ...interface{}) {
	log.SqlQuery(v...)
}

func (logger XOrmLoggerAdapter) Infof(format string, v ...interface{}) {
	log.SqlQueryf(format, v...)
}

func (logger XOrmLoggerAdapter) Warn(v ...interface{}) {
	log.SqlQuery(v...)
}

func (logger XOrmLoggerAdapter) Warnf(format string, v ...interface{}) {
	log.SqlQueryf(format, v...)
}

func (logger XOrmLoggerAdapter) Error(v ...interface{}) {
	log.SqlQuery(v...)
}

func (logger XOrmLoggerAdapter) Errorf(format string, v ...interface{}) {
	log.SqlQueryf(format, v...)
}

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

func (logger XOrmLoggerAdapter) ShowSQL(show ...bool) {
	logger.enable = len(show) > 0 && show[0]
}

func (logger XOrmLoggerAdapter) IsShowSQL() bool {
	return logger.enable
}

func NewXOrmLoggerAdapter(showSql bool, logLevel settings.Level) xorm.Logger {
	return XOrmLoggerAdapter{
		enable:   showSql,
		logLevel: logLevel,
	}
}
