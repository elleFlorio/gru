package communication

import (
	"errors"
	"math/rand"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/discovery"
	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/network"
	"github.com/elleFlorio/gru/node"
	"github.com/elleFlorio/gru/storage"
)

//TODO get these in an automatic way from api
const c_PATH_NODES string = "/nodes/"
const c_ROUTE_ANALYTICS string = "/gru/v1/analytics"

var (
	ErrInvalidFriendsNumber error = errors.New("Friends number should be > 0")
	ErrInvalidDataType      error = errors.New("Invalid data type to retrieve")
	ErrNoPeers              error = errors.New("There are no peers to reach")
	ErrNoFriends            error = errors.New("There are no friends to reach")
)

func KeepAlive(ttl int) {
	if discovery.Name() != "noservice" {
		err := discovery.Set(createAgentKey(), createAgentAddress(), ttl)
		if err != nil {
			log.WithField("err", err).Errorln("Cannot keep the agent alive")
		}
	}
}

func createAgentAddress() string {
	return "http://" + network.Config().IpAddress + ":" + network.Config().Port
}

func createAgentKey() string {
	return c_PATH_NODES + node.Config().UUID
}

func UpdateFriendsData(nFriends int) error {
	log.Debugln("Updating friends data")
	storage.DeleteAllData(enum.ANALYTICS)

	peers, err := getAllPeers()
	if err != nil {
		return err
	}

	log.WithField("peers", len(peers)).Debugln("Number of peers")

	if len(peers) == 0 {
		return ErrNoPeers
	}

	friends, err := chooseRandomFriends(peers, nFriends)
	if err != nil {
		return err
	}

	log.WithField("friends", friends).Debugln("Friends to connect with")

	err = getFriendsData(friends)
	if err != nil {
		return err
	}

	return nil
}

func getAllPeers() (map[string]string, error) {
	peers, err := discovery.Get(c_PATH_NODES)
	if err != nil {
		return nil, err
	}

	return peers, nil
}

// Is there a more efficient way to do this?
func chooseRandomFriends(peers map[string]string, n int) (map[string]string, error) {
	nPeers := len(peers)
	if nPeers <= 1 {
		return nil, ErrNoFriends
	}

	if n <= 0 {
		return nil, ErrInvalidFriendsNumber
	} else if n > nPeers {
		n = nPeers
	}

	friends := make(map[string]string, n)

	peersKeys := make([]string, 0, len(peers))
	for peerKey, _ := range peers {
		peersKeys = append(peersKeys, peerKey)
	}

	friendsKeys := make([]string, 0, n)
	indexes := rand.Perm(nPeers)[:n]
	for _, index := range indexes {
		if peersKeys[index] != node.Config().UUID {
			friendsKeys = append(friendsKeys, peersKeys[index])
		}
	}

	for _, friendKey := range friendsKeys {
		friends[friendKey] = peers[friendKey]
	}

	return friends, nil
}

func getFriendsData(friends map[string]string) error {
	var err error

	for friend, address := range friends {
		// This should not be possible, but in a local test (multiple nodes on the same node with
		//different ports) it happened.
		if address != createAgentAddress() {
			friendRoute := address + c_ROUTE_ANALYTICS
			friendData, err := network.DoRequest("GET", friendRoute, nil)
			if err != nil {
				log.WithField("address", address).Debugln("Error retrieving friend stats")
			}
			err = storage.StoreData(friend, friendData, enum.ANALYTICS)
			if err != nil {
				log.WithField("err", err).Debugln("Error storing friend stats")
			}
		}
	}

	return err
}
