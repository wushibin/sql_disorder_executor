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

// 通过gorm实现数据库连接客户端
type DBClient struct {
	*gorm.DB
}

func (s *DBClient) Execute(sql string) error {
	return s.Exec(sql).Error
}

// Mock的数据库连接客户端，不执行具体SQL，只打印出SQL
type MockClient struct {
	DSN string
}

// 创建连接真实数据库的客户端
func BuildDBClient(dsn string) Client {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Panic(err)
		panic(err)
	}

	return &DBClient{db}
}

// 创建Mock数据库客户端
func BuildMockClient(dsn string) Client  {
	return &MockClient{DSN: dsn}
}

// Mock客户端，打印执行的SQL语句
func (s *MockClient) Execute(sql string) error {
	logrus.Debugf("MockClient: execute sql:%v", sql)
	return nil
}

// 数据库客户端管理
type ClientManager interface {
	GetClient(idx int) Client
	ClientCount() int
}

type ClientManagerImpl struct {
	dbs []Client
}

// 方法会注册到Container， 创建数据库客户端管理的interface。
// 根据SqlFileManager管理的SQL文件数量，创建相应数量的SQL客户端。
// SQL客户端实例与文件实例一一对应
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