package log

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/mayswind/lab/pkg/core"
	"github.com/mayswind/lab/pkg/settings"
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

func SetLoggerConfiguration(config *settings.Config) error {
	var bootWriters []io.Writer
	var writers []io.Writer

	bootWriters = append(bootWriters, os.Stdout)

	if config.EnableConsoleLog {
		writers = append(writers, os.Stdout)
	}

	if config.EnableFileLog {
		logFile, err := os.OpenFile(config.FileLogPath, os.O_CREATE|os.O_WRONLY, 0666)

		if err != nil {
			return err
		}

		bootWriters = append(bootWriters, logFile)
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

func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(getFinalLog(format, args...))
}

func DebugfWithRequestId(c *core.Context, format string, args ...interface{}) {
	defaultLogger.WithField(logFieldRequestId, c.GetRequestId()).Debugf(getFinalLog(format, args...))
}

func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(getFinalLog(format, args...))
}

func InfofWithRequestId(c *core.Context, format string, args ...interface{}) {
	defaultLogger.WithField(logFieldRequestId, c.GetRequestId()).Infof(getFinalLog(format, args...))
}

func Warnf(format string, args ...interface{}) {
	defaultLogger.Warnf(getFinalLog(format, args...))
}

func WarnfWithRequestId(c *core.Context, format string, args ...interface{}) {
	defaultLogger.WithField(logFieldRequestId, c.GetRequestId()).Warnf(getFinalLog(format, args...))
}

func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(getFinalLog(format, args...))
}

func ErrorfWithRequestId(c *core.Context, format string, args ...interface{}) {
	defaultLogger.WithField(logFieldRequestId, c.GetRequestId()).Errorf(getFinalLog(format, args...))
}

func ErrorfWithRequestIdAndExtra(c *core.Context, extraString string, format string, args ...interface{}) {
	defaultLogger.WithField(logFieldRequestId, c.GetRequestId()).WithField(logFieldExtra, extraString).Errorf(getFinalLog(format, args...))
}

func BootInfof(format string, args ...interface{}) {
	if bootLogger != nil {
		bootLogger.Infof(getFinalLog(format, args...))
	}
}

func BootWarnf(format string, args ...interface{}) {
	if bootLogger != nil {
		bootLogger.Warnf(getFinalLog(format, args...))
	}
}

func BootErrorf(format string, args ...interface{}) {
	if bootLogger != nil {
		bootLogger.Errorf(getFinalLog(format, args...))
	}
}

func Requestf(c *core.Context, format string, args ...interface{}) {
	if requestLogger != nil {
		requestLogger.WithField(logFieldRequestId, c.GetRequestId()).Infof(getFinalLog(format, args...))
	}
}

func SqlQuery(args ...interface{}) {
	if sqlQueryLogger != nil {
		sqlQueryLogger.Info(args...)
	}
}

func SqlQueryf(format string, args ...interface{}) {
	if sqlQueryLogger != nil {
		sqlQueryLogger.Infof(getFinalLog(format, args...))
	}
}

func getFinalLog(format string, args ...interface{}) string {
	result := fmt.Sprintf(format, args...)
	result = strings.Replace(result, "\n", " ", -1)

	return result
}
