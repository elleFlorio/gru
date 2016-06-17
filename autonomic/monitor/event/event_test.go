package event

import (
	"testing"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	cfg "github.com/elleFlorio/gru/configuration"
	dsc "github.com/elleFlorio/gru/discovery"
	srv "github.com/elleFlorio/gru/service"
)

func init() {
	dsc.New("noservice", "")
	cfg.GetAgentDiscovery().TTL = 5
	resetMockServices()
	Initialize(srv.List())
}

func TestHandleStartEvent(t *testing.T) {
	defer resetMockServices()
	var e Event

	etype := "start"
	esrv := "service2"
	eimg := "img"
	id2_s := "instance2_s"
	id2_p := "instance2_p"
	id2_r := "instance2_r"
	status2_s := "stopped"
	status2_p := "pending"
	status2_r := "running"
	service, _ := srv.GetServiceByName(esrv)

	e = createEvent(etype, esrv, eimg, id2_s, status2_s)
	HanldeStartEvent(e)
	assert.Contains(t, service.Instances.All, id2_s,
		"(new -> stopped) Service 2 - instances - all, should contain added instance")
	assert.Contains(t, service.Instances.Stopped, id2_s,
		"(new -> stopped) Service 2 - instances - stopped, should contain added instance")

	// check add pending
	e = createEvent(etype, esrv, eimg, id2_p, status2_p)
	HanldeStartEvent(e)
	assert.Contains(t, service.Instances.All, id2_p,
		"(new -> pending) Service 2 - instances - all, should contain added instance")
	assert.Contains(t, service.Instances.Pending, id2_p,
		"(new -> pending) Service 2 - instances - pending, should contain added instance")
	assert.Contains(t, events.Service[esrv].Start, id2_p,
		"(new -> pending) Service 2 - events - start, should contain added instance")

	// check add running
	e = createEvent(etype, esrv, eimg, id2_r, status2_r)
	HanldeStartEvent(e)
	assert.Contains(t, service.Instances.All, id2_r,
		"(new -> running) Service 2 - instances - all, should contain added instance")
	assert.Contains(t, service.Instances.Running, id2_r,
		"(new -> running) Service 2 - instances - running, should contain added instance")

	//check stopped -> pending
	e = createEvent(etype, esrv, eimg, id2_s, status2_p)
	HanldeStartEvent(e)
	assert.Contains(t, service.Instances.Pending, id2_s,
		"(stopped -> pending) Service 2 - instances - pending, should contain added instance")
	assert.Contains(t, events.Service[esrv].Start, id2_s,
		"(stopped -> pending) Service 2 - events - start, should contain added instance")
	assert.NotContains(t, service.Instances.Stopped, id2_s,
		"(stopped -> pending) Service 2 - instances - stopped, should not contain added instance")

	//check pending -> running
	e = createEvent(etype, esrv, eimg, id2_s, status2_r)
	HanldeStartEvent(e)
	assert.Contains(t, service.Instances.Running, id2_s,
		"(pending -> running) Service 2 - instances - running, should contain added instance")
	assert.NotContains(t, service.Instances.Pending, id2_s,
		"(pending -> running) Service 2 - instances - pending, should not contain added instance")
}

func TestHandleStopEvent(t *testing.T) {
	defer resetMockServices()
	var e Event

	etype := "stop"
	esrv := "service2"
	eimg := "img"
	mockInstId_r := "instance2_1"
	mockInstId_p := "instance1_3"
	estat := "status"
	service, _ := srv.GetServiceByName(esrv)

	// check error
	e = createEvent(etype, "pippo", eimg, mockInstId_p, estat)
	HandleStopEvent(e)

	// check running
	e = createEvent(etype, esrv, eimg, mockInstId_r, estat)
	HandleStopEvent(e)
	assert.NotContains(t, service.Instances.Running, mockInstId_r,
		"(running) Service stats should not contain 'instance2_1'")
	assert.Contains(t, events.Service[esrv].Stop, mockInstId_r,
		"(running) Events Stop should contain 'instance2_1'")

	// check pending
	e = createEvent(etype, esrv, eimg, mockInstId_p, estat)
	HandleStopEvent(e)

	assert.NotContains(t, service.Instances.Pending, mockInstId_p,
		"(pending) Service stats should not contain 'instance1_3'")
	assert.Contains(t, events.Service[esrv].Stop, mockInstId_r,
		"(running) Events Stop should contain 'instance1_3'")
}

func TestHandleRemoveEvent(t *testing.T) {
	defer resetMockServices()
	defer log.SetLevel(log.ErrorLevel)
	var e Event
	var service *cfg.Service

	etype := "destroy"
	eimg := "img"
	estat := "status"
	service1 := "service1"
	mockInstId_s := "instance1_0"
	service2 := "service2"
	mockInstId_r := "instance2_1"
	mockInstId_wrong := "pippo"

	e = createEvent(etype, service1, eimg, mockInstId_s, estat)
	HandleRemoveEvent(e)
	service, _ = srv.GetServiceByName(service1)
	assert.NotContains(t, service.Instances.Running, mockInstId_s)
	assert.NotContains(t, service.Instances.Pending, mockInstId_s)
	assert.NotContains(t, service.Instances.Stopped, mockInstId_s)

	e = createEvent(etype, service2, eimg, mockInstId_r, estat)
	HandleRemoveEvent(e)
	service, _ = srv.GetServiceByName(service2)
	assert.NotContains(t, service.Instances.Running, mockInstId_r)
	assert.NotContains(t, service.Instances.Pending, mockInstId_r)
	assert.NotContains(t, service.Instances.Stopped, mockInstId_r)

	// Check the log for this test
	log.SetLevel(log.DebugLevel)
	e = createEvent(etype, "pippo", eimg, mockInstId_wrong, estat)
	HandleRemoveEvent(e)
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

func createEvent(etype string, esrv string, eimg string, einst string, estat string) Event {
	return Event{
		Type:     etype,
		Service:  esrv,
		Image:    eimg,
		Isntance: einst,
		Status:   estat,
	}
}

func resetMockServices() {
	mockServices := srv.CreateMockServices()
	cfg.SetServices(mockServices)
}
