package log

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func Init(path string) error {
	logger = logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		ForceQuote:      true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})

	logger.SetLevel(logrus.InfoLevel)

	output, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {

		return err
	}

	logger.SetOutput(output)

	return nil
}

func Debug(args ...interface{}) {

	logger.Debugln(args...)
}

func Info(args ...interface{}) {

	logger.Infoln(args...)
}

func Error(args ...interface{}) {

	logger.Errorln(args...)
}

func Warn(args ...interface{}) {

	logger.Warnln(args...)
}

func GetWriter() *io.PipeWriter {

	return logger.Writer()
}
