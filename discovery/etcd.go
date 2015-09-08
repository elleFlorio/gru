package discovery

import (
	log "github.com/Sirupsen/logrus"
	"github.com/coreos/go-etcd/etcd"
)

// make these parameters?
const sortResults = false
const recursiveResults = false

type EtcdDiscovery struct {
	uuid   string
	client *etcd.Client
}

func (p *EtcdDiscovery) Name() string {
	return "etcd"
}

func (p *EtcdDiscovery) Initialize(uuid string, uri string) error {
	p.uuid = uuid
	p.client = etcd.NewClient([]string{uri})
	return nil
}

func (p *EtcdDiscovery) Register(myAddress string, ttl uint64) error {
	path := "/nodes/" + p.uuid

	_, err := p.client.Set(path, myAddress, ttl)
	if err != nil {
		log.WithField("error", err).Errorln("Registering to discovery service")
		return err
	}

	return err
}

func (p *EtcdDiscovery) Get(key string) (map[string]string, error) {
	result := make(map[string]string)

	resp, err := p.client.Get(key, sortResults, recursiveResults)
	if err != nil {
		log.WithField("error", err).Errorln("Querying discovery service")
		return nil, err
	}

	for _, entry := range resp.Node.Nodes {
		result[entry.Key] = entry.Value
	}

	return result, nil
}

func (p *EtcdDiscovery) Set(key string, value string, ttl uint64) error {
	_, err := p.client.Set(key, value, ttl)
	if err != nil {
		log.WithField("error", err).Errorln("Setting value to discovery service")
		return err
	}

	return err
}
