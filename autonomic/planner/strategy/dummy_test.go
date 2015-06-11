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
	mockAnalytics := createMockAnalytics()
	mockServices := createMockServices()
	mockServiceCont := mockServices[0]
	mockServiceImg := mockServices[1]
	mockServiceErr := mockServices[2]

	target, err := dummy.choseTarget(mockPlanCont.TargetType, &mockAnalytics, &mockServiceCont)
	assert.Contains(t, mockAnalytics.Service["service1"].Instances, target, "first plan should have as target an instance of service1")

	target, err = dummy.choseTarget(mockPlanImg.TargetType, &mockAnalytics, &mockServiceImg)
	assert.Equal(t, mockServiceImg.Image, target, "secon plan should have as target the image of service2")

	target, err = dummy.choseTarget(mockPlanErr.TargetType, &mockAnalytics, &mockServiceErr)
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

func createMockAnalytics() analyzer.GruAnalytics {
	instances1 := []string{"instance1", "instance2"}
	service1 := analyzer.ServiceAnalytics{
		CpuAvg:    0.2,
		Instances: instances1}
	instances2 := []string{"instance3"}
	service2 := analyzer.ServiceAnalytics{
		CpuAvg:    0.4,
		Instances: instances2,
	}
	services := map[string]analyzer.ServiceAnalytics{
		"service1": service1,
		"service2": service2,
	}

	instAnalytics1 := analyzer.InstanceAnalytics{
		Cpu:     10000,
		CpuPerc: 0.1,
	}
	instAnalytics2 := analyzer.InstanceAnalytics{
		Cpu:     20000,
		CpuPerc: 0.2,
	}
	instAnalytics3 := analyzer.InstanceAnalytics{
		Cpu:     30000,
		CpuPerc: 0.3,
	}
	instances := map[string]analyzer.InstanceAnalytics{
		"instance1": instAnalytics1,
		"instance2": instAnalytics2,
		"instance3": instAnalytics3,
	}

	system := analyzer.SystemAnalytics{50000}

	mockAnalytics := analyzer.GruAnalytics{
		Service:  services,
		Instance: instances,
		System:   system,
	}

	return mockAnalytics
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
