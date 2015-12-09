package node

type Node struct {
	UUID        string      `json:"uuid"`
	Name        string      `json:"name"`
	Address     string      `json:"address"`
	Active      bool        `json:"active"`
	Constraints Constraints `json:"constraints"`
	Resources   Resources   `json:resources`
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
