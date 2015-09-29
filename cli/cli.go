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
		/*if !c.IsSet("log-level") && !c.IsSet("l") {
			log.SetLevel(log.ErrorLevel)
		}*/

		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:   "start",
			Usage:  "Start the GRU agent",
			Action: start,
		},
	}

	app.Run(os.Args)
}
