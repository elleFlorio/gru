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

func TestRetrieveAnalytics(t *testing.T) {
	defer storage.DeleteAllData(enum.ANALYTICS)
	var err error

	_, err = retrieveAnalytics()
	assert.Error(t, err)

	analyzer.StoreMockAnalytics()
	_, err = retrieveAnalytics()
	assert.NoError(t, err)
}

func TestBuildPlans(t *testing.T) {
	srvcs := service.CreateMockServices()
	service.UpdateServices(srvcs)
	analytics := analyzer.CreateMockAnalytics()
	plans := buildPlans(analytics)
	nPlans := len(srvcs) * len(policy.List())

	assert.Len(t, plans, nPlans)
}

func TestSavePlan(t *testing.T) {
	defer storage.DeleteAllData(enum.PLANS)
	plan := strategy.CreateRandomPlans(1)[0]

	err := savePlan(&plan)
	assert.NoError(t, err)
}

func TestRun(t *testing.T) {
	srvcs := service.CreateMockServices()
	service.UpdateServices(srvcs)
	assert.NotPanics(t, Run)

	analyzer.StoreMockAnalytics()
	assert.NotPanics(t, Run)

}
