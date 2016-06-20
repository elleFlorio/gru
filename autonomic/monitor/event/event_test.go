package event

import (
	"testing"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	cfg "github.com/elleFlorio/gru/configuration"
	dsc "github.com/elleFlorio/gru/discovery"
	"github.com/elleFlorio/gru/enum"
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
	id2_ps := "instance2_ps"
	status2_s := enum.STOPPED
	status2_p := enum.PENDING
	status2_ps := enum.PAUSED
	service, _ := srv.GetServiceByName(esrv)

	e = createEvent(etype, esrv, eimg, id2_s, status2_s)
	HanldeStartEvent(e)
	assert.Contains(t, service.Instances.All, id2_s)
	assert.Contains(t, service.Instances.Stopped, id2_s)

	// check add pending
	e = createEvent(etype, esrv, eimg, id2_p, status2_p)
	HanldeStartEvent(e)
	assert.Contains(t, service.Instances.All, id2_p)
	assert.Contains(t, service.Instances.Pending, id2_p)
	assert.Contains(t, events.Service[esrv].Start, id2_p)

	//check add paused
	e = createEvent(etype, esrv, eimg, id2_ps, status2_ps)
	HanldeStartEvent(e)
	assert.Contains(t, service.Instances.Paused, id2_ps)
}

func TestHandlePromoteEvent(t *testing.T) {
	defer resetMockServices()
	var e Event

	etype := "start"
	esrv := "service1"
	eimg := "img"
	einst := "instance1_3"
	estat := enum.PENDING
	service, _ := srv.GetServiceByName(esrv)

	e = createEvent(etype, esrv, eimg, einst, estat)
	HandlePromoteEvent(e)
	assert.Contains(t, service.Instances.Running, einst)
	assert.NotContains(t, service.Instances.Pending, einst)
}

func TestHandleStopEvent(t *testing.T) {
	defer resetMockServices()
	var e Event

	etype := "stop"
	esrv := "service2"
	eimg := "img"
	mockInstId_r := "instance2_1"
	mockInstId_p := "instance1_3"
	estat := enum.RUNNING
	service, _ := srv.GetServiceByName(esrv)

	// check error
	e = createEvent(etype, "pippo", eimg, mockInstId_p, estat)
	HandleStopEvent(e)

	// check running
	e = createEvent(etype, esrv, eimg, mockInstId_r, estat)
	HandleStopEvent(e)
	assert.NotContains(t, service.Instances.Running, mockInstId_r)
	assert.Contains(t, events.Service[esrv].Stop, mockInstId_r)

	// check pending
	e = createEvent(etype, esrv, eimg, mockInstId_p, estat)
	HandleStopEvent(e)

	assert.NotContains(t, service.Instances.Pending, mockInstId_p)
	assert.Contains(t, events.Service[esrv].Stop, mockInstId_r)
}

func TestHandleRemoveEvent(t *testing.T) {
	defer resetMockServices()
	defer log.SetLevel(log.ErrorLevel)
	var e Event
	var service *cfg.Service

	etype := "destroy"
	eimg := "img"
	service1 := "service1"
	mockInstId_s := "instance1_0"
	estat_s := enum.STOPPED
	service2 := "service2"
	mockInstId_r := "instance2_1"
	estat_r := enum.RUNNING
	mockInstId_wrong := "pippo"

	e = createEvent(etype, service1, eimg, mockInstId_s, estat_s)
	HandleRemoveEvent(e)
	service, _ = srv.GetServiceByName(service1)
	assert.NotContains(t, service.Instances.Running, mockInstId_s)
	assert.NotContains(t, service.Instances.Pending, mockInstId_s)
	assert.NotContains(t, service.Instances.Stopped, mockInstId_s)

	e = createEvent(etype, service2, eimg, mockInstId_r, estat_r)
	HandleRemoveEvent(e)
	service, _ = srv.GetServiceByName(service2)
	assert.NotContains(t, service.Instances.Running, mockInstId_r)
	assert.NotContains(t, service.Instances.Pending, mockInstId_r)
	assert.NotContains(t, service.Instances.Stopped, mockInstId_r)

	// Check the log for this test
	log.SetLevel(log.DebugLevel)
	e = createEvent(etype, "pippo", eimg, mockInstId_wrong, estat_s)
	HandleRemoveEvent(e)
}

func createEvent(etype string, esrv string, eimg string, einst string, estat enum.Status) Event {
	return Event{
		Type:     etype,
		Service:  esrv,
		Image:    eimg,
		Instance: einst,
		Status:   estat,
	}
}

func resetMockServices() {
	mockServices := srv.CreateMockServices()
	cfg.SetServices(mockServices)
}
