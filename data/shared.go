package data

type Shared struct {
	Service map[string]ServiceShared `json:"service"`
	System  SystemShared             `json:"system"`
}

type ServiceShared struct {
	Load      float64 `json:"load"`
	Cpu       float64 `json:"cpu"`
	Memory    float64 `json:"memory"`
	Resources float64 `json:"resources"`
	Active    bool    `json:"active"`
}

type SystemShared struct {
	Cpu            float64  `json:"cpu"`
	Memory         float64  `json:"memory"`
	Health         float64  `json:"health"`
	ActiveServices []string `json:"activeservices"`
}
