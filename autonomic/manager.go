package autonomic

import (
	log "github.com/Sirupsen/logrus"
	"github.com/samalba/dockerclient"
)

type autoManager struct {
	Timer  float32
	Client *dockerclient.DockerClient
}

var manager autoManager

func NewAutoManager(timer float32, client *dockerclient.DockerClient) *autoManager {
	manager = autoManager{
		timer,
		client,
	}
	return &manager
}

func (man *autoManager) RunLoop() {
	man.loop()
}

func (man *autoManager) loop() {
	log.Debugln("Started autonomic loop")
	//Channels initialization
	//TODO channel type
	monitorChannel := make(chan *monitorData)

	//Start the loop. Monitor works in the background
	go monitor(monitorChannel)

	for {
		select {
		case m_stats := <-monitorChannel:
			a_data := analyze(m_stats)
			p_act := plan(a_data)
			execute(p_act)
		}
	}
}
