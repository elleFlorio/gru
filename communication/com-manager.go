package communication

import (
	"errors"
	"math/rand"
	"time"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/cluster"
	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/data"
	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/network"
	"github.com/elleFlorio/gru/storage"
)

// api
const c_ROUTE_SHARED string = "/gru/v1/shared"

var (
	ErrInvalidFriendsNumber error = errors.New("Friends number should be > 0")
	ErrInvalidDataType      error = errors.New("Invalid data type to retrieve")
	ErrNoPeers              error = errors.New("There are no peers to reach")
	ErrNoFriends            error = errors.New("There are no friends to reach")
)

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
}

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
	var err error
	log.Debugln("Updating friends data")
	err = clearFriendsData()
	if err != nil {
		return err
	}

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

	friendsData, err := getFriendsData(friends)
	if err != nil {
		return err
	}

	err = mergeFriendsData(friendsData)
	if err != nil {
		return err
	}

	return nil
}

func clearFriendsData() error {
	allData, err := storage.GetAllData(enum.SHARED)
	if err != nil {
		log.Errorln("Error cleaning up friends data")
		return err
	}

	for key, _ := range allData {
		if key != enum.CLUSTER.ToString() && key != enum.LOCAL.ToString() {
			log.WithField("friend", key).Debugln("cleared friend data")
			storage.DeleteData(key, enum.SHARED)
		}
	}

	return nil
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

func getFriendsData(friends map[string]string) ([]data.Shared, error) {
	var err error
	friendsData := make([]data.Shared, 0, len(friends))

	for friend, address := range friends {
		friendRoute := address + c_ROUTE_SHARED
		friendData, err := network.DoRequest("GET", friendRoute, nil)
		if err != nil {
			log.WithField("address", address).Debugln("Error retrieving friend stats")
		}
		err = storage.StoreData(friend, friendData, enum.SHARED)
		if err != nil {
			log.WithField("err", err).Debugln("Error storing friends data")
		}

		sharedData, err := data.ByteToShared(friendData)
		if err != nil {
			log.WithField("address", address).Debugln("Friend data not stored")
		} else {
			friendsData = append(friendsData, sharedData)
		}

	}

	return friendsData, err
}

func mergeFriendsData(friendsData []data.Shared) error {
	merged, err := data.MergeShared(friendsData)
	if err != nil {
		log.WithField("err", err).Debugln("Cannot store friends data")
		return err
	}

	data.SaveShared(merged)

	return nil
}
