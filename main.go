package main

import (
	logrus "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
	"sql_disorder_executor/executor"
)

func setLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05"})
	logrus.SetOutput(os.Stderr)
	logrus.SetLevel(logrus.DebugLevel)
}

type SqlDisorderExecutor struct {
	ClientManager  *executor.ClientManager
	SqlFileManager *executor.SqlFileManager
}

func main() {
	setLogger()

	logrus.Info("application started")

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "conf",
				Required: true,
				Usage: "sql disorder executor config file",
			},
		},
		Action: func(c *cli.Context) error {
			name := "sql disorder executor"
			if c.NArg() > 0 {
				name = c.Args().Get(0)
			}

			path := c.String("conf")
			logrus.Infof("execute: %v with config:%v", name, path)

			executor.InitConfig(path)
			logrus.Info(executor.Config.DB)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
