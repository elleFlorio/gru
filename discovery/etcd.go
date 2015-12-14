package discovery

import (
	"time"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/coreos/etcd/client"
	"github.com/elleFlorio/gru/Godeps/_workspace/src/golang.org/x/net/context"

	"github.com/elleFlorio/gru/utils"
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

func (p *etcdDiscovery) Register(nodePath string, nodeAddress string) error {

	_, err := p.kAPI.Set(
		context.Background(),
		nodePath,
		nodeAddress,
		nil,
	)
	if err != nil {
		log.WithField("err", err).Errorln("Registering to discovery service")
		return err
	}

	return err
}

func (p *etcdDiscovery) Get(key string, opt Options) (map[string]string, error) {
	var err error
	result := make(map[string]string)
	cli_opt := &client.GetOptions{}
	err = utils.FillStruct(cli_opt, opt)
	if err != nil {
		log.WithField("err", err).Errorln("Error reading discovery get options")
		return nil, err
	}

	resp, err := p.kAPI.Get(context.Background(), key, cli_opt)
	if err != nil {
		log.WithField("err", err).Errorln("Querying discovery service")
		return nil, err
	}

	result = exploreNode(resp.Node, result)

	// Clean empty results
	for k, v := range result {
		if k == key && v == "" {
			delete(result, k)
		}
	}

	return result, nil
}

func exploreNode(node *client.Node, result map[string]string) map[string]string {
	if len(node.Nodes) > 0 {
		for _, next := range node.Nodes {
			result = exploreNode(next, result)
		}
	} else {
		result[node.Key] = node.Value
		log.WithFields(log.Fields{
			"key":   node.Key,
			"value": node.Value,
		}).Debugln("Node entry")

	}
	return result
}

func (p *etcdDiscovery) Set(key string, value string, opt Options) error {
	var err error

	cli_opt := &client.SetOptions{}
	if _, ok := opt["PrevExist"]; ok {
		opt["PrevExist"] = client.PrevExist
	}
	err = utils.FillStruct(cli_opt, opt)
	if err != nil {
		log.WithField("err", err).Errorln("Error reading discovery set options")
		return err
	}

	_, err = p.kAPI.Set(
		context.Background(),
		key,
		value,
		cli_opt)
	if err != nil {
		log.WithField("err", err).Errorln("Error setting value to discovery service")
		return err
	}

	return err
}
