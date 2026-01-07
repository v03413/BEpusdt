package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var (
	Task     *logrus.Logger
	bepusdt  *logrus.Logger
	logFiles []*os.File
)

func newLogger(file string) (*logrus.Logger, error) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		ForceQuote:      true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})

	logger.SetLevel(logrus.InfoLevel)

	output, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	// 保存文件句柄
	logFiles = append(logFiles, output)
	logger.SetOutput(output)

	return logger, nil
}

func Init(dir string) error {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {

		return fmt.Errorf("创建日志目录失败：%w", err)
	}

	var err error

	bepusdt, err = newLogger(filepath.Join(dir, "bepusdt.log"))
	if err != nil {
		return err
	}

	Task, err = newLogger(filepath.Join(dir, "task.log"))
	if err != nil {
		return err
	}

	return nil
}

func Debug(args ...interface{}) {
	bepusdt.Debugln(args...)
}

func Info(args ...interface{}) {
	bepusdt.Infoln(args...)
}

func Error(args ...interface{}) {
	bepusdt.Errorln(args...)
}

func Warn(args ...interface{}) {
	bepusdt.Warnln(args...)
}

func GetWriter() *io.PipeWriter {

	return bepusdt.Writer()
}

func Close() {
	var lastErr error
	for _, f := range logFiles {
		if f != nil {
			if err := f.Close(); err != nil {
				lastErr = err
			}
		}
	}

	logFiles = nil

	if lastErr != nil {
		_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf("日志句柄资源关闭错误：%s", lastErr.Error()))
	}
}
