package strategy

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/service"
)

func TestConvertPlanToData(t *testing.T) {
	plan := CreateMockPlan(enum.WHITE, service.Service{}, []enum.Action{enum.START})

	_, err := ConvertPlanToData(plan)
	assert.NoError(t, err)
}

func TestConvertDataToPlan(t *testing.T) {
	plan := CreateMockPlan(enum.WHITE, service.Service{}, []enum.Action{enum.START})
	data_ok, err := ConvertPlanToData(plan)
	data_bad := []byte{}

	_, err = ConvertDataToPlan(data_ok)
	assert.NoError(t, err)

	_, err = ConvertDataToPlan(data_bad)
	assert.Error(t, err)
}
