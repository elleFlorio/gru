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
			err := communication.UpdateFriendsData(manager.MaxFrineds, manager.DataToShare)
			if err != nil {
				log.WithField("waring", err).Warnln("Running autonomic loop")
			}

			stats := monitor.Run()

			log.WithFields(log.Fields{
				"status":    "received stats",
				"instances": len(stats.Instance),
				"services":  len(stats.Service),
			}).Debugln("Running autonomic loop")

			analytics := analyzer.Run(stats)

			log.WithFields(log.Fields{
				"status":    "received analytics",
				"instances": len(analytics.Instance),
				"services":  len(analytics.Service),
			}).Debugln("Running autonomic loop")

			plan := planner.Run(analytics)

			log.WithFields(log.Fields{
				"status":     "received plan",
				"Service":    plan.Service,
				"TargetType": plan.TargetType,
				"Target":     plan.Target,
				"Weight":     plan.Weight,
				"Actions":    plan.Actions,
			}).Debugln("Running autonomic loop")

			executor.Run(plan)
		case <-c_err:
			log.WithField("status", "error").Errorln("Running autonomic loop")
		case <-c_stop:
			ticker.Stop()
		}
	}
}
