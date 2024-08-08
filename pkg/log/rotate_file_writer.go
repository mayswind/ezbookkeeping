package log

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

const (
	logRotateSuffixDateFormat = "20060102150405"
)

type RotateFileWriter struct {
	EnableRotate bool
	MaxFileSize  int64
	MaxFileDays  uint32

	filePath  string
	file      *os.File
	totalSize int64

	mutex                 sync.Mutex
	lastRemoveOldFilesDay int
}

var logFallbackLogger = logrus.New()

func init() {
	logFallbackLogger.SetFormatter(&LogFormatter{})
	logFallbackLogger.SetOutput(os.Stdout)
	logFallbackLogger.SetLevel(logrus.InfoLevel)
}

// NewRotateFileWriter returns a new rotate file writer
func NewRotateFileWriter(filePath string, enableRotate bool, maxFileSize int64, maxFileDays uint32) (*RotateFileWriter, error) {
	writer := &RotateFileWriter{
		EnableRotate: enableRotate,
		MaxFileSize:  maxFileSize,
		MaxFileDays:  maxFileDays,
		filePath:     filePath,
		totalSize:    0,
	}

	err := writer.openFile()

	if err != nil {
		return nil, err
	}

	return writer, nil
}

// Write does log data to specified file
func (w *RotateFileWriter) Write(p []byte) (n int, err error) {
	dataSize := int64(len(p))

	if w.EnableRotate && w.totalSize > 0 && w.totalSize+dataSize >= w.MaxFileSize {
		w.mutex.Lock()
		defer w.mutex.Unlock()

		if w.EnableRotate && w.totalSize > 0 && w.totalSize+dataSize >= w.MaxFileSize {
			err := w.rotateFile()

			if err != nil {
				logFallbackLogger.Errorf("[rotate_file_writer.Write] cannot rotate log file \"%s\", because %s", w.file.Name(), err.Error())
				return 0, err
			}
		}
	}

	writeSize, err := w.file.Write(p)

	if err != nil {
		return 0, err
	}

	w.totalSize += int64(writeSize)

	if w.EnableRotate {
		today := time.Now().Day()

		if today != w.lastRemoveOldFilesDay && w.MaxFileDays > 0 {
			w.lastRemoveOldFilesDay = today
			go w.removeOldFiles()
		}
	}

	return writeSize, err
}

func (w *RotateFileWriter) rotateFile() error {
	currentFileName := w.file.Name()
	err := w.file.Close()

	if err != nil {
		return errs.NewLoggingError(fmt.Sprintf("cannot close log file \"%s\", because %s", w.file.Name(), err.Error()), err)
	}

	w.file = nil
	archiveFileName := fmt.Sprintf("%s.%s", currentFileName, time.Now().Format(logRotateSuffixDateFormat))
	err = os.Rename(currentFileName, archiveFileName)

	if err != nil {
		return errs.NewLoggingError(fmt.Sprintf("cannot rename log file \"%s\" to \"%s\", because %s", currentFileName, archiveFileName, err.Error()), err)
	}

	err = w.openFile()

	if err != nil {
		return err
	}

	return nil
}

func (w *RotateFileWriter) openFile() error {
	if w.file != nil {
		logFallbackLogger.Warnf("[rotate_file_writer.removeOldFiles] cannot reopen log file \"%s\"", w.file.Name())
		return nil
	}

	file, err := os.OpenFile(w.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	if err != nil {
		return errs.NewLoggingError(fmt.Sprintf("cannot open log file \"%s\", because %s", w.filePath, err.Error()), err)
	}

	w.file = file
	return nil
}

func (w *RotateFileWriter) removeOldFiles() {
	dir := filepath.Dir(w.filePath)
	logBaseFileName := filepath.Base(w.filePath) + "."

	allLogFiles, err := os.ReadDir(dir)

	if err != nil {
		return
	}

	retainMinUnixTime := int64(0)

	if w.MaxFileDays > 0 {
		retainMinUnixTime = time.Now().AddDate(0, 0, -int(w.MaxFileDays)).Unix()
	}

	for _, file := range allLogFiles {
		if file.IsDir() {
			continue
		}

		logFileName := filepath.Base(file.Name())

		if !strings.HasPrefix(logFileName, logBaseFileName) {
			continue
		}

		rotateDate := logFileName[len(logBaseFileName):]
		dotIndex := strings.Index(rotateDate, ".")

		if dotIndex > 0 {
			rotateDate = rotateDate[0:dotIndex]
		}

		if len(rotateDate) != len(logRotateSuffixDateFormat) {
			logFallbackLogger.Errorf("[rotate_file_writer.removeOldFiles] date suffix of old log file \"%s\" is invalid", file.Name())
			continue
		}

		rotateDateTime, err := time.ParseInLocation(logRotateSuffixDateFormat, rotateDate, time.Now().Location())

		if err != nil {
			logFallbackLogger.Errorf("[rotate_file_writer.removeOldFiles] cannot parse rotate date of old log file \"%s\", because %s", file.Name(), err.Error())
			continue
		}

		if rotateDateTime.Unix() >= retainMinUnixTime {
			continue
		}

		err = os.Remove(filepath.Join(dir, file.Name()))

		if err != nil {
			logFallbackLogger.Errorf("[rotate_file_writer.removeOldFiles] cannot remove old log file \"%s\", because %s", file.Name(), err.Error())
		}
	}
}
