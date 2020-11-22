package executor

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
)

type Executor interface {
	Run() error
}

func NewExecutor(fileManager SqlFileManager) Executor {
	return &ExecutorImpl{
		sqlFileManager: fileManager,
	}
}

type ExecutorImpl struct {
	sqlFileManager SqlFileManager
}

func (s *ExecutorImpl) Run() error {
	var loopInfoList  []LoopInfo
	for idx, sqlFile := range s.sqlFileManager.ListSqlFiles() {
		loopInfoList = append(loopInfoList, LoopInfo{TagIndex: idx, Count: sqlFile.SqlCount()})
	}

	if len(loopInfoList) == 0 {
		logrus.Warn("not found the loop info")
		return fmt.Errorf("no sql file to execute")
	}

	generator := NewCombinatorGenerator(loopInfoList)
	defer DestroyCombinatorGenerator(generator)

	var wg sync.WaitGroup

	for {
		combinator := generator.Generate()
		if combinator.EOF {
			wg.Wait()
			logrus.Info("finished")
			return nil
		}

		logrus.Info(combinator.InstructionFlagList)
	}

	return nil
}
