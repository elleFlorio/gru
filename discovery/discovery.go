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
	discoveries []Discovery

	ErrNotSupported = errors.New("discovery service is not supported")
)

func init() {
	discoveries = []Discovery{
		&EtcdDiscovery{},
	}
}

func New(name string, uri string) (Discovery, error) {
	nodeUUID := node.GetNodeConfig().UUID
	for _, dscvr := range discoveries {
		if dscvr.Name() == name {
			log.WithField("name", name).Debugln("Initializing discovery")
			err := dscvr.Initialize(nodeUUID, uri)
			return dscvr, err
		}
	}

	return nil, ErrNotSupported
}
