package data

import (
	"encoding/json"
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/storage"
)

func init() {
	//Initialize storage
	storage.New("internal")
}

func TestSaveStats(t *testing.T) {
	defer storage.DeleteAllData(enum.STATS)
	stats := CreateMockStats()
	SaveStats(stats)

	encoded, _ := storage.GetClusterData(enum.STATS)
	decoded := GruStats{}
	json.Unmarshal(encoded, &decoded)

	assert.Equal(t, stats, decoded)
}

func TestSaveAnalytics(t *testing.T) {
	defer storage.DeleteAllData(enum.ANALYTICS)
	analytics := CreateMockAnalytics()
	SaveAnalytics(analytics)

	encoded, _ := storage.GetClusterData(enum.ANALYTICS)
	decoded := GruAnalytics{}
	json.Unmarshal(encoded, &decoded)

	assert.Equal(t, analytics, decoded)
}

func TestSavePolicy(t *testing.T) {
	defer storage.DeleteAllData(enum.POLICIES)
	policy := CreateRandomMockPolicies(1)[0]
	SavePolicy(policy)

	encoded, _ := storage.GetClusterData(enum.POLICIES)
	decoded := Policy{}
	json.Unmarshal(encoded, &decoded)

	assert.Equal(t, policy, decoded)
}

func TestSaveInfo(t *testing.T) {
	//TODO
}

func TestByteToStats(t *testing.T) {
	var encoded []byte
	var decoded GruStats
	var err error

	stats := CreateMockStats()
	encoded, _ = json.Marshal(stats)
	decoded, err = ByteToStats(encoded)
	assert.NoError(t, err)
	assert.Equal(t, stats, decoded)

	bad := 0
	encoded, _ = json.Marshal(bad)
	decoded, err = ByteToStats(encoded)
	assert.Error(t, err)

}

func TestByteToAnalytics(t *testing.T) {
	var encoded []byte
	var decoded GruAnalytics
	var err error

	analytics := CreateMockAnalytics()
	encoded, _ = json.Marshal(analytics)
	decoded, err = ByteToAnalytics(encoded)
	assert.NoError(t, err)
	assert.Equal(t, analytics, decoded)

	bad := 0
	encoded, _ = json.Marshal(bad)
	decoded, err = ByteToAnalytics(encoded)
	assert.Error(t, err)
}

func TestByteToPolicy(t *testing.T) {
	var encoded []byte
	var decoded Policy
	var err error

	policy := CreateRandomMockPolicies(1)[0]
	encoded, _ = json.Marshal(policy)
	decoded, err = ByteToPolicy(encoded)
	assert.NoError(t, err)
	assert.Equal(t, policy, decoded)

	bad := 0
	encoded, _ = json.Marshal(bad)
	decoded, err = ByteToPolicy(encoded)
	assert.Error(t, err)
}

func TestByteToInfo(t *testing.T) {
	//TODO
}

func TestGetStats(t *testing.T) {
	defer storage.DeleteAllData(enum.STATS)
	var err error

	_, err = GetStats()
	assert.Error(t, err)

	expected := CreateMockStats()
	StoreMockStats()
	stats, err := GetStats()
	assert.NoError(t, err)
	assert.Equal(t, expected, stats)
}

func TestGetAnalytics(t *testing.T) {
	defer storage.DeleteAllData(enum.ANALYTICS)
	var err error

	_, err = GetAnalytics()
	assert.Error(t, err)

	expected := CreateMockAnalytics()
	StoreMockAnalytics()
	analytics, err := GetAnalytics()
	assert.NoError(t, err)
	assert.Equal(t, expected, analytics)
}

func TestGetPolicy(t *testing.T) {
	defer storage.DeleteAllData(enum.POLICIES)
	var err error

	_, err = GetPolicy()
	assert.Error(t, err)

	expected := CreateRandomMockPolicies(1)[0]
	SavePolicy(expected)
	policy, err := GetPolicy()
	assert.NoError(t, err)
	assert.Equal(t, expected, policy)
}
