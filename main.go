package main

import (
	logrus "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
	"sql_disorder_executor/di"
	"sql_disorder_executor/executor"
)

func setLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05"})
	logrus.SetOutput(os.Stderr)
	logrus.SetLevel(logrus.DebugLevel)
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

			container := di.GetContainer()
			container.Register(executor.GetConfig)
			container.Register(executor.NewSqlFileManager)
			container.Register(executor.NewClientManager)
			container.Register(executor.NewSqlGroupRunner)
			container.Register(func() executor.ClientBuilder{return executor.BuildMockClient})
			container.Register(executor.NewExecutor)

			err := container.Call(func(executor executor.Executor) error {
				return executor.Run()
			})
			if err != nil {
				logrus.Errorf("execute sql disorder error: %v", err)
				return err
			}

			logrus.Info("execute sql disorder success")
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
