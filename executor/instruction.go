package executor

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
)

type SqlFile struct {
	FileName     string
	Instructions []string
}

func (s *SqlFile) SqlCount() int {
	return len(s.Instructions)
}

func (s *SqlFile) GetInstruction(idx int) string {
	if idx >= len(s.Instructions) {
		err := fmt.Errorf("sql index exceed max sql file instraction count, file:%v, idx:%v, size:%v", s.FileName, idx, len(s.Instructions))
		logrus.Error(err)
		panic(err)
	}

	return s.Instructions[idx]
}

type SqlFileManager interface {
	GetSqlFile(idx int) SqlFile
	SqlFileCount() int
	ListSqlFiles() []SqlFile
}
type SqlFileManagerImpl struct {
	Files []SqlFile
}

func NewSqlFileManager(cfg _Config) SqlFileManager {
	manager := SqlFileManagerImpl{}

	for _, file := range cfg.SqlConfig.SqlFiles {
		_ = func() error {
			ff, err := os.Open(file)
			defer ff.Close()

			if err != nil {
				logrus.Error(err)
				panic(err)
			}

			si := SqlFile{
				FileName: file,
			}

			fi := bufio.NewReader(ff)
			for {
				l, _, err := fi.ReadLine()
				if err == io.EOF {
					break
				}

				ins := strings.TrimSpace(string(l))
				if len(ins) == 0 {
					continue
				}

				si.Instructions = append(si.Instructions, ins)
			}

			manager.Files = append(manager.Files, si)
			return nil
		}()
	}


	return &manager
}

func (s *SqlFileManagerImpl) GetSqlFile(idx int) SqlFile {
	if idx >= s.SqlFileCount() {
		err := fmt.Errorf("sql index exceed max sql file count, idx:%v, count:%v", idx, s.SqlFileCount())
		logrus.Error(err)
		panic(err)
	}

	return s.Files[idx]
}

func (s *SqlFileManagerImpl) SqlFileCount() int {
	return len(s.Files)
}

func (s *SqlFileManagerImpl) ListSqlFiles() []SqlFile {
	return s.Files
}
