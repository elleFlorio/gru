package configuration

type Node struct {
	Configuration NodeConfig      `json:"configuration"`
	Constraints   NodeConstraints `json:"constraints"`
	Resources     NodeResources   `json:resources`
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
	TotalMemory int64 `json:"totalmemory"`
	TotalCpus   int64 `json:"totalcpus"`
	UsedMemory  int64 `json:"usedmemory"`
	UsedCpu     int64 `json:"usedcpu"`
}
