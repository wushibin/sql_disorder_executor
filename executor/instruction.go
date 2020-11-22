package executor

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"strings"
)

type SqlFile struct {
	FileName     string
	Instructions []string
}

func (s *SqlFile) SqlCount() int {
	return len(s.Instructions)
}

func (s *SqlFile) Instruction(idx int) string {
	if idx >= len(s.Instructions) {
		err := fmt.Errorf("sql index exceed max sql file instraction count, file:%v, idx:%v, size:%v", s.FileName, idx, len(s.Instructions))
		logrus.Error(err)
		panic(err)
	}

	return s.Instructions[idx]
}

type SqlFileManager struct {
	Files []SqlFile
}

func (s *SqlFileManager) AddSqlInstructionFile(name string, file bufio.Reader) {
	si := SqlFile{
		FileName: name,
	}

	for {
		l, _, err := file.ReadLine()
		if err == io.EOF {
			return
		}

		ins := strings.TrimSpace(BytesToString(l))
		if len(ins) == 0 {
			continue
		}

		si.Instructions = append(si.Instructions, BytesToString(l))
	}

	s.Files = append(s.Files, si)
}

func NewSqlFileManager() *SqlFileManager {
	return &SqlFileManager{}
}
