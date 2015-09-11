package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStoreData(t *testing.T) {
	test := createInternalStorage()
	statsType := "stats"
	notSupported := "notSupported"
	mockData := []byte("pippo")

	err := test.StoreData("test", mockData, statsType)
	assert.NoError(t, err, "(stats type) Store data should produce no error")

	err = test.StoreData("test", mockData, notSupported)
	assert.Error(t, err, "(not supported type) Store data should produce an error")
}

func TestGetData(t *testing.T) {
	test := createInternalStorage()
	statsType := "stats"
	notSupported := "notSupported"
	mockKey := "default"

	data, err := test.GetData(mockKey, statsType)
	assert.NoError(t, err, "(stats type) Get data should produce no error")
	assert.Equal(t, data, test.statsData[mockKey], "(stats type) Retrieved data should be equal to stored one")

	data, err = test.GetData(mockKey, notSupported)
	assert.Error(t, err, "(not supported type) Get data should produce an error")

}

func TestGetAllData(t *testing.T) {
	test := createInternalStorage()
	statsType := "stats"
	notSupported := "notSupported"

	_, err := test.GetAllData(statsType)
	assert.NoError(t, err, "(stats type) Get data should produce no error")

	_, err = test.GetAllData(notSupported)
	assert.Error(t, err, "(not supported type) Get data should produce an error")
}

func TestDeleteData(t *testing.T) {
	test := createInternalStorage()
	statsType := "stats"
	notSupported := "notSupported"
	mockKey := "default"

	err := test.DeleteData(mockKey, statsType)
	assert.NoError(t, err, "(stats type) Delete data should produce no error")
	assert.Nil(t, test.statsData[mockKey], "(stats type) Deleted data should return nil")

	err = test.DeleteData(mockKey, notSupported)
	assert.Error(t, err, "(not supported type) Delete data should produce an error")
}

func TestDeleteAllData(t *testing.T) {
	test := createInternalStorage()
	statsType := "stats"
	notSupported := "notSupported"

	err := test.DeleteAllData(statsType)
	assert.NoError(t, err, "(stats type) Delete all data should produce no error")
	assert.Len(t, test.statsData, 0, "(stats type) Stats storage should be empty")

	err = test.DeleteAllData(notSupported)
	assert.Error(t, err, "(not supported type) Delete all data should produce an error")
}

func createInternalStorage() internal {
	test := internal{}
	test.Initialize()
	test.statsData["default"] = []byte("default")
	return test
}
