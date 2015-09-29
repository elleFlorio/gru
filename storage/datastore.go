package storage

import (
	"errors"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

type Storage interface {
	Name() string
	Initialize() error
	StoreData(string, []byte, string) error
	GetData(string, string) ([]byte, error)
	GetAllData(string) (map[string][]byte, error)
	DeleteData(string, string) error
	DeleteAllData(string) error
}

var (
	dataStores      []Storage
	dataStore       int
	ErrNotSupported = errors.New("Storage system not supported")
)

func init() {
	dataStores = []Storage{
		&internal{},
	}
}

func New(name string) (Storage, error) {
	dataStore = 0
	for index, dtstr := range dataStores {
		if name == dtstr.Name() {
			dataStore = index
			log.WithField("name", dtstr.Name()).Debugln("Initializing storage")
			return dataStores[index], nil
		}
	}

	return dataStores[dataStore], ErrNotSupported
}

func DataStore() Storage {
	return dataStores[dataStore]
}
