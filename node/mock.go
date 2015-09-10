package node

import (
	"io/ioutil"
)

func CreateMockNode() Node {
	mockConstraints := Constraints{
		CpuMin:       0.2,
		CpuMax:       0.8,
		MaxInstances: 10,
	}
	mockNode := Node{
		UUID:        "abcdefghi",
		Name:        "mockNode",
		IpAddr:      "127.0.0.1",
		Port:        "8080",
		Constraints: mockConstraints,
	}

	return mockNode
}

func createMockConfigFileNode(ipAddr bool, port bool) string {
	mockConfigFileNode := `{
		"Name":"mockNode",
		"Constraints":{
			"CpuMax":0.8,
			"CpuMin":0.2,
			"MaxInstances":10
		}
	}`

	if ipAddr && port {
		mockConfigFileNode = `{
		"Name":"mockNode",
		"IpAddr":"127.0.0.1",
		"Port":"8080",
		"Constraints":{
			"CpuMax":0.8,
			"CpuMin":0.2,
			"MaxInstances":10
			}
		}`
	} else if ipAddr {
		mockConfigFileNode = `{
		"Name":"mockNode",
		"IpAddr":"127.0.0.1",
		"Constraints":{
			"CpuMax":0.8,
			"CpuMin":0.2,
			"MaxInstances":10
			}
		}`
	} else if port {
		mockConfigFileNode = `{
		"Name":"mockNode",
		"Port":"8080",
		"Constraints":{
			"CpuMax":0.8,
			"CpuMin":0.2,
			"MaxInstances":10
			}
		}`
	}

	tmpfile, err := ioutil.TempFile(".", "gru_test_node_config")
	if err != nil {
		panic(err)
	}

	ioutil.WriteFile(tmpfile.Name(), []byte(mockConfigFileNode), 0644)

	return tmpfile.Name()

}
