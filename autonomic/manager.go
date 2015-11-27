package autonomic

import (
	"time"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/autonomic/analyzer"
	"github.com/elleFlorio/gru/autonomic/executor"
	"github.com/elleFlorio/gru/autonomic/monitor"
	"github.com/elleFlorio/gru/autonomic/planner"
	"github.com/elleFlorio/gru/communication"
)

var manager AutonomicConfig

func Initialize(loopTimeInterval int, nFriends int, dataType string) {
	manager.LoopTimeInterval = loopTimeInterval
	manager.MaxFrineds = nFriends
	manager.DataToShare = dataType
}

func RunLoop() {
	log.WithField("status", "start").Infoln("Running autonomic loop")
	defer log.WithField("status", "done").Infoln("Running autonomic loop")

	c_err := make(chan error)
	c_stop := make(chan struct{})

	go monitor.Start(c_err, c_stop)
	planner.SetPlannerStrategy("probabilistic")

	// Set the ticker for the periodic execution
	ticker := time.NewTicker(time.Duration(manager.LoopTimeInterval) * time.Second)

	for {
		select {
		case <-ticker.C:

			communication.KeepAlive(manager.LoopTimeInterval)
			err := communication.UpdateFriendsData(manager.MaxFrineds)
			if err != nil {
				log.WithField("warning", err).Warnln("Cannot update friends data")
			}

			stats := monitor.Run()
			analytics := analyzer.Run(stats)
			plan := planner.Run(analytics)
			executor.Run(plan)

		case <-c_err:
			log.WithField("status", "error").Errorln("Running autonomic loop")
		case <-c_stop:
			ticker.Stop()
		}
	}
}
