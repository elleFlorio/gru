package autonomic

import (
	"time"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/autonomic/analyzer"
	"github.com/elleFlorio/gru/autonomic/executor"
	"github.com/elleFlorio/gru/autonomic/monitor"
	"github.com/elleFlorio/gru/autonomic/planner"
	"github.com/elleFlorio/gru/communication"
	"github.com/elleFlorio/gru/metric"
)

var manager AutonomicConfig

func Initialize(loopTimeInterval int, nFriends int, dataType string) {
	manager.LoopTimeInterval = loopTimeInterval
	manager.MaxFrineds = nFriends
	manager.DataToShare = dataType
}

func RunLoop() {
	c_err := make(chan error)
	c_stop := make(chan struct{})

	monitor.Start(c_err, c_stop)
	planner.SetPlannerStrategy("probabilistic")

	// Set the ticker for the periodic execution
	ticker := time.NewTicker(time.Duration(manager.LoopTimeInterval) * time.Second)

	log.Infoln("Running autonomic loop")
	for {
		select {
		case <-ticker.C:
			communication.KeepAlive(manager.LoopTimeInterval)
			err := communication.UpdateFriendsData(manager.MaxFrineds)
			if err != nil {
				log.WithField("err", err).Debugln("Cannot update friends data")
			}

			stats := monitor.Run()
			analytics := analyzer.Run(stats)
			plan := planner.Run(analytics)
			executor.Run(plan)

			collectMetrics()

			log.Infoln("-------------------------")

		case <-c_err:
			log.Errorln("Error running autonomic loop")
		case <-c_stop:
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
