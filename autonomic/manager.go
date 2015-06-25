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
	log.WithField("status", "start").Infoln("Running autonomic loop")

	man.loop()

	log.WithField("status", "done").Infoln("Running autonomic loop")
}

func (man *autoManager) loop() {
	c_err := make(chan error)
	c_stop := make(chan struct{})

	m := monitor.NewMonitor(c_stop, c_err)
	a := analyzer.NewAnalyzer(c_err)
	p := planner.NewPlanner("probabilistic", c_err)
	e := executor.NewExecutor(c_err)

	go m.Start(man.Docker)

	// Set the ticker for the periodic execution
	ticker := time.NewTicker(time.Duration(man.LoopTimeInterval) * time.Second)

	for {
		select {
		case <-ticker.C:
			stats := m.Run()

			log.WithFields(log.Fields{
				"status":    "received stats",
				"instances": len(stats.Instance),
				"services":  len(stats.Service),
			}).Debugln("Running autonomic loop")

			analytics := a.Run(stats)

			log.WithFields(log.Fields{
				"status":    "received analytics",
				"instances": len(analytics.Instance),
				"services":  len(analytics.Service),
			}).Debugln("Running autonomic loop")

			plan := p.Run(analytics)

			log.WithFields(log.Fields{
				"status":     "received plan",
				"Service":    plan.Service,
				"TargetType": plan.TargetType,
				"Target":     plan.Target,
				"Weight":     plan.Weight,
				"Actions":    plan.Actions,
			}).Debugln("Running autonomic loop")

			e.Run(plan, man.Docker)
		case <-c_err:
			log.WithField("status", "error").Errorln("Running autonomic loop")
		case <-c_stop:
			ticker.Stop()
		}
	}
}
