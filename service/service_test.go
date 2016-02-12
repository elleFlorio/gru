package service

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	cfg "github.com/elleFlorio/gru/configuration"
)

func TestGetServiceByType(t *testing.T) {
	defer cfg.CleanServices()
	cfg.SetServices(CreateMockServices())

	ws := GetServiceByType("webserver")
	if assert.Len(t, ws, 2, "there should be 2 services with type webserver") {
		assert.Equal(t, "webserver", ws[0].Type, "service type should be webserver")
		img := [2]string{ws[0].Image, ws[1].Image}
		assert.Contains(t, img, "test/tomcat", "images should contain test/tomcat")
	}

	db := GetServiceByType("database")
	if assert.Len(t, db, 1, "there should be 1 services with type database") {
		assert.Equal(t, "database", db[0].Type, "service type should be database")
	}

	ap := GetServiceByType("application")
	assert.Len(t, ap, 0, "there should be 0 services with type application")
}

func TestGetServiceByName(t *testing.T) {
	defer cfg.CleanServices()
	cfg.SetServices(CreateMockServices())

	s1, err := GetServiceByName("service2")
	assert.Equal(t, "service2", s1.Name, "service name should be service1")

	_, err = GetServiceByName("pippo")
	assert.Error(t, err, "There should be no service with name 'pippo'")
}

func TestGetServiceByImage(t *testing.T) {
	defer cfg.CleanServices()
	cfg.SetServices(CreateMockServices())

	img1, err := GetServiceByImage("test/mysql")
	assert.Equal(t, "test/mysql", img1.Image, "service image should be test/tomcat")

	_, err = GetServiceByImage("test/pippo")
	assert.Error(t, err, "There should be no image 'test/pippo'")
}

func TestGetServiceById(t *testing.T) {
	defer cfg.CleanServices()

	var err error
	cfg.SetServices(CreateMockServices())

	srv1, err := GetServiceById("instance1_2")
	assert.NoError(t, err)
	assert.Equal(t, "service1", srv1.Name)

	srv2, err := GetServiceById("instance2_1")
	assert.NoError(t, err)
	assert.Equal(t, "service2", srv2.Name)

	_, err = GetServiceById("pippo")
	assert.Error(t, err)
}

func TestAddServices(t *testing.T) {
	defer cfg.CleanServices()
	cfg.SetServices(CreateMockServices())

	newService := cfg.Service{
		Name:  "newService",
		Type:  "mockService",
		Image: "noImage",
	}
	newServices := []cfg.Service{newService}

	cfg.AddServices(newServices)
	assert.Contains(t, cfg.GetServices(), newService, "services should contain the added service")
}

func TestRemoveServices(t *testing.T) {
	defer cfg.CleanServices()
	cfg.SetServices(CreateMockServices())
	rmService := "service2"
	rmServices := []string{rmService}

	cfg.RemoveServices(rmServices)
	assert.NotContains(t, List(), rmService, "services should not contain removed service 'service2'")
}

func TestUpdateServices(t *testing.T) {
	defer cfg.CleanServices()
	cfg.SetServices(CreateMockServices())
	newService := cfg.Service{
		Name:  "newService",
		Type:  "mockService",
		Image: "noImage",
	}
	newServices := []cfg.Service{newService}

	cfg.SetServices(newServices)
	assert.Len(t, cfg.GetServices(), 1, "services should have lenght = 1 after the update")
	assert.Contains(t, List(), "newService", "services should contain service 'newService' after the update")
}
