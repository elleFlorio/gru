package service

import (
	"time"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	ch "github.com/elleFlorio/gru/channels"
	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/discovery"
)

var discoveryConf = cfg.GetAgentDiscovery()

func RegisterServiceInstanceId(name string, id string, address string) {
	var err error
	opt := discovery.Options{
		"TTL": time.Duration(discoveryConf.TTL) * time.Second,
	}

	isntanceKey := discoveryConf.AppRoot + "/" + name + "/" + id
	instanceValue := address

	err = discovery.Set(isntanceKey, instanceValue, opt)
	if err != nil {
		log.WithFields(log.Fields{
			"service":  name,
			"instance": id,
			"address":  address,
			"err":      err,
		}).Errorln("Error registering service instance to discovery service")
	}
}

func UnregisterServiceInstance(name string, id string) {
	isntanceKey := discoveryConf.AppRoot + "/" + name + "/" + id
	discovery.Delete(isntanceKey)
	ch_stop, err := ch.GetInstanceChannel(id)
	if err != nil {
		log.WithFields(log.Fields{
			"service":  name,
			"instance": id,
			"err":      err,
		}).Errorln("Error getting instance stop channel")
		return
	}
	ch_stop <- struct{}{}
}

func KeepAlive(name string, id string, address string) {
	go keepAlive(name, id, address)
}

func keepAlive(name string, id string, address string) {
	var err error
	ticker := time.NewTicker(time.Duration(discoveryConf.TTL-1) * time.Second)
	opt := discovery.Options{
		"TTL": time.Duration(discoveryConf.TTL) * time.Second,
	}

	ch_stop := ch.CreateInstanceChannel(id)

	isntanceKey := discoveryConf.AppRoot + "/" + name + "/" + id
	instanceValue := address

	for {
		select {
		case <-ticker.C:
			err = discovery.Set(isntanceKey, instanceValue, opt)
			if err != nil {
				log.WithFields(log.Fields{
					"service":  name,
					"instance": id,
					"address":  address,
					"err":      err,
				}).Errorln("Error keeping instance alive")
			}
		case <-ch_stop:
			return
		}
	}
}
