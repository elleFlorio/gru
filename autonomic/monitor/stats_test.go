package monitor

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/elleFlorio/gru/node"
)

//TODO ok, I know it is not a very good test...
func TestMergeStats(t *testing.T) {
	node.UpdateNode(node.CreateMockNode())
	stats1 := CreateMockStats()
	stats2 := CreateMockStats()
	stats3 := CreateMockStats()

	mockStats := map[string]GruStats{
		"stats1":           stats1,
		"stats2":           stats2,
		node.Config().UUID: stats3,
	}

	test := mergeStats(mockStats)
	assert.NotEmpty(t, test, "Merge stats should produce a non empty map")
}
