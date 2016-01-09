package communication

import (
	"math/rand"
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/node"
	"github.com/elleFlorio/gru/storage"
	"github.com/elleFlorio/gru/utils"
)

func TestClearFriendsData(t *testing.T) {
	var err error
	storage.New("internal")
	data := []byte{}
	storage.StoreLocalData(data, enum.ANALYTICS)
	storage.StoreClusterData(data, enum.ANALYTICS)
	storage.StoreData("node1", data, enum.ANALYTICS)
	storage.StoreData("node2", data, enum.ANALYTICS)
	storage.StoreData("node3", data, enum.ANALYTICS)

	err = clearFriendsData()
	assert.NoError(t, err)
	stored, _ := storage.GetAllData(enum.ANALYTICS)
	assert.Equal(t, 2, len(stored))
}

func TestChooseRandomFriends(t *testing.T) {
	cfg.SetNode(node.CreateMockNode())
	mockPeers := createMockPeers(100)
	nFriends := 10
	test, err := chooseRandomFriends(mockPeers, nFriends)
	assert.NoError(t, err, "(nFrineds > 0) Choose friends should produce no error")

	nFriends = 150
	test, err = chooseRandomFriends(mockPeers, nFriends)
	friendsKeys := make([]string, 0, len(mockPeers)-1)
	for key, _ := range test {
		friendsKeys = append(friendsKeys, key)
	}

	assert.NoError(t, err, "(nFrineds > nPeers) Choose friends should produce no error")
	assert.Len(t, test, len(mockPeers)-1, "(nFrineds > nPeers) Choose peers should return the map of all peers except me")
	assert.NotContains(t, friendsKeys, cfg.GetNodeConfig().Name, "(nFrineds > nPeers) Choose friends should not contain my key")

	nFriends = 0
	test, err = chooseRandomFriends(mockPeers, nFriends)
	assert.Error(t, err, "(nFriends < 0) Choose friends should produce an error")

	mockPeers = createMockPeers(0)
	nFriends = 10
	_, err = chooseRandomFriends(mockPeers, nFriends)
	assert.Error(t, err, "(peers = 0) Choose friend should produce an error")

	mockPeers = createMockPeers(1)
	nFriends = 10
	_, err = chooseRandomFriends(mockPeers, nFriends)
	assert.Error(t, err, "(peers = me) Choose friend should produce an error")

	mockPeers = createMockPeers(2)
	nFriends = 10
	test, err = chooseRandomFriends(mockPeers, nFriends)
	assert.NoError(t, err, "(nFrineds == 2) Choose friends should produce no error")
	assert.Len(t, test, len(mockPeers)-1, "(nFrineds == 2) Choose peers should return the map of all peers except me")
	assert.NotContains(t, friendsKeys, cfg.GetNodeConfig().UUID, "(nFrineds == 2) Choose friends should not contain my key")

}

func createMockPeers(nPeers int) map[string]string {
	myMockName := node.CreateMockNode().Configuration.Name
	mockPeers := make(map[string]string, nPeers)
	for i := 0; i < nPeers-1; i++ {
		name := utils.GetRandomName(0)
		mockPeers[name] = string(rand.Intn(nPeers))
	}
	mockPeers[myMockName] = "myValue"
	return mockPeers
}
