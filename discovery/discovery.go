package discovery

import (
	"errors"

	log "github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/node"
)

type Discovery interface {
	Name() string
	Initialize(string, string) error
	Register(string, uint64) error
	Get(string) (map[string]string, error)
	Set(string, string, uint64) error
}

var (
	discoveries     []Discovery
	discService     int
	ErrNotSupported = errors.New("discovery service not supported")
)

func init() {
	discoveries = []Discovery{
		&noService{},
		&etcdDiscovery{},
	}
}

func New(name string, uri string) (Discovery, error) {
	discService = 0
	nodeUUID := node.Config().UUID
	for index, dscvr := range discoveries {
		if dscvr.Name() == name {
			err := dscvr.Initialize(nodeUUID, uri)
			discService = index
			log.WithField("name", name).Debugln("Initializing discovery")
			return discoveries[discService], err
		}
	}

	return discoveries[discService], ErrNotSupported
}

func Service() Discovery {
	return discoveries[discService]
}
