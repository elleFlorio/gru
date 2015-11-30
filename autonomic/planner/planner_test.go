package planner

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	"github.com/elleFlorio/gru/autonomic/analyzer"
	"github.com/elleFlorio/gru/autonomic/planner/policy"
	"github.com/elleFlorio/gru/autonomic/planner/strategy"
	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/service"
	"github.com/elleFlorio/gru/storage"
)

func init() {
	storage.New("internal")
}

func TestSetPlannerStrategy(t *testing.T) {
	probabilistic := "probabilistic"
	notSupported := "notsupported"

	SetPlannerStrategy(probabilistic)
	assert.Equal(t, probabilistic, currentStrategy.Name(), "(probabilistic) Current strategy should be probabilistic")

	SetPlannerStrategy(notSupported)
	assert.Equal(t, "dummy", currentStrategy.Name(), "(notsupported) Current strategy should be dummy")
}

func TestBuildPlans(t *testing.T) {
	srvcs := service.CreateMockServices()
	service.UpdateServices(srvcs)
	analytics := analyzer.GruAnalytics{}
	plans := buildPlans(analytics)
	assert.Len(t, plans, 1)
	analytics = analyzer.CreateMockAnalytics()
	plans = buildPlans(analytics)
	nPlans := len(srvcs)*len(policy.List()) + 1 // count also the noAction plan
	assert.Len(t, plans, nPlans)
}

func TestSavePlan(t *testing.T) {
	defer storage.DeleteAllData(enum.PLANS)
	plan := strategy.CreateRandomPlans(1)[0]

	err := savePlan(&plan)
	assert.NoError(t, err)
}

func TestConvertPlanToData(t *testing.T) {
	plan := strategy.CreateMockPlan(0.0, service.Service{}, enum.Actions{enum.START})

	_, err := convertPlanToData(plan)
	assert.NoError(t, err)
}

func TestConvertDataToPlan(t *testing.T) {
	plan := strategy.CreateMockPlan(0.0, service.Service{}, enum.Actions{enum.START})
	data_ok, err := convertPlanToData(plan)
	data_bad := []byte{}

	_, err = convertDataToPlan(data_ok)
	assert.NoError(t, err)

	_, err = convertDataToPlan(data_bad)
	assert.Error(t, err)
}

func TestGetPlannerData(t *testing.T) {
	defer storage.DeleteAllData(enum.PLANS)
	var err error

	_, err = GetPlannerData()
	assert.Error(t, err)

	strategy.StoreMockPlan(0.8, service.Service{}, enum.Actions{})
	_, err = GetPlannerData()
	assert.NoError(t, err)
}

func TestRun(t *testing.T) {
	srvcs := service.CreateMockServices()
	service.UpdateServices(srvcs)
	assert.Nil(t, Run(analyzer.GruAnalytics{}))

	analytics := analyzer.CreateMockAnalytics()
	assert.NotNil(t, Run(analytics))
}
