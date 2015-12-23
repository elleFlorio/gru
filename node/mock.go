package node

import (
	cfg "github.com/elleFlorio/gru/configuration"
)

func CreateMockNode() cfg.Node {
	mockConfiguration := cfg.NodeConfig{
		Name:    "topolino",
		UUID:    "abcdefghi",
		Address: "http://localhost:8080",
	}

	mockConstraints := cfg.NodeConstraints{
		CpuMin: 0.2,
		CpuMax: 0.8,
	}

	mockResources := cfg.NodeResources{
		TotalCpus:   int64(4),
		TotalMemory: int64(8 * 1024 * 1024 * 1024),
	}

	mockNode := cfg.Node{
		Configuration: mockConfiguration,
		Constraints:   mockConstraints,
		Resources:     mockResources,
		Active:        false,
	}

	return mockNode
}
