package discovery

import (
	"errors"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

type Discovery interface {
	Name() string
	Initialize(string) error
	Register(string, string, int) error
	Get(string) (map[string]string, error)
	Set(string, string, int) error
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
	for index, dscvr := range discoveries {
		if dscvr.Name() == name {
			err := dscvr.Initialize(uri)
			if err != nil {
				return discoveries[discService], err
			}
			discService = index
			log.WithField("name", name).Debugln("Initializing discovery")
			return discoveries[discService], nil
		}
	}

	return discoveries[discService], ErrNotSupported
}

func Service() Discovery {
	return discoveries[discService]
}
