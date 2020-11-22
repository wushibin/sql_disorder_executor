package executor

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBManager struct {
	dbs []*gorm.DB
}

func NewDBManager() *DBManager {
	return &DBManager{}
}

func (s *DBManager) AddDBInstance(dsn string) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Panic(err)
		panic(err)
	}

	s.dbs = append(s.dbs, db)
}

func (s *DBManager) GetDB(idx int) *gorm.DB {
	if idx >= len(s.dbs) {
		logrus.Errorf("index is exceed db connections, index:%v, count:%v", idx, len(s.dbs))
		panic(fmt.Errorf("index is excced db connections"))
	}

	return s.dbs[idx]
}