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
	// 根据第一个SQL文件的数量及第个SQL文件中语句的数量，创建枚举所有执行序列时需要的信息
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

	idx := 0
	for {
		// 获取当前标识的语句执行序列
		combinator := generator.Generate()
		if combinator.EOF {
			s.SqlGroupRunner.Waiting()
			logrus.Infof("finished, loop count:%v", idx)
			return nil
		}

		// 根据当前语句的执行序列，调用相应的数据库客户端依次执行SQL
		// 如：标识的序列是：（文件1， 文件2， 文件1), SQL的执行顺序是：
		// 执行文件1的第1条SQL语句，执行文件2的第1条SQL， 执行文件1的第2条SQL语句
		s.SqlGroupRunner.RunInstruction(fmt.Sprintf("loop_%v", idx), combinator.InstructionFlagList)
		idx++
	}

	return nil
}
