package cli

import (
	"github.com/elleFlorio/gru/manager"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/codegangsta/cli"
)

func manage(c *cli.Context) {
	etcdAddress := c.String("etcdserver")
	man, err := manager.New(etcdAddress)
	if err != nil {
		log.WithField("err", err).Fatalln("Cannot start manager")
	}
	man.Run()
}
