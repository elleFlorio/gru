package analyzer

import (
	"github.com/elleFlorio/gru/service"
)

type GruAnalytics struct {
	Service map[string]ServiceAnalytics `json:"service"`
	System  SystemAnalytics             `json:"system"`
	Cluster ClusterAnalytics            `json:"cluster"`
	Health  float64                     `json:"health"`
}

type ServiceAnalytics struct {
	Load      float64                `json:"load"`
	Resources ResourcesAnalytics     `json:"resources"`
	Instances service.InstanceStatus `json:"instances"`
	Health    float64                `json:"health"`
}

type ResourcesAnalytics struct {
	Cpu       float64 `json:"cpu"`
	Memory    float64 `json:"memory"`
	Available float64 `json:"available"`
}

type SystemAnalytics struct {
	Services  []string               `json:"services"`
	Resources ResourcesAnalytics     `json:"resources"`
	Instances service.InstanceStatus `json:"instances"`
	Health    float64                `json:"health"`
}

type ClusterAnalytics struct {
	Services           []string           `json:"services"`
	ResourcesAnalytics ResourcesAnalytics `json:"resources"`
	Health             float64            `json:"health"`
}
