package resources

type Resource struct {
	CPU    CpuResource    `json:"cpu"`
	Memory MemoryResource `json:"memory"`
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
