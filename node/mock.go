package node

import (
	"io/ioutil"
)

func CreateMockNode() Node {
	mockConstraints := Constraints{
		CpuMin: 0.2,
		CpuMax: 0.8,
	}

	mockResources := Resources{
		TotalMemory: int64(8 * 1024 * 1024 * 1024),
		TotalCpus:   int64(4),
		UsedMemory:  int64(0),
		UsedCpu:     int64(0),
	}

	mockNode := Node{
		UUID:        "abcdefghi",
		Constraints: mockConstraints,
		Resources:   mockResources,
	}

	return mockNode
}

func createMockConfigFileNode() string {
	mockConfigFileNode := `{
		"Name":"mockNode",
		"Constraints":{
			"CpuMax":0.8,
			"CpuMin":0.2,
			"MaxInstances":10
		}
	}`

	tmpfile, err := ioutil.TempFile(".", "gru_test_node_config")
	if err != nil {
		panic(err)
	}

	ioutil.WriteFile(tmpfile.Name(), []byte(mockConfigFileNode), 0644)

	return tmpfile.Name()

}
