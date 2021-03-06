package autonomic

import (
	"time"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/autonomic/analyzer"
	"github.com/elleFlorio/gru/autonomic/executor"
	"github.com/elleFlorio/gru/autonomic/monitor"
	"github.com/elleFlorio/gru/autonomic/planner"
	chn "github.com/elleFlorio/gru/channels"
	"github.com/elleFlorio/gru/metric"
)

var (
	ch_err  chan error
	ch_stop chan struct{}
)

func init() {
	ch_err = chn.GetAutonomicErrChannel()
}

func UpdatePlannerStrategy(plannerStrategy string) {
	planner.SetPlannerStrategy(plannerStrategy)
}

func Start() {
	monitor.StartMonitor()
	executor.ListenToActionMessages()
}

func RunLoop(loopTimeInterval int) {
	// Start the metric collector
	metric.StartMetricCollector()
	// Set the ticker for the periodic execution
	ticker := time.NewTicker(time.Duration(loopTimeInterval) * time.Second)
	log.Infoln("Running autonomic loop")
	for {
		select {
		case <-ticker.C:
			stats := monitor.Run()
			analytics := analyzer.Run(stats)
			policy := planner.Run(analytics)
			executor.Run(policy)

			log.Infoln("-------------------------")

		case <-ch_err:
			log.Debugln("Error running autonomic loop")
		}
	}
}
