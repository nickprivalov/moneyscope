package common

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func InitLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:   false,
		ForceColors:     true,
		TimestampFormat: time.DateTime, // "2006-01-02 15:04:05"
	})

	logrus.SetOutput(os.Stdout)

	logrus.SetLevel(logrus.InfoLevel)

	logrus.WithFields(logrus.Fields{
		"package":  "common",
		"function": "InitLogger",
	}).Printf("Logrus successfully initialized.")
}
