package data

type Info struct {
	Service map[string]ServiceInfo `json:"service"`
	System  SystemInfo             `json:"system"`
}

type ServiceInfo struct {
	Load      float64 `json:"load"`
	Cpu       float64 `json:"cpu"`
	Memory    float64 `json:"memory"`
	Resources float64 `json:"resources"`
	Active    bool    `json:"active"`
}

type SystemInfo struct {
	Cpu            float64  `json:"cpu"`
	Memory         float64  `json:"memory"`
	Health         float64  `json:"health"`
	ActiveServices []string `json:"activeservices"`
}
