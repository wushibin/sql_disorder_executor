package executor

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ClientManager struct {
	dbs []*gorm.DB
}

func NewClientManager() *ClientManager {
	return &ClientManager{}
}

func (s *ClientManager) AddClientInstance(dsn string) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Panic(err)
		panic(err)
	}

	s.dbs = append(s.dbs, db)
}

func (s *ClientManager) GetDB(idx int) *gorm.DB {
	if idx >= len(s.dbs) {
		logrus.Errorf("index is exceed db connections, index:%v, count:%v", idx, len(s.dbs))
		panic(fmt.Errorf("index is excced db connections"))
	}

	return s.dbs[idx]
}