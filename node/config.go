package node

type Node struct {
	Configuration Config      `json:"configuration"`
	Constraints   Constraints `json:"constraints"`
	Resources     Resources   `json:resources`
	Active        bool        `json:"active"`
}

type Config struct {
	UUID    string `json:"uuid"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Cluster string `json:"cluster"`
}

// Is this still necessary?
type Constraints struct {
	CpuMin       float64  `json:"cpumin"`
	CpuMax       float64  `json:"cpumax"`
	BaseServices []string `json:"baseservices"`
}

type Resources struct {
	TotalMemory int64 `json:"totalmemory"`
	TotalCpus   int64 `json:"totalcpus"`
	UsedMemory  int64 `json:"usedmemory"`
	UsedCpu     int64 `json:"usedcpu"`
}
