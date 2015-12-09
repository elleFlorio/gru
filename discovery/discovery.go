package discovery

import (
	"errors"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

type Discovery interface {
	Name() string
	Initialize(string) error
	Register(string, string) error
	Get(string, Options) (map[string]string, error)
	Set(string, string, Options) error
}

type Options map[string]interface{}

var (
	discoveries []Discovery
	discService int

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

func service() Discovery {
	return discoveries[discService]
}

func Name() string {
	return service().Name()
}

func Initialize(uri string) error {
	return service().Initialize(uri)
}

func Register(nodePath string, nodeAddress string) error {
	return service().Register(nodePath, nodeAddress)
}

func Get(key string, opt Options) (map[string]string, error) {
	return service().Get(key, opt)
}

func Set(key string, value string, opt Options) error {
	return service().Set(key, value, opt)
}
