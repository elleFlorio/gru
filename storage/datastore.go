package storage

import (
	"errors"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/enum"
)

//Local data key
const local string = "local"

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
			log.WithField("name", dtstr.Name()).Debugln("Initializing storage")
			return dataStores[index], nil
		}
	}

	return dataStores[dataStore], ErrNotSupported
}

func DataStore() Storage {
	return dataStores[dataStore]
}

func GetLocalData(dataType enum.Datatype) ([]byte, error) {
	return DataStore().GetData(local, dataType)
}

func StoreLocalData(data []byte, dataType enum.Datatype) error {
	return DataStore().StoreData(local, data, dataType)
}
