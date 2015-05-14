package autonomic

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/samalba/dockerclient"

	"github.com/elleFlorio/gru/autonomic/analyzer"
	"github.com/elleFlorio/gru/autonomic/executor"
	"github.com/elleFlorio/gru/autonomic/monitor"
	"github.com/elleFlorio/gru/autonomic/planner"
)

type autoManager struct {
	Docker           *dockerclient.DockerClient
	LoopTimeInterval int
}

var manager autoManager

func NewAutoManager(docker *dockerclient.DockerClient, loopTimeInternval int) *autoManager {
	manager = autoManager{
		docker,
		loopTimeInternval,
	}
	return &manager
}

func (man *autoManager) RunLoop() {
	log.WithField("status", "start").Debugln("Running autonomic loop")

	man.loop()

	log.WithField("status", "done").Debugln("Running autonomic loop")
}

func (man *autoManager) loop() {
	c_err := make(chan error)
	c_stop := make(chan struct{})

	m := monitor.NewMonitor(c_stop, c_err)
	a := analyzer.NewAnalyzer(c_err)
	p := planner.NewPlanner()
	e := executor.NewExecutor()

	m.Start(man.Docker)

	// Set the ticker for the periodic execution
	ticker := time.NewTicker(time.Duration(man.LoopTimeInterval) * time.Second)

	for {
		select {
		case <-ticker.C:
			stats := m.Run()
			analytics := a.Run(stats)
			p.Run(analytics)
			e.Run()
		case <-c_err:
			log.WithField("status", "error").Debugln("Running autonomic loop")
		case <-c_stop:
			ticker.Stop()
		}
	}
}
