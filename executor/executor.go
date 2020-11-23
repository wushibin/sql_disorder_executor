package executor

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

type Executor interface {
	Run() error
}

func NewExecutor(SqlFileManager SqlFileManager, runner SqlGroupRunner) Executor {
	return &ExecutorImpl{
		SqlFileManager: SqlFileManager,
		SqlGroupRunner: runner,
	}
}

type ExecutorImpl struct {
	SqlFileManager SqlFileManager
	SqlGroupRunner SqlGroupRunner
}

func (s *ExecutorImpl) Run() error {
	var loopInfoList []LoopInfo
	for idx, sqlFile := range s.SqlFileManager.ListSqlFiles() {
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
			s.SqlGroupRunner.Waiting()
			logrus.Info("finished")
			return nil
		}

		s.SqlGroupRunner.RunInstruction(combinator.InstructionFlagList)
	}

	return nil
}
