package cli

import (
	"fmt"
	"os"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/codegangsta/cli"
)

func Run() {
	app := cli.NewApp()
	app.Name = "Gru"
	app.Usage = "Self-managing container system"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "log-level, l",
			Value: "info",
			Usage: fmt.Sprintf("Log level (options: debug, info, warn, error, fatal, panic)"),
		},
	}

	// logs
	app.Before = func(c *cli.Context) error {
		log.SetOutput(os.Stderr)
		level, err := log.ParseLevel(c.String("log-level"))
		if err != nil {
			log.Fatalf(err.Error())
		}
		log.SetLevel(level)

		// If a log level wasn't specified enforce log-level=error
		if !c.IsSet("log-level") && !c.IsSet("l") {
			log.SetLevel(log.ErrorLevel)
		}

		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:   "start",
			Usage:  "Start the GRU agent",
			Action: start,
		},
		{
			Name:   "create",
			Usage:  "Create a GRU cluster",
			Action: create,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "etcdserver, e",
					Value: "http://localhost:4001",
					Usage: fmt.Sprintf("url of etcd server. Default is 'http://localhost:4001'"),
				},
			},
		},
		{
			Name:   "join",
			Usage:  "join a GRU cluster. Need as argument the ID of the cluster.",
			Action: join,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "etcdserver, e",
					Value: "http://localhost:4001",
					Usage: fmt.Sprintf("url of etcd server. Default is 'http://localhost:4001'"),
				},
				cli.StringFlag{
					Name:  "name, n",
					Value: "random_name",
					Usage: fmt.Sprintf("Name of the node. Default is an awesome random-generated name"),
				},
				cli.StringFlag{
					Name:  "address, a",
					Value: "",
					Usage: fmt.Sprintf("Address of the node. If not provided is taken automatically from the host"),
				},
				cli.StringFlag{
					Name:  "port, p",
					Value: "8080",
					Usage: fmt.Sprintf("Port for the rest api server. Default is 8080"),
				},
			},
		},
	}

	app.Run(os.Args)
}
