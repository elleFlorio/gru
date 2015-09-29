package node

import (
	"io/ioutil"
)

func CreateMockNode() Node {
	mockConstraints := Constraints{
		CpuMin:       0.2,
		CpuMax:       0.8,
		MaxInstances: 10, //TODO this will be removed
	}
	mockNode := Node{
		UUID:        "abcdefghi",
		Constraints: mockConstraints,
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
