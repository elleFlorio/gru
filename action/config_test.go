package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
testConfig.json

{
	"Services":
	{
		"MinActive":1,
		"MaxActive":10
		"Service":[
			{
				"service1":3,
				"service2":2,
			}
		]
	},

	"Nodes": {
		"MaxInstances":10
	},

	"Cpu": {
		"Min":0.3,
		"Max":0.8
	},

	"Docker": {
		"StopTimeout":3
	}
}
*/

func TestBuildGruActionConfig(t *testing.T) {
	correctPath := "testConfig.json"
	wrongPath := ""

	_, err := BuildGruActionConfig(wrongPath)

	//Test Error in case of wrong path
	assert.Error(t, err, "Wrong path return a error")

	result, err := BuildGruActionConfig(correctPath)

	//Test correct reading and parsing
	assert.NoError(t, err, "The function should return no error")

	//Test a single field
	assert.Equal(t, 1, result.Services.MinActive, "result.Services.MinActive should be 1")

	//Test map creation
	assert.Equal(t, 3, result.Services.Service["service1"], "service 1 should have a limit of 3")

	//Test the empty map
	assert.Empty(t, result.Nodes.Node, "Node map should be empty")
}
