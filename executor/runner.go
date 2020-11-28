package executor

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
)

type SqlGroupRunner interface {
	RunInstruction(task string, instructionList []int) error
	Waiting()
}

func NewSqlGroupRunner(fileManager SqlFileManager, clientManager ClientManager) SqlGroupRunner {
	if fileManager.SqlFileCount() > clientManager.ClientCount() {
		logrus.Error("sql file count is bigger than db client count")
		panic(fmt.Errorf("sql file count is bigger than db client count"))
	}

	return &SqlGroupRunnerImpl{
		SqlFileManager: fileManager,
		ClientManager:  clientManager,
		WaitGroup:      sync.WaitGroup{},
	}
}

type SqlGroupRunnerImpl struct {
	SqlFileManager SqlFileManager
	ClientManager  ClientManager
	WaitGroup      sync.WaitGroup
}

// 根据当前语句的执行序列，调用相应的数据库客户端依次执行SQL
// 如：标识的序列是：（文件1， 文件2， 文件1), SQL的执行顺序是：
// (执行文件1的第1条SQL语句，执行文件2的第1条SQL， 执行文件1的第2条SQL语句)
func (s *SqlGroupRunnerImpl) RunInstruction(taskName string, instructionFlagList []int) error {
	var recordList []*SqlRunner

	// 将SQL文件与相应的数据库客户端对应
	for idx, sqlFile := range s.SqlFileManager.ListSqlFiles() {
		client := s.ClientManager.GetClient(idx)
		runner := SqlRunner{
			Current: 0,
			SqlFile: sqlFile,
			Client:  client,
		}

		recordList = append(recordList, &runner)
	}

	// 创建go routine执行序列对应的SQL
	s.WaitGroup.Add(1)
	go func() {

		for _, flagIndex := range instructionFlagList {
			runner := recordList[flagIndex]
			// 通过SQL文件对应的数据库客户端，执行SQL文件的下一条SQL语名
			if err := runner.ExecNextSqlStatement(taskName); err != nil {
				logrus.Error(err)
				panic(err)
			}
		}

		s.WaitGroup.Done()
	}()

	return nil
}

func (s *SqlGroupRunnerImpl) Waiting() {
	s.WaitGroup.Wait()
}

type SqlRunner struct {
	Current int
	SqlFile SqlFile
	Client  Client
}

func (s *SqlRunner) ExecNextSqlStatement(task string) error {
	// 获取SQL文件的s.Current条SQL语句
	statement := s.SqlFile.GetInstruction(s.Current)
	logrus.Infof("[SqlRunner]: task:%v, sql_file:%v, current:%v, statement:(%v)", task, s.SqlFile.FileName, s.Current, statement)

	err := s.Client.Execute(statement)
	if err != nil {
		return err
	}

	// 标识下一条需要执行的SQL语句
	s.Current++
	return nil
}

