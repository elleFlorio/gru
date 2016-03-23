package data

import (
	"encoding/json"
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/service"
	"github.com/elleFlorio/gru/storage"
)

func init() {
	//Initialize storage
	storage.New("internal")
	//simply check if the functions return without errors
	CreateMockHistory()
	StoreRandomMockPolicy()
	ListMockServices()
	MaxNumberOfEntryInHistory()
}

func TestSaveStats(t *testing.T) {
	defer storage.DeleteAllData(enum.STATS)
	stats := CreateMockStats()
	SaveStats(stats)

	encoded, _ := storage.GetLocalData(enum.STATS)
	decoded := GruStats{}
	json.Unmarshal(encoded, &decoded)

	assert.Equal(t, stats, decoded)
}

func TestSaveAnalytics(t *testing.T) {
	defer storage.DeleteAllData(enum.ANALYTICS)
	analytics := CreateMockAnalytics()
	SaveAnalytics(analytics)

	encoded, _ := storage.GetLocalData(enum.ANALYTICS)
	decoded := GruAnalytics{}
	json.Unmarshal(encoded, &decoded)

	assert.Equal(t, analytics, decoded)
}

func TestSavePolicy(t *testing.T) {
	defer storage.DeleteAllData(enum.POLICIES)
	policy := CreateRandomMockPolicies(1)[0]
	SavePolicy(policy)

	encoded, _ := storage.GetLocalData(enum.POLICIES)
	decoded := Policy{}
	json.Unmarshal(encoded, &decoded)

	assert.Equal(t, policy, decoded)
}

func TestSaveShared(t *testing.T) {
	defer storage.DeleteAllData(enum.SHARED)
	info := CreateMockInfo()
	SaveShared(info)

	encoded, _ := storage.GetClusterData(enum.SHARED)
	decoded := Shared{}
	json.Unmarshal(encoded, &decoded)

	assert.Equal(t, info, decoded)
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
	var encoded []byte
	var decoded Shared
	var err error

	info := CreateMockInfo()
	encoded, _ = json.Marshal(info)
	decoded, err = ByteToShared(encoded)
	assert.NoError(t, err)
	assert.Equal(t, info, decoded)

	bad := 0
	encoded, _ = json.Marshal(bad)
	decoded, err = ByteToShared(encoded)
	assert.Error(t, err)
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

func TestGeInfo(t *testing.T) {
	defer storage.DeleteAllData(enum.SHARED)
	var err error

	_, err = GetShared()
	assert.Error(t, err)

	expected := CreateMockInfo()
	SaveShared(expected)
	info, err := GetShared()
	assert.NoError(t, err)
	assert.Equal(t, expected, info)
}

func TestCheckAndAppend(t *testing.T) {
	list0 := []string{}
	list1 := []string{"pippo", "topolino"}
	list2 := []string{"paperino"}
	list3 := []string{"topolino", "paperino", "paperone"}

	list0 = checkAndAppend(list0, list1)
	assert.Len(t, list0, 2)
	assert.Contains(t, list0, "pippo")
	assert.Contains(t, list0, "topolino")

	list0 = checkAndAppend(list0, list2)
	assert.Len(t, list0, 3)
	assert.Contains(t, list0, "paperino")

	list0 = checkAndAppend(list0, list3)
	assert.Len(t, list0, 4)
	assert.Contains(t, list0, "paperone")
}

func TestMergeInfo(t *testing.T) {
	defer service.ClearMockServices()
	var err error

	service.SetMockServices()
	info1 := CreateMockInfo()
	info2 := CreateMockInfo()
	info3 := CreateMockInfo()

	peers := []Shared{info1, info2, info3}
	merged, err := MergeShared(peers)
	assert.NoError(t, err)
	assert.Equal(t, info1.Service["service1"].Cpu, merged.Service["service1"].Cpu)
	assert.Equal(t, info1.System.ActiveServices, merged.System.ActiveServices)

	empty := []Shared{}
	_, err = MergeShared(empty)
	assert.Error(t, err)

	one := []Shared{info1}
	_, err = MergeShared(one)
	assert.NoError(t, err)
}
