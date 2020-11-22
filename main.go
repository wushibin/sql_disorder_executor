package main

import (
	logrus "github.com/sirupsen/logrus"
	"os"
	"sql_disorder_executor/executor"
)

func setLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05"})
	logrus.SetOutput(os.Stderr)
	logrus.SetLevel(logrus.DebugLevel)
}

type SqlDisorderExecutor struct {
	ClientManager  *executor.ClientManager
	SqlFileManager *executor.SqlFileManager
}

func main() {
	setLogger()

	logrus.Info("hello")
}
