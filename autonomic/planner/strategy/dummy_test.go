package strategy

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elleFlorio/gru/autonomic/analyzer"
	"github.com/elleFlorio/gru/service"
)

var dummy DummyStrategy

func init() {
	dummy = DummyStrategy{}
}

func TestChosePlan(t *testing.T) {
	mockPlans := createMockPlans()
	theMockPlan := dummy.chosePlan(mockPlans)

	assert.Equal(t, "service1", theMockPlan.Service, "Chosen service should be pippo")
}

func TestChoseTarget(t *testing.T) {
	mockPlans := createMockPlans()
	mockPlanCont := mockPlans[0]
	mockPlanImg := mockPlans[1]
	mockPlanErr := mockPlans[2]
	//mockAnalytics := createMockAnalytics()
	mockAnalytics := analyzer.CreateMockAnalytics()
	mockServices := createMockServices()
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

func createMockPlans() []GruPlan {
	p1 := GruPlan{
		Service:    "service1",
		Weight:     0.8,
		TargetType: "container",
		Actions:    []string{"start, stop"},
	}

	p2 := GruPlan{
		Service:    "service2",
		Weight:     0.5,
		TargetType: "image",
		Actions:    []string{"open"},
	}

	p3 := GruPlan{
		Service:    "service3",
		Weight:     0.3,
		TargetType: "notExist",
		Actions:    []string{"close, shutdown"},
	}

	return []GruPlan{p1, p2, p3}

}

func createMockServices() []service.Service {
	service1 := service.Service{
		Name:  "service1",
		Type:  "webserver",
		Image: "test/tomcat",
	}

	service2 := service.Service{
		Name:  "service2",
		Type:  "webserver",
		Image: "test/jetty",
	}

	service3 := service.Service{
		Name:  "service3",
		Type:  "database",
		Image: "test/mysql",
	}

	return []service.Service{service1, service2, service3}
}
