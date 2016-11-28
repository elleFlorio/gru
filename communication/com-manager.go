package communication

import (
	"errors"
	"math/rand"
	"time"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/cluster"
	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/data"
	"github.com/elleFlorio/gru/network"
)

// api
const c_ROUTE_SHARED string = "/gru/v1/shared"

var (
	ErrInvalidFriendsNumber error = errors.New("Friends number should be > 0")
	ErrNoFriends            error = errors.New("There are no friends to reach")
)

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
}

func UpdateFriends() {
	var err error
	nFriends := cfg.GetAgentCommunication().MaxFriends

	log.Debugln("Updating friends...")

	peers := getAllPeers()
	log.WithField("peers", len(peers)).Debugln("Number of peers")
	if len(peers) == 0 {
		log.Warnln("There are no peers to reach")
		return
	}

	friends, err := chooseRandomFriends(peers, nFriends)
	if err != nil {
		log.WithField("err", err).Warnln("Error choosing friends to update")
		return
	}
	log.WithField("friends", friends).Debugln("Friends to connect with")

	sendDataToFriends(friends)
}

func getAllPeers() map[string]string {
	myCluster, err := cluster.GetMyCluster()
	if err != nil {
		log.WithField("err", err).Errorln("Error getting my cluster")
		return map[string]string{}
	}
	peers := cluster.ListNodes(myCluster.Name, true)
	delete(peers, cfg.GetNodeConfig().Name)
	return peers
}

// Is there a more efficient way to do this?
func chooseRandomFriends(peers map[string]string, n int) (map[string]string, error) {
	nPeers := len(peers)
	if nPeers < 1 {
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

	indexes := rand.Perm(nPeers)[:n]
	for _, index := range indexes {
		friendKey := peersKeys[index]
		friends[friendKey] = peers[friendKey]
	}

	return friends, nil
}

func sendDataToFriends(friends map[string]string) {
	var err error
	local, err := data.GetSharedLocal()
	if err != nil {
		log.WithField("err", err).Errorln("Cannot get local shared data to send")
		return
	}

	encoded := data.SharedToByte(local)

	for _, address := range friends {
		friendRoute := address + c_ROUTE_SHARED
		_, err := network.DoRequest("POST", friendRoute, encoded)
		if err != nil {
			log.WithField("address", address).Warnln("Error sending data to friends")
		}
	}
}
