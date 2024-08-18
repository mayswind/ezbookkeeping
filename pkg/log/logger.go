package log

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

const logFieldRequestId = "REQUEST_ID"
const logFieldExtra = "EXTRA"

var bootLogger = logrus.New()
var defaultLogger = logrus.New()
var requestLogger = logrus.New()
var sqlQueryLogger = logrus.New()

func init() {
	bootLogger.SetFormatter(&LogFormatter{})
	bootLogger.SetOutput(os.Stdout)
	bootLogger.SetLevel(logrus.InfoLevel)

	defaultLogger.SetFormatter(&LogFormatter{})
	defaultLogger.SetOutput(os.Stdout)
	defaultLogger.SetLevel(logrus.InfoLevel)

	requestLogger.SetFormatter(&LogFormatter{Prefix: "[REQUEST]", DisableLevel: true})
	requestLogger.SetOutput(os.Stdout)
	requestLogger.SetLevel(logrus.InfoLevel)

	sqlQueryLogger.SetFormatter(&LogFormatter{Prefix: "[SQLQUERY]", DisableLevel: true})
	sqlQueryLogger.SetOutput(os.Stdout)
	sqlQueryLogger.SetLevel(logrus.InfoLevel)
}

// SetLoggerConfiguration sets the logger according to the config
func SetLoggerConfiguration(config *settings.Config, isDisableBootLog bool) error {
	var bootWriters []io.Writer
	var defaultWriters []io.Writer
	var requestWriters []io.Writer
	var queryWriters []io.Writer

	if !isDisableBootLog {
		bootWriters = append(bootWriters, os.Stdout)
	}

	if config.EnableConsoleLog {
		defaultWriters = append(defaultWriters, os.Stdout)
		requestWriters = append(requestWriters, os.Stdout)
		queryWriters = append(queryWriters, os.Stdout)
	}

	if config.EnableFileLog {
		defaultWriter, err := NewRotateFileWriter(config.FileLogPath, config.LogFileRotate, int64(config.LogFileMaxSize), config.LogFileMaxDays)

		if err != nil {
			return err
		}

		if !isDisableBootLog {
			bootWriters = append(bootWriters, defaultWriter)
		}

		defaultWriters = append(defaultWriters, defaultWriter)

		if config.EnableRequestLog {
			if config.RequestFileLogPath != "" && config.RequestFileLogPath != config.FileLogPath {
				requestWriter, err := NewRotateFileWriter(config.RequestFileLogPath, config.LogFileRotate, int64(config.LogFileMaxSize), config.LogFileMaxDays)

				if err != nil {
					return err
				}

				requestWriters = append(requestWriters, requestWriter)
			} else {
				requestWriters = append(requestWriters, defaultWriter)
			}
		}

		if config.EnableQueryLog {
			if config.QueryFileLogPath != "" && config.QueryFileLogPath != config.FileLogPath {
				queryWriter, err := NewRotateFileWriter(config.QueryFileLogPath, config.LogFileRotate, int64(config.LogFileMaxSize), config.LogFileMaxDays)

				if err != nil {
					return err
				}

				queryWriters = append(queryWriters, queryWriter)
			} else {
				queryWriters = append(queryWriters, defaultWriter)
			}
		}
	}

	bootMultipleWriter := io.MultiWriter(bootWriters...)
	defaultMultipleWriter := io.MultiWriter(defaultWriters...)
	requestMultipleWriter := io.MultiWriter(requestWriters...)
	queryMultipleWriter := io.MultiWriter(queryWriters...)

	bootLogger.SetOutput(bootMultipleWriter)
	defaultLogger.SetOutput(defaultMultipleWriter)
	requestLogger.SetOutput(requestMultipleWriter)
	sqlQueryLogger.SetOutput(queryMultipleWriter)

	if config.LogLevel == settings.LOGLEVEL_DEBUG {
		bootLogger.SetLevel(logrus.DebugLevel)
		defaultLogger.SetLevel(logrus.DebugLevel)
	} else if config.LogLevel == settings.LOGLEVEL_INFO {
		bootLogger.SetLevel(logrus.InfoLevel)
		defaultLogger.SetLevel(logrus.InfoLevel)
	} else if config.LogLevel == settings.LOGLEVEL_WARN {
		bootLogger.SetLevel(logrus.WarnLevel)
		defaultLogger.SetLevel(logrus.WarnLevel)
	} else if config.LogLevel == settings.LOGLEVEL_ERROR {
		bootLogger.SetLevel(logrus.ErrorLevel)
		defaultLogger.SetLevel(logrus.ErrorLevel)
	}

	if !config.EnableRequestLog {
		requestLogger = nil
	}

	if !config.EnableQueryLog {
		sqlQueryLogger = nil
	}

	return nil
}

// DebugfWithRequestId logs debug log with custom format
func Debugf(c core.Context, format string, args ...any) {
	if c == nil {
		defaultLogger.Debug(getFinalLog(format, args...))
	} else {
		defaultLogger.WithField(logFieldRequestId, c.GetContextId()).Debug(getFinalLog(format, args...))
	}
}

// Infof logs info log with custom format
func Infof(c core.Context, format string, args ...any) {
	if c == nil {
		defaultLogger.Info(getFinalLog(format, args...))
	} else {
		defaultLogger.WithField(logFieldRequestId, c.GetContextId()).Info(getFinalLog(format, args...))
	}
}

// Warnf logs warn log with custom format
func Warnf(c core.Context, format string, args ...any) {
	if c == nil {
		defaultLogger.Warn(getFinalLog(format, args...))
	} else {
		defaultLogger.WithField(logFieldRequestId, c.GetContextId()).Warn(getFinalLog(format, args...))
	}
}

// Errorf logs error log with custom format
func Errorf(c core.Context, format string, args ...any) {
	if c == nil {
		defaultLogger.Error(getFinalLog(format, args...))
	} else {
		defaultLogger.WithField(logFieldRequestId, c.GetContextId()).Error(getFinalLog(format, args...))
	}
}

// ErrorfWithExtra logs error log with custom format and extra info
func ErrorfWithExtra(c core.Context, extraString string, format string, args ...any) {
	if c == nil {
		defaultLogger.WithField(logFieldExtra, extraString).Error(getFinalLog(format, args...))
	} else {
		defaultLogger.WithField(logFieldRequestId, c.GetContextId()).WithField(logFieldExtra, extraString).Error(getFinalLog(format, args...))
	}
}

// BootInfof logs boot info log
func BootInfof(c core.Context, format string, args ...any) {
	if bootLogger != nil {
		if c == nil {
			bootLogger.Info(getFinalLog(format, args...))
		} else {
			bootLogger.WithField(logFieldRequestId, c.GetContextId()).Info(getFinalLog(format, args...))
		}
	}
}

// BootWarnf logs boot warn log
func BootWarnf(c core.Context, format string, args ...any) {
	if bootLogger != nil {
		if c == nil {
			bootLogger.Warn(getFinalLog(format, args...))
		} else {
			bootLogger.WithField(logFieldRequestId, c.GetContextId()).Warn(getFinalLog(format, args...))
		}
	}
}

// BootErrorf logs boot error log
func BootErrorf(c core.Context, format string, args ...any) {
	if bootLogger != nil {
		if c == nil {
			bootLogger.Error(getFinalLog(format, args...))
		} else {
			bootLogger.WithField(logFieldRequestId, c.GetContextId()).Error(getFinalLog(format, args...))
		}
	}}

// Requestf logs http request log with custom format
func Requestf(c core.Context, format string, args ...any) {
	if requestLogger != nil {
		requestLogger.WithField(logFieldRequestId, c.GetContextId()).Info(getFinalLog(format, args...))
	}
}

// SqlQuery logs sql query log
func SqlQuery(args ...any) {
	if sqlQueryLogger != nil {
		sqlQueryLogger.Info(args...)
	}
}

// SqlQueryf logs sql query log with custom format
func SqlQueryf(format string, args ...any) {
	if sqlQueryLogger != nil {
		sqlQueryLogger.Info(getFinalLog(format, args...))
	}
}

func getFinalLog(format string, args ...any) string {
	result := fmt.Sprintf(format, args...)
	result = strings.Replace(result, "\n", " ", -1)

	return result
}
