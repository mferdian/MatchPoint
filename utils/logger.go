package utils

import (
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func SetUpLogger() {
	err := os.MkdirAll("temp", os.ModePerm)
	if err != nil {
		panic("failed to create log directory: " + err.Error())
	}

	logFile := filepath.Join("temp", "app.log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("failed to open log file: " + err.Error())
	}

	// Output ke console dan file sekaligus
	multiWriter := io.MultiWriter(os.Stdout, file)
	Log.SetOutput(multiWriter)

	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	Log.SetLevel(logrus.InfoLevel)
}