package configuration

type Node struct {
	Configuration NodeConfig      `json:"configuration"`
	Constraints   NodeConstraints `json:"constraints"`
	Resources     NodeResources   `json:"resources"`
	Instances     ServiceStatus   `json:"instances"`
	Active        bool            `json:"active"`
}

type NodeConfig struct {
	UUID    string `json:"uuid"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Cluster string `json:"cluster"`
	Remote  string `json:"remote"`
}

// Is this still necessary?
type NodeConstraints struct {
	CpuMin       float64  `json:"cpumin"`
	CpuMax       float64  `json:"cpumax"`
	BaseServices []string `json:"baseservices"`
}

type NodeResources struct {
	TotalCpus   int64 `json:"totalcpus"`
	TotalMemory int64 `json:"totalmemory"`
}
