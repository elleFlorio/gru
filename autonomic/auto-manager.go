package autonomic

import (
	"time"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/autonomic/analyzer"
	"github.com/elleFlorio/gru/autonomic/executor"
	"github.com/elleFlorio/gru/autonomic/monitor"
	"github.com/elleFlorio/gru/autonomic/planner"
	chn "github.com/elleFlorio/gru/channels"
	com "github.com/elleFlorio/gru/communication"
	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/metric"
)

var (
	ch_err  chan error
	ch_stop chan struct{}

	interval int
)

func init() {
	ch_err = chn.GetAutonomicErrChannel()
	interval = cfg.GetAgentAutonomic().LoopTimeInterval
}

func UpdatePlannerStrategy(plannerStrategy string) {
	planner.SetPlannerStrategy(plannerStrategy)
}

func Start() {
	monitor.StartMonitor()
	executor.ListenToActionMessages()
}

func RunLoop() {
	// Start the metric collector
	metric.StartMetricCollector()
	// Set the ticker for the periodic execution
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	log.Infoln("Running autonomic loop")
	for {
		select {
		case <-ticker.C:
			stats := monitor.Run()
			analytics := analyzer.Run(stats)
			policy := planner.Run(analytics)
			executor.Run(policy)
			com.UpdateFriends()
			if checkNewInterval() {
				ticker = time.NewTicker(time.Duration(interval) * time.Second)
				log.WithField("interval", interval).Debugln("Updated autonomic loop time interval")
			}

			log.Infoln("-------------------------")

		case <-ch_err:
			log.Debugln("Error running autonomic loop")
		}
	}
}

func checkNewInterval() bool {
	if interval != cfg.GetAgentAutonomic().LoopTimeInterval {
		interval = cfg.GetAgentAutonomic().LoopTimeInterval
		return true
	}

	return false
}
