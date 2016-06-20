package service

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/enum"
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

	srv1, err := GetServiceById("instance1_0")
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

func TestFindIdIndex(t *testing.T) {
	instances := []string{
		"instance1_1",
		"instance1_2",
		"instance1_3",
		"instance1_4",
		"instance2_1",
	}

	index, _ := findIdIndex("instance1_3", instances)
	assert.Equal(t, 2, index, "index of 'instance3' should be 2")
}

func TestAddServiceInstance(t *testing.T) {
	defer cfg.CleanServices()
	cfg.SetServices(CreateMockServices())
	name := "service1"
	service, _ := GetServiceByName(name)
	instance_p := "instance1_100"
	instance_ps := "instance1_101"
	instance_s := "instance1_102"
	var err error
	var status enum.Status

	status = enum.PENDING
	err = AddServiceInstance(name, instance_p, status)
	assert.NoError(t, err)
	assert.Contains(t, service.Instances.All, instance_p)
	assert.Contains(t, service.Instances.Pending, instance_p)

	status = enum.PAUSED
	err = AddServiceInstance(name, instance_ps, status)
	assert.NoError(t, err)
	assert.Contains(t, service.Instances.All, instance_ps)
	assert.Contains(t, service.Instances.Paused, instance_ps)

	status = enum.STOPPED
	err = AddServiceInstance(name, instance_s, status)
	assert.NoError(t, err)
	assert.Contains(t, service.Instances.All, instance_s)
	assert.Contains(t, service.Instances.Stopped, instance_s)

	status = enum.PENDING
	err = AddServiceInstance("pippo", "pippo", status)
	assert.Error(t, err)
}

func TestChangeServiceInstanceStatus(t *testing.T) {
	defer cfg.CleanServices()
	cfg.SetServices(CreateMockServices())
	name := "service1"
	service, _ := GetServiceByName(name)
	instance := "instance1_3"
	var prev enum.Status
	var upd enum.Status
	var err error

	prev = enum.PENDING
	upd = enum.RUNNING
	err = ChangeServiceInstanceStatus(name, instance, prev, upd)
	assert.NotContains(t, service.Instances.Pending, instance)
	assert.Contains(t, service.Instances.Running, instance)
	assert.NoError(t, err)

	prev = enum.RUNNING
	upd = enum.PAUSED
	err = ChangeServiceInstanceStatus(name, instance, prev, upd)
	assert.NotContains(t, service.Instances.Running, instance)
	assert.Contains(t, service.Instances.Paused, instance)
	assert.NoError(t, err)

	prev = enum.PAUSED
	upd = enum.STOPPED
	err = ChangeServiceInstanceStatus(name, instance, prev, upd)
	assert.NotContains(t, service.Instances.Paused, instance)
	assert.Contains(t, service.Instances.Stopped, instance)
	assert.NoError(t, err)

	prev = enum.STOPPED
	upd = enum.PENDING
	err = ChangeServiceInstanceStatus(name, instance, prev, upd)
	assert.NotContains(t, service.Instances.Stopped, instance)
	assert.Contains(t, service.Instances.Pending, instance)
	assert.NoError(t, err)

	prev = enum.PENDING
	upd = enum.RUNNING
	err = ChangeServiceInstanceStatus("pippo", instance, prev, upd)
	assert.Error(t, err)

	prev = enum.PENDING
	upd = enum.RUNNING
	err = ChangeServiceInstanceStatus(name, "pippo", prev, upd)
	assert.Error(t, err)
}

func TestRemoveServiceIntance(t *testing.T) {
	defer cfg.CleanServices()
	cfg.SetServices(CreateMockServices())
	name := "service1"
	service, _ := GetServiceByName(name)
	instance := "instance1_0"
	var err error

	err = RemoveServiceInstance(name, instance)
	assert.NotContains(t, service.Instances.Stopped, instance)
	assert.NotContains(t, service.Instances.All, instance)
	assert.NoError(t, err)

	err = RemoveServiceInstance("pippo", instance)
	assert.Error(t, err)

	err = RemoveServiceInstance(name, "pippo")
	assert.Error(t, err)

}

func TestGetServiceInstanceStatus(t *testing.T) {
	defer cfg.CleanServices()
	cfg.SetServices(CreateMockServices())
	name := "service1"
	instance_p := "instance1_3"
	instance_r := "instance1_1"
	instance_s := "instance1_0"
	instance_ps := "instance1_5"
	var status enum.Status

	status = GetServiceInstanceStatus(name, instance_p)
	assert.Equal(t, enum.PENDING, status)

	status = GetServiceInstanceStatus(name, instance_r)
	assert.Equal(t, enum.RUNNING, status)

	status = GetServiceInstanceStatus(name, instance_s)
	assert.Equal(t, enum.STOPPED, status)

	status = GetServiceInstanceStatus(name, instance_ps)
	assert.Equal(t, enum.PAUSED, status)

	status = GetServiceInstanceStatus(name, "pippo")
	assert.Equal(t, enum.UNKNOWN, status)

	status = GetServiceInstanceStatus("pippo", "pippo")
	assert.Equal(t, enum.UNKNOWN, status)
}
