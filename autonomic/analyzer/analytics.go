package analyzer

import (
	"encoding/json"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/service"
)

type GruAnalytics struct {
	Service map[string]ServiceAnalytics `json:"service"`
	System  SystemAnalytics             `json:"system"`
	Cluster ClusterAnalytics            `json:"cluster"`
	Health  enum.Label                  `json:"health"`
}

type ServiceAnalytics struct {
	Load      enum.Label             `json:"load"`
	Resources ResourcesAnalytics     `json:"resources"`
	Instances service.InstanceStatus `json:"instances"`
	Health    enum.Label             `json:"health"`
}

type ResourcesAnalytics struct {
	Cpu       enum.Label `json:"cpu"`
	Memory    enum.Label `json:"memory"`
	Available enum.Label `json:"available"`
}

type SystemAnalytics struct {
	Services  []string               `json:"services"`
	Resources ResourcesAnalytics     `json:"resources"`
	Instances service.InstanceStatus `json:"instances"`
	Health    enum.Label             `json:"health"`
}

type ClusterAnalytics struct {
	Services           []string           `json:"services"`
	ResourcesAnalytics ResourcesAnalytics `json:"resources"`
	Health             enum.Label         `json:"health"`
}

func convertAnalyticsToData(analytics GruAnalytics) ([]byte, error) {
	data, err := json.Marshal(analytics)
	if err != nil {
		log.WithField("error", err).Errorln("Error marshaling analytics data")
		return nil, err
	}

	return data, nil
}

func ConvertDataToAnalytics(data []byte) (GruAnalytics, error) {
	analytics := GruAnalytics{}
	err := json.Unmarshal(data, &analytics)
	if err != nil {
		log.WithField("error", err).Errorln("Error unmarshaling analytics data")
	}

	return analytics, err
}
