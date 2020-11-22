package executor

import (
	"bufio"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"os"
)

type DBConfig struct {
	DSN string `json:"dsn"`
}

type FileConfig struct {
	SqlFiles []string `json:"files"`
}

type _Config struct {
	DB        DBConfig   `json:"db"`
	SqlConfig FileConfig `json:"sql_config"`
}

var config = _Config{}

func InitConfig(path string) {
	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		logrus.Fatal(err)
		panic(err)
	}

	decoder := json.NewDecoder(bufio.NewReader(file))
	if err := decoder.Decode(&config); err != nil {
		logrus.Fatal(err)
		panic(err)
	}
}

func GetConfig() _Config {
	return config
}
