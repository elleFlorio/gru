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
	Status        map[string]PortStatus `json:"status"`
	LastRequested map[string]string     `json:lastrequested`
	LastAssigned  map[string][]string   `json:lastassigned`
}

type PortStatus struct {
	Available []string `json:"available"`
	Occupied  []string `json:"occupied"`
}
