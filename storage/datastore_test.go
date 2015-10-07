package storage

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	"github.com/elleFlorio/gru/enum"
)

func TestNew(t *testing.T) {
	supported := "internal"
	notSuported := "notSupported"

	test, err := New(supported)
	assert.NoError(t, err, "Supported storage should produce no error")
	assert.Equal(t, "internal", test.Name(), "Storage should be 'internal'")

	test = DataStore()
	assert.Equal(t, "internal", test.Name(), "(supported) Retrieved datastore should be 'internal'")

	test, err = New(notSuported)
	assert.Error(t, err, "Not supported storage should produce an error")
	assert.Equal(t, "internal", test.Name(), "If storage is not supported the default one should be 'internal'")

	test = DataStore()
	assert.Equal(t, "internal", test.Name(), "(not supported) retrieved datastore should be 'internal'")
}

func TestStoreData(t *testing.T) {
	intern := 0
	key := "test"
	data := []byte("pippo")
	var err error

	//INTERNAL
	dataStore = intern
	DataStore().Initialize()
	err = DataStore().StoreData(key, data, enum.STATS)
	assert.NoError(t, err)
	err = DataStore().StoreData(key, data, enum.ANALYTICS)
	assert.NoError(t, err)
}

func TestGetData(t *testing.T) {
	intern := 0
	key := "test"
	data := []byte("pippo")
	var value []byte

	//INTERNAL
	dataStore = intern
	DataStore().Initialize()
	DataStore().StoreData(key, data, enum.STATS)
	DataStore().StoreData(key, data, enum.ANALYTICS)
	value, _ = DataStore().GetData(key, enum.STATS)
	assert.Equal(t, data, value)
	value, _ = DataStore().GetData(key, enum.ANALYTICS)
	assert.Equal(t, data, value)
}

func TestGetAllData(t *testing.T) {
	intern := 0

	//INTERNAL
	dataStore = intern
	DataStore().Initialize()
	_, err := DataStore().GetAllData(enum.STATS)
	assert.NoError(t, err)
	_, err = DataStore().GetAllData(enum.ANALYTICS)
	assert.NoError(t, err)
}

func TestDeleteData(t *testing.T) {
	intern := 0
	key := "test"
	data := []byte("pippo")
	var value []byte
	var err error

	//INTERNAL
	dataStore = intern
	DataStore().Initialize()
	DataStore().StoreData(key, data, enum.STATS)
	DataStore().StoreData(key, data, enum.ANALYTICS)

	err = DataStore().DeleteData(key, enum.STATS)
	assert.NoError(t, err)
	value, _ = DataStore().GetData(key, enum.STATS)
	assert.Nil(t, value)
	err = DataStore().DeleteData(key, enum.ANALYTICS)
	assert.NoError(t, err)
	value, _ = DataStore().GetData(key, enum.ANALYTICS)
	assert.Nil(t, value)
}

func TestDeleteAllData(t *testing.T) {
	intern := 0
	key := "test"
	data := []byte("pippo")
	var value map[string][]byte
	var err error

	//INTERNAL
	dataStore = intern
	DataStore().Initialize()
	DataStore().StoreData(key, data, enum.STATS)
	DataStore().StoreData(key, data, enum.ANALYTICS)

	err = DataStore().DeleteAllData(enum.STATS)
	assert.NoError(t, err)
	value, _ = DataStore().GetAllData(enum.STATS)
	assert.Empty(t, value)
	err = DataStore().DeleteAllData(enum.ANALYTICS)
	assert.NoError(t, err)
	value, _ = DataStore().GetAllData(enum.ANALYTICS)
	assert.Empty(t, value)
}
