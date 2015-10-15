package storage

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	"github.com/elleFlorio/gru/enum"
)

func TestNew(t *testing.T) {
	supported := "internal"
	notSuported := "notSupported"

	_, err := New(supported)
	assert.NoError(t, err)
	assert.Equal(t, "internal", Name())

	test := client()
	assert.Equal(t, "internal", test.Name())

	_, err = New(notSuported)
	assert.Error(t, err)
	assert.Equal(t, "internal", Name())

	test = client()
	assert.Equal(t, "internal", test.Name())
}

func TestStoreData(t *testing.T) {
	intern := 0
	key := "test"
	data := []byte("pippo")
	var err error

	//INTERNAL
	dataStore = intern
	Initialize()
	err = StoreData(key, data, enum.STATS)
	assert.NoError(t, err)
	err = StoreData(key, data, enum.ANALYTICS)
	assert.NoError(t, err)
	err = StoreData(key, data, enum.PLANS)
	assert.NoError(t, err)
}

func TestGetData(t *testing.T) {
	intern := 0
	key := "test"
	data := []byte("pippo")
	var value []byte

	//INTERNAL
	dataStore = intern
	Initialize()
	StoreData(key, data, enum.STATS)
	StoreData(key, data, enum.ANALYTICS)
	StoreData(key, data, enum.PLANS)
	value, _ = GetData(key, enum.STATS)
	assert.Equal(t, data, value)
	value, _ = GetData(key, enum.ANALYTICS)
	assert.Equal(t, data, value)
	value, _ = GetData(key, enum.PLANS)
	assert.Equal(t, data, value)
}

func TestGetAllData(t *testing.T) {
	intern := 0

	//INTERNAL
	dataStore = intern
	Initialize()
	_, err := GetAllData(enum.STATS)
	assert.NoError(t, err)
	_, err = GetAllData(enum.ANALYTICS)
	assert.NoError(t, err)
	_, err = GetAllData(enum.PLANS)
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
	Initialize()
	StoreData(key, data, enum.STATS)
	StoreData(key, data, enum.ANALYTICS)
	StoreData(key, data, enum.PLANS)

	err = DeleteData(key, enum.STATS)
	assert.NoError(t, err)
	value, _ = GetData(key, enum.STATS)
	assert.Nil(t, value)
	err = DeleteData(key, enum.ANALYTICS)
	assert.NoError(t, err)
	value, _ = GetData(key, enum.ANALYTICS)
	assert.Nil(t, value)
	err = DeleteData(key, enum.PLANS)
	assert.NoError(t, err)
	value, _ = GetData(key, enum.PLANS)
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
	Initialize()
	StoreData(key, data, enum.STATS)
	StoreData(key, data, enum.ANALYTICS)
	StoreData(key, data, enum.PLANS)

	err = DeleteAllData(enum.STATS)
	assert.NoError(t, err)
	value, _ = GetAllData(enum.STATS)
	assert.Empty(t, value)
	err = DeleteAllData(enum.ANALYTICS)
	assert.NoError(t, err)
	value, _ = GetAllData(enum.ANALYTICS)
	assert.Empty(t, value)
	err = DeleteAllData(enum.PLANS)
	assert.NoError(t, err)
	value, _ = GetAllData(enum.PLANS)
	assert.Empty(t, value)
}

func TestGetLocalData(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5}
	New("intenal")

	StoreData(enum.LOCAL.ToString(), data, enum.PLANS)
	test, err := GetLocalData(enum.PLANS)
	assert.NoError(t, err)
	assert.Equal(t, data, test)
}

func TestStoreLocalData(t *testing.T) {
	New("internal")
	var err error
	data := []byte{1, 2, 3, 4, 5}

	err = StoreLocalData(data, enum.PLANS)
	assert.NoError(t, err)
}

func TestGetClusterData(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5}
	New("intenal")

	StoreData(enum.CLUSTER.ToString(), data, enum.PLANS)
	test, err := GetClusterData(enum.PLANS)
	assert.NoError(t, err)
	assert.Equal(t, data, test)
}

func TestStoreClusterData(t *testing.T) {
	New("internal")
	var err error
	data := []byte{1, 2, 3, 4, 5}

	err = StoreClusterData(data, enum.PLANS)
	assert.NoError(t, err)
}
