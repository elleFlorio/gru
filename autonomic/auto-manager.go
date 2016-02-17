package autonomic

import (
	"time"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/autonomic/analyzer"
	"github.com/elleFlorio/gru/autonomic/executor"
	"github.com/elleFlorio/gru/autonomic/monitor"
	"github.com/elleFlorio/gru/autonomic/planner"
	"github.com/elleFlorio/gru/metric"
)

var (
	ch_err  chan error
	ch_stop chan struct{}
)

func init() {
	ch_err = make(chan error)
	ch_stop = make(chan struct{})
}

func Initialize(plannerStrategy string) {
	planner.SetPlannerStrategy("probabilistic")
}

func Start() {
	monitor.Start(ch_err, ch_stop)
	executor.ListenToActionMessages()
}

func RunLoop(loopTimeInterval int) {
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

			collectMetrics()

			log.Infoln("-------------------------")

		case <-ch_err:
			log.Errorln("Error running autonomic loop")
		case <-ch_stop:
			ticker.Stop()
		}
	}
}

func collectMetrics() {
	log.Debugln("Collecting metrics")
	metric.UpdateMetrics()
	err := metric.StoreMetrics(metric.Metrics())
	if err != nil {
		log.WithField("errr", err).Errorln("Error collecting agent metrics")
	}
}
