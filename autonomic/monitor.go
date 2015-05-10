package autonomic

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/samalba/dockerclient"
)

var m_data monitorData

func statCallBack(id string, stats *dockerclient.Stats, ec chan error, args ...interface{}) {
	//In the callback I should update my stats
	m_data.stats[id] = *stats
}

// Events are: create, destroy, die, exec_create, exec_start, export, kill, oom, pause, restart, start, stop, unpause
func eventCallback(event *dockerclient.Event, ec chan error, args ...interface{}) {
	log.WithFields(log.Fields{
		"from":   event.From,
		"status": event.Status,
	}).Debug("Received event")

	client := manager.Client
	switch event.Status {
	case "start":
		client.StartMonitorStats(event.Id, statCallBack, nil, m_data)
		info, err := client.InspectContainer(event.Id)
		if err != nil {
			log.Errorln("Error in inspecting containers: ", err)
			panic(err)
		}
		m_data.info[event.Id] = *info

		log.WithFields(log.Fields{
			"id":    event.Id,
			"image": event.From,
		}).Debug("Started monitor on new container")

	case "die":
		m_data.removeData(event.Id)

		log.WithFields(log.Fields{
			"id":    event.Id,
			"image": event.From,
		}).Debug("Removed monitor from container")
	}
}

func monitor(channel chan *monitorData) {
	// Monitor stuff
	log.Debug("I'm monitoring...")
	client := manager.Client
	m_data = monitorData{
		make(dckrStats),
		make(dckrInfo),
	}

	// Listen to events
	client.StartMonitorEvents(eventCallback, nil)

	// Get the list of active containers to monitor
	containers, err := client.ListContainers(false, false, "")
	if err != nil {
		log.Errorln("Error in listing containers: ", err)
		panic(err)
	}

	// Start the monitor for each active container
	for _, c := range containers {
		log.WithFields(log.Fields{
			"id":    c.Id,
			"image": c.Image,
		}).Debug("Started monitor on container")

		client.StartMonitorStats(c.Id, statCallBack, nil)

		// Get the info on containers
		info, err := client.InspectContainer(c.Id)
		if err != nil {
			log.Errorln("Error in inspecting containers: ", err)
			panic(err)
		}
		m_data.info[c.Id] = *info
	}

	// Set the ticker for the periodic update
	timeInterval := manager.Timer
	ticker := time.NewTicker(time.Duration(timeInterval) * time.Second)
	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			channel <- &m_data
		case <-quit:
			ticker.Stop()
			return
		}
	}
}
