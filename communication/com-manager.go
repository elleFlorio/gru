package communication

import (
	"errors"
	"math/rand"
	"time"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/cluster"
	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/network"
	"github.com/elleFlorio/gru/storage"
)

// api
const c_ROUTE_ANALYTICS string = "/gru/v1/analytics"

var (
	ErrInvalidFriendsNumber error = errors.New("Friends number should be > 0")
	ErrInvalidDataType      error = errors.New("Invalid data type to retrieve")
	ErrNoPeers              error = errors.New("There are no peers to reach")
	ErrNoFriends            error = errors.New("There are no friends to reach")
)

func Start(nFriends int, timeInterval int) {
	go start(nFriends, timeInterval)
}

func start(nFriends int, timeInterval int) {
	ticker := time.NewTicker(time.Duration(timeInterval) * time.Second)
	for {
		select {
		case <-ticker.C:
			err := updateFriendsData(nFriends)
			if err != nil {
				log.WithField("err", err).Debugln("Cannot update friends data")
			}
		}
	}
}

func updateFriendsData(nFriends int) error {
	log.Debugln("Updating friends data")
	storage.DeleteAllData(enum.ANALYTICS)

	peers := getAllPeers()
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

func getAllPeers() map[string]string {
	myCluster, err := cluster.GetMyCluster()
	if err != nil {
		log.WithField("err", err).Errorln("Error getting my cluster")
		return map[string]string{}
	}
	return cluster.ListNodes(myCluster.Name, true)
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
		if peersKeys[index] != cfg.GetNodeConfig().Name {
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
		friendRoute := address + c_ROUTE_ANALYTICS
		friendData, err := network.DoRequest("GET", friendRoute, nil)
		if err != nil {
			log.WithField("address", address).Debugln("Error retrieving friend stats")
		}
		err = storage.StoreData(friend, friendData, enum.ANALYTICS)
		if err != nil {
			log.WithField("err", err).Debugln("Error storing friends data")
		}
	}

	return err
}
