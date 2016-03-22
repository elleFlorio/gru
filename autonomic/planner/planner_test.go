package planner

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

func TestSetPlannerStrategy(t *testing.T) {
	probabilistic := "probabilistic"
	notSupported := "notsupported"

	SetPlannerStrategy(probabilistic)
	assert.Equal(t, probabilistic, currentStrategy.Name(), "(probabilistic) Current strategy should be probabilistic")

	SetPlannerStrategy(notSupported)
	assert.Equal(t, "dummy", currentStrategy.Name(), "(notsupported) Current strategy should be dummy")
}

// func TestSavePlan(t *testing.T) {
// 	defer storage.DeleteAllData(enum.POLICIES)
// 	plc := policy.CreateRandomMockPolicies(1)[0]

// 	err := savePolicy(&plc)
// 	assert.NoError(t, err)
// }

// func TestConvertPolicyToData(t *testing.T) {
// 	plc := policy.CreateRandomMockPolicies(1)[0]
// 	_, err := convertPolicyToData(&plc)
// 	assert.NoError(t, err)
// }

// func TestConvertDataToPolicy(t *testing.T) {
// 	plc := policy.CreateRandomMockPolicies(1)[0]
// 	data_ok, err := convertPolicyToData(&plc)
// 	data_bad := []byte{}

// 	plc_data, err := convertDataToPolicy(data_ok)
// 	assert.NoError(t, err)
// 	assert.Equal(t, plc, plc_data)

// 	_, err = convertDataToPolicy(data_bad)
// 	assert.Error(t, err)
// }

// func TestGetPlannerData(t *testing.T) {
// 	defer storage.DeleteAllData(enum.POLICIES)
// 	var err error

// 	_, err = GetPlannerData()
// 	assert.Error(t, err)

// 	plc := policy.CreateRandomMockPolicies(1)[0]
// 	StoreMockPolicy(plc)
// 	_, err = GetPlannerData()
// 	assert.NoError(t, err)
// }
