package main

import (
	logrus "github.com/sirupsen/logrus"
	"os"
	"sql_disorder_executor/executor"
)

func setLoger() {
	logrus.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05"})
	logrus.SetOutput(os.Stderr)
	logrus.SetLevel(logrus.DebugLevel)
}

type SqlDisorderExecutor struct {
	DBManager             *executor.DBManager
	SqlInstructionManager *executor.SqlInstructionManager
	//InstractionSortor
}

func main() {
	setLoger()

	logrus.Info("hello")
}
