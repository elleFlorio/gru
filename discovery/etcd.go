package discovery

import (
	"time"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/coreos/etcd/client"
	"github.com/elleFlorio/gru/Godeps/_workspace/src/golang.org/x/net/context"
)

type etcdDiscovery struct {
	kAPI client.KeysAPI
}

func (p *etcdDiscovery) Name() string {
	return "etcd"
}

func (p *etcdDiscovery) Initialize(uri string) error {
	log.WithField("uri", uri).Debugln("Trying to connect to etcd")

	cfg := client.Config{
		Endpoints: []string{uri},
	}

	etcd, err := client.New(cfg)
	if err != nil {
		return err
	}

	p.kAPI = client.NewKeysAPI(etcd)

	//This is needed to probe if the etcd server is reachable
	_, err = p.kAPI.Set(
		context.Background(),
		"/probe",
		"etcd",
		&client.SetOptions{TTL: time.Duration(1) * time.Millisecond},
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *etcdDiscovery) Register(myUUID string, myAddress string, ttl int) error {
	path := "/nodes/" + myUUID

	_, err := p.kAPI.CreateInOrder(
		context.Background(),
		path,
		myAddress,
		&client.CreateInOrderOptions{TTL: time.Duration(ttl) * time.Second},
	)
	if err != nil {
		log.WithField("error", err).Errorln("Registering to discovery service")
		return err
	}

	return err
}

func (p *etcdDiscovery) Get(key string) (map[string]string, error) {
	result := make(map[string]string)

	resp, err := p.kAPI.Get(context.Background(), key, nil)
	if err != nil {
		log.WithField("error", err).Errorln("Querying discovery service")
		return nil, err
	}

	log.WithField("metadata", resp).Debugln("Get etcd")

	for _, entry := range resp.Node.Nodes {
		log.WithFields(log.Fields{
			"key":   entry.Key,
			"value": entry.Value,
		}).Debugln("Node entries")
		result[entry.Key] = entry.Value
	}

	return result, nil
}

func (p *etcdDiscovery) Set(key string, value string, ttl int) error {

	_, err := p.kAPI.Set(
		context.Background(),
		key,
		value,
		&client.SetOptions{TTL: time.Duration(ttl) * time.Second})
	if err != nil {
		log.WithField("error", err).Errorln("Setting value to discovery service")
		return err
	}

	return err
}
