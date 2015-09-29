package strategy

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	"github.com/elleFlorio/gru/autonomic/analyzer"
	"github.com/elleFlorio/gru/service"
)

var dummy DummyStrategy

func init() {
	dummy = DummyStrategy{}
}

func TestChosePlanDummy(t *testing.T) {
	mockPlans := CreateMockPlans(0.8, 0.5, 0.3)
	theMockPlan := dummy.chosePlan(mockPlans)

	assert.Equal(t, "service1", theMockPlan.Service, "Chosen service should be pippo")
}

func TestChoseTargetDummy(t *testing.T) {
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
