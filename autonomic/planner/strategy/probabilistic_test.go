package strategy

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	"github.com/elleFlorio/gru/autonomic/analyzer"
	"github.com/elleFlorio/gru/service"
)

var prob ProbabilisticStrategy

func init() {
	prob = ProbabilisticStrategy{}
}

func TestRandomUniform(t *testing.T) {
	val := prob.randUniform(0, 1)

	assert.InDelta(t, 0.5, val, 0.5, "Expected value should be in (0,1))")
}

func TestWeightedRandomElement(t *testing.T) {
	mockPlans := CreateMockPlans(0.8, 0.5, 0.3)
	_, err := prob.weightedRandomElement(mockPlans)
	assert.NoError(t, err, "No error should be returned")

	mockPlans = CreateMockPlans(0.0, 0.0, 0.0)
	_, err = prob.weightedRandomElement(mockPlans)
	assert.Error(t, err, "Error is expected for total weight equals 0")

}

func TestChosePlanProb(t *testing.T) {
	mockPlans := CreateMockPlans(0.8, 0.5, 0.3)
	mockPlan := prob.chosePlan(mockPlans)
	assert.NotEqual(t, "none", mockPlan.Service, "Chosen plan - service should not be 'none'")

	mockPlans = CreateMockPlans(0.0, 0.0, 0.0)
	mockPlan = prob.chosePlan(mockPlans)
	assert.Equal(t, "none", mockPlan.Service, "Chosen plan - service should be 'none'")
}

func TestChoseTargetProb(t *testing.T) {
	mockPlans := CreateMockPlans(0.8, 0.5, 0.3)
	mockPlanCont := mockPlans[0]
	mockPlanImg := mockPlans[1]
	mockPlanErr := mockPlans[2]
	mockAnalytics := analyzer.CreateMockAnalytics()
	mockServices := service.CreateMockServices()
	mockServiceCont := mockServices[0]
	mockServiceImg := mockServices[1]
	mockServiceErr := mockServices[2]

	target, err := dummy.choseTarget(mockPlanCont.TargetType, "running", mockAnalytics, &mockServiceCont)
	assert.Contains(t, mockAnalytics.Service["service1"].Instances.Active, target, "first plan should have as target an instance of service1")

	target, err = dummy.choseTarget(mockPlanImg.TargetType, "running", mockAnalytics, &mockServiceImg)
	assert.Equal(t, mockServiceImg.Image, target, "secon plan should have as target the image of service2")

	target, err = dummy.choseTarget(mockPlanErr.TargetType, "running", mockAnalytics, &mockServiceErr)
	assert.Error(t, err, "third plan should produce an error in the choice of the target")

}
