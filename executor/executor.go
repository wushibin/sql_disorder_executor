package executor

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

type Executor interface {
	Run() error
}

func NewExecutor(fileManager SqlFileManager) Executor {
	return &ExecutorImpl{
		sqlManager: fileManager,
	}
}

type ExecutorImpl struct {
	sqlManager SqlFileManager
}

func (s *ExecutorImpl) Run() error {
	var loopInfoList  []LoopInfo
	for idx, sqlFile := range s.sqlManager.ListSqlFiles() {
		loopInfoList = append(loopInfoList, LoopInfo{TagIndex: idx, Count: sqlFile.SqlCount()})
	}

	if len(loopInfoList) == 0 {
		logrus.Warn("not found the loop info")
		return fmt.Errorf("no sql file to execute")
	}

	generator := NewCombinatorGenerator(loopInfoList)
	defer DestroyCombinatorGenerator(generator)

	for {
		combinator := generator.Generate()
		if combinator.EOF {
			logrus.Info("finished")
			return nil
		}

		logrus.Info("%v", combinator.InstructionFlagList)
	}

	return nil
}
