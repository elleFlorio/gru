package communication

import (
	"errors"
	"math/rand"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/discovery"
	"github.com/elleFlorio/gru/network"
	"github.com/elleFlorio/gru/node"
	"github.com/elleFlorio/gru/storage"
)

//TODO get these in an automatic way?
const path string = "/nodes/"
const routeStats string = "/gru/v1/stats"

var (
	ErrInvalidFriendsNumber error = errors.New("Friends number should be > 0")
	ErrInvalidDataType      error = errors.New("Invalid data type to retrieve")
	ErrNoPeers              error = errors.New("There are no peers to reach")
	ErrNoFriends            error = errors.New("There are no friends to reach")
)

func KeepAlive(ttl int) {
	agentAddress := "http://" + network.Config().IpAddress + ":" + network.Config().Port
	key := path + node.Config().UUID
	err := discovery.Service().Set(key, agentAddress, ttl)
	if err != nil {
		log.WithField("error", err).Errorln("Keeping alive the agent")
	}
	log.WithFields(log.Fields{
		"key":     key,
		"address": agentAddress,
	}).Debugln("Agent is alive")
}

func UpdateFriendsData(nFriends int, dataType string) error {
	log.WithField("status", "start").Debugln("Updating friends data")
	defer log.WithField("status", "done").Debugln("Updating friends data")
	storage.DataStore().DeleteAllData(dataType)

	peers, err := getAllPeers()
	if err != nil {
		return err
	}

	if len(peers) == 0 {
		return ErrNoPeers
	}

	friends, err := chooseRandomFriends(peers, nFriends)
	if err != nil {
		return err
	}

	err = getFriendsData(friends, dataType)
	if err != nil {
		return err
	}

	return nil
}

func getAllPeers() (map[string]string, error) {
	peers, err := discovery.Service().Get(path)
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

func getFriendsData(friends map[string]string, dataType string) error {
	switch dataType {
	case "stats":
		for friend, address := range friends {
			// This should not be possible, but in a local test (multiple nodes on the same node with
			//different ports) it happened.
			myAddress := "http://" + network.Config().IpAddress + ":" + network.Config().Port
			if address != myAddress {
				friendRoute := address + routeStats
				friendData, err := network.DoRequest("GET", friendRoute, nil)
				if err != nil {
					log.WithField("address", address).Warnln("Error retrieving friend stats")
				}
				err = storage.DataStore().StoreData(friend, friendData, dataType)
				if err != nil {
					log.WithField("error", err).Warnln("Error storing friend stats")
				}
			}
		}
		return nil
	case "analytics":
		//TODO
	}
	return ErrInvalidDataType
}
