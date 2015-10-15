package storage

import (
	"errors"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/enum"
)

type Storage interface {
	Name() string
	Initialize() error
	StoreData(string, []byte, enum.Datatype) error
	GetData(string, enum.Datatype) ([]byte, error)
	GetAllData(enum.Datatype) (map[string][]byte, error)
	DeleteData(string, enum.Datatype) error
	DeleteAllData(enum.Datatype) error
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
			log.WithField("name", name).Debugln("Initializing datastore")
			err := dataStores[dataStore].Initialize()
			return dataStores[index], err
		}
	}

	return dataStores[dataStore], ErrNotSupported
}

func client() Storage {
	return dataStores[dataStore]
}

func Initialize() error {
	return client().Initialize()
}

func Name() string {
	return client().Name()
}

func StoreData(key string, data []byte, dataType enum.Datatype) error {
	return client().StoreData(key, data, dataType)
}

func GetData(key string, dataType enum.Datatype) ([]byte, error) {
	return client().GetData(key, dataType)
}

func GetAllData(dataType enum.Datatype) (map[string][]byte, error) {
	return client().GetAllData(dataType)
}

func DeleteData(key string, dataType enum.Datatype) error {
	return client().DeleteData(key, dataType)
}

func DeleteAllData(dataType enum.Datatype) error {
	return client().DeleteAllData(dataType)
}

func GetLocalData(dataType enum.Datatype) ([]byte, error) {
	return client().GetData(enum.LOCAL.ToString(), dataType)
}

func StoreLocalData(data []byte, dataType enum.Datatype) error {
	return client().StoreData(enum.LOCAL.ToString(), data, dataType)
}

func GetClusterData(dataType enum.Datatype) ([]byte, error) {
	return client().GetData(enum.CLUSTER.ToString(), dataType)
}

func StoreClusterData(data []byte, dataType enum.Datatype) error {
	return client().StoreData(enum.CLUSTER.ToString(), data, dataType)
}
