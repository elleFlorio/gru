package data

type Shared struct {
	Service map[string]ServiceShared `json:"service"`
	System  SystemShared             `json:"system"`
}

type ServiceShared struct {
	Data   SharedData
	Active bool `json:"active"`
}

type SystemShared struct {
	Data           SharedData
	ActiveServices []string `json:"activeservices"`
}

type SharedData struct {
	BaseShared map[string]float64
	UserShared map[string]float64
}
