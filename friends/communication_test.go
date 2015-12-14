package friends

import (
	"math/rand"
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	"github.com/elleFlorio/gru/node"
	"github.com/elleFlorio/gru/utils"
)

func TestChooseRandomFriends(t *testing.T) {
	node.UpdateNodeConfig(node.CreateMockNode())
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
	assert.NotContains(t, friendsKeys, node.Config().UUID, "(nFrineds > nPeers) Choose friends should not contain my key")

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
	assert.NotContains(t, friendsKeys, node.Config().UUID, "(nFrineds == 2) Choose friends should not contain my key")

}

func createMockPeers(nPeers int) map[string]string {
	myMockUUID := node.CreateMockNode().UUID
	mockPeers := make(map[string]string, nPeers)
	for i := 0; i < nPeers-1; i++ {
		uuid, _ := utils.GenerateUUID()
		mockPeers[uuid] = string(rand.Intn(nPeers))
	}
	mockPeers[myMockUUID] = "myValue"
	return mockPeers
}
