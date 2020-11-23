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

func (s *SqlGroupRunnerImpl) RunInstruction(taskName string, instructionFlagList []int) error {
	var recordList []*SqlRunner

	for idx, sqlFile := range s.SqlFileManager.ListSqlFiles() {
		client := s.ClientManager.GetClient(idx)
		runner := SqlRunner{
			Current: 0,
			SqlFile: sqlFile,
			Client:  client,
		}

		recordList = append(recordList, &runner)
	}

	s.WaitGroup.Add(1)
	go func() {

		for _, instruction := range instructionFlagList {
			runner := recordList[instruction]
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
	statement := s.SqlFile.GetInstruction(s.Current)
	logrus.Infof("[SqlRunner]: task:%v, sql_file:%v, current:%v, statement:(%v)", task, s.SqlFile.FileName, s.Current, statement)

	err := s.Client.Execute(statement)
	if err != nil {
		return err
	}

	s.Current++
	return nil
}

