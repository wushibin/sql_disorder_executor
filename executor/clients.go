package executor

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ClientBuilder func(dsn string) Client

type Client interface {
	Execute(sql string) error
}

type DBClient struct {
	*gorm.DB
}

func (s *DBClient) Execute(sql string) error {
	return s.Exec(sql).Error
}

type MockClient struct {
	DSN string
}

func BuildDBClient(dsn string) Client {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Panic(err)
		panic(err)
	}

	return &DBClient{db}
}

func BuildMockClient(dsn string) Client  {
	return &MockClient{DSN: dsn}
}

func (s *MockClient) Execute(sql string) error {
	logrus.Debugf("MockClient: DSN:%v, Sql:%v", s.DSN, sql)
	return nil
}

type ClientManager interface {
	GetClient(idx int) Client
	ClientCount() int
}

type ClientManagerImpl struct {
	dbs []Client
}

func NewClientManager(cfg _Config, fileManager SqlFileManager, clientBuilder ClientBuilder) ClientManager {
	m := &ClientManagerImpl{}

	for i:=0; i<fileManager.SqlFileCount(); i++ {
		m.dbs = append(m.dbs, clientBuilder(cfg.DB.DSN))
	}

	return m
}

func (s *ClientManagerImpl) GetClient(idx int) Client {
	if idx >= len(s.dbs) {
		logrus.Errorf("index is exceed db connections, index:%v, count:%v", idx, len(s.dbs))
		panic(fmt.Errorf("index is excced db connections"))
	}

	return s.dbs[idx]
}

func (s *ClientManagerImpl) ClientCount() int {
	return len(s.dbs)
}