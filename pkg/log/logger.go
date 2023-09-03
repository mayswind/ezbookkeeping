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
	var writers []io.Writer

	if !isDisableBootLog {
		bootWriters = append(bootWriters, os.Stdout)
	}

	if config.EnableConsoleLog {
		writers = append(writers, os.Stdout)
	}

	if config.EnableFileLog {
		logFile, err := os.OpenFile(config.FileLogPath, os.O_CREATE|os.O_WRONLY, 0666)

		if err != nil {
			return err
		}

		if !isDisableBootLog {
			bootWriters = append(bootWriters, logFile)
		}

		writers = append(writers, logFile)
	}

	bootMultipleWriter := io.MultiWriter(bootWriters...)
	multipleWriter := io.MultiWriter(writers...)

	bootLogger.SetOutput(bootMultipleWriter)
	defaultLogger.SetOutput(multipleWriter)
	requestLogger.SetOutput(multipleWriter)
	sqlQueryLogger.SetOutput(multipleWriter)

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

// Debugf logs debug log with custom format
func Debugf(format string, args ...interface{}) {
	defaultLogger.Debug(getFinalLog(format, args...))
}

// DebugfWithRequestId logs debug log with custom format and request id
func DebugfWithRequestId(c *core.Context, format string, args ...interface{}) {
	defaultLogger.WithField(logFieldRequestId, c.GetRequestId()).Debug(getFinalLog(format, args...))
}

// Infof logs info log with custom format
func Infof(format string, args ...interface{}) {
	defaultLogger.Info(getFinalLog(format, args...))
}

// InfofWithRequestId logs info log with custom format and request id
func InfofWithRequestId(c *core.Context, format string, args ...interface{}) {
	defaultLogger.WithField(logFieldRequestId, c.GetRequestId()).Info(getFinalLog(format, args...))
}

// Warnf logs warn log with custom format
func Warnf(format string, args ...interface{}) {
	defaultLogger.Warn(getFinalLog(format, args...))
}

// WarnfWithRequestId logs warn log with custom format and request id
func WarnfWithRequestId(c *core.Context, format string, args ...interface{}) {
	defaultLogger.WithField(logFieldRequestId, c.GetRequestId()).Warn(getFinalLog(format, args...))
}

// Errorf logs error log with custom format
func Errorf(format string, args ...interface{}) {
	defaultLogger.Error(getFinalLog(format, args...))
}

// ErrorfWithRequestId logs error log with custom format and request id
func ErrorfWithRequestId(c *core.Context, format string, args ...interface{}) {
	defaultLogger.WithField(logFieldRequestId, c.GetRequestId()).Error(getFinalLog(format, args...))
}

// ErrorfWithRequestIdAndExtra logs error log with custom format and request id and extra info
func ErrorfWithRequestIdAndExtra(c *core.Context, extraString string, format string, args ...interface{}) {
	defaultLogger.WithField(logFieldRequestId, c.GetRequestId()).WithField(logFieldExtra, extraString).Error(getFinalLog(format, args...))
}

// BootInfof logs boot info log
func BootInfof(format string, args ...interface{}) {
	if bootLogger != nil {
		bootLogger.Info(getFinalLog(format, args...))
	}
}

// BootWarnf logs boot warn log
func BootWarnf(format string, args ...interface{}) {
	if bootLogger != nil {
		bootLogger.Warn(getFinalLog(format, args...))
	}
}

// BootErrorf logs boot error log
func BootErrorf(format string, args ...interface{}) {
	if bootLogger != nil {
		bootLogger.Error(getFinalLog(format, args...))
	}
}

// Requestf logs http request log with custom format
func Requestf(c *core.Context, format string, args ...interface{}) {
	if requestLogger != nil {
		requestLogger.WithField(logFieldRequestId, c.GetRequestId()).Info(getFinalLog(format, args...))
	}
}

// SqlQuery logs sql query log
func SqlQuery(args ...interface{}) {
	if sqlQueryLogger != nil {
		sqlQueryLogger.Info(args...)
	}
}

// SqlQueryf logs sql query log with custom format
func SqlQueryf(format string, args ...interface{}) {
	if sqlQueryLogger != nil {
		sqlQueryLogger.Info(getFinalLog(format, args...))
	}
}

func getFinalLog(format string, args ...interface{}) string {
	result := fmt.Sprintf(format, args...)
	result = strings.Replace(result, "\n", " ", -1)

	return result
}
