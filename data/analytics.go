package data

import (
	cfg "github.com/elleFlorio/gru/configuration"
)

type GruAnalytics struct {
	Service map[string]AnalyticData `json:"service"`
	System  AnalyticData            `json:"system"`
}

type ServiceAnalytics struct {
	Load      float64            `json:"load"`
	Resources ResourcesAnalytics `json:"resources"`
	Instances cfg.ServiceStatus  `json:"instances"`
}

type ResourcesAnalytics struct {
	Cpu       float64 `json:"cpu"`
	Memory    float64 `json:"memory"`
	Available float64 `json:"available"`
}

type SystemAnalytics struct {
	Services  []string           `json:"services"`
	Resources ResourcesAnalytics `json:"resources"`
	Instances cfg.ServiceStatus  `json:"instances"`
}

type AnalyticData struct {
	BaseAnalytics map[string]float64
	UserAnalytics map[string]float64
}
