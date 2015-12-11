package cli

import (
	"fmt"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/codegangsta/cli"

	"github.com/elleFlorio/gru/cluster"
	"github.com/elleFlorio/gru/discovery"
	"github.com/elleFlorio/gru/utils"
)

func create(c *cli.Context) {
	var clusterName string
	if !c.Args().Present() {
		log.Fatalln("Error: missing cluster name")
	} else {
		clusterName = c.Args().First()
	}
	// Get etcd address
	etcdAddress := c.String("etcdserver")

	// Initialize etcd client
	discovery.New("etcd", etcdAddress)

	// Generate cluster
	id, err := utils.GenerateUUID()
	if err != nil {
		log.WithField("err", err).Fatalln("Error generating cluster UUID")
	}
	cluster.RegisterCluster(clusterName, id)
	// Print cluster UUID for join
	fmt.Println(clusterName + ":" + id)
}
