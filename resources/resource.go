package resources

type Resource struct {
	CPU     CpuResource     `json:"cpu"`
	Memory  MemoryResource  `json:"memory"`
	Network NetworkResource `json:"network"`
}

type CpuResource struct {
	Cores map[int]bool `json:"cores"`
	Total int64        `json:"total"`
	Used  int64        `json:"used"`
}

type MemoryResource struct {
	Total int64 `json:"total"`
	Used  int64 `json:"used"`
}

type NetworkResource struct {
	ServicePorts map[string]Ports `json:"ports"`
}

type Ports struct {
	Available    []string `json:"available"`
	Occupied     []string `json:"occupied"`
	LastAssigned string   `json:lastassigned`
}
