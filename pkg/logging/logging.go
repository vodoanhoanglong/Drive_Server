package logging

import (
	"github.com/sirupsen/logrus"
)

// InitLogger initialize logrus with level and json format
func InitLogger(level string) {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	if level == "" {
		return
	}
	iLevel, err := logrus.ParseLevel(level)

	if err != nil {
		logrus.Fatalf("invalid log level %s", level)
	}
	logrus.SetLevel(iLevel)
}

func LogSuccess(response interface{}, metadata map[string]interface{}) {
	logrus.WithFields(metadata).Info("success")
}

func LogError(err error, metadata map[string]interface{}) {
	logrus.WithError(err).Error(err)
}
