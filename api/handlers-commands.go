package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/agent"
	ch "github.com/elleFlorio/gru/channels"
	com "github.com/elleFlorio/gru/communication"
	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/service"
)

type Command struct {
	Name      string
	Target    string
	Object    interface{}
	Result    string
	Timestamp time.Time
}

const c_GRU_REMOTE = "/gru/"
const c_CONFIG_REMOTE = "config"
const c_SERVICES_REMOTE = "services"
const c_TUNING_REMOTE = "tuning"

func PostCommand(w http.ResponseWriter, r *http.Request) {
	var err error
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	cmd, err := readCommand(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	executeCommand(cmd)
	w.WriteHeader(http.StatusAccepted)
}

func readCommand(r *http.Request) (Command, error) {
	var err error
	var cmd Command

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.WithField("err", err).Errorln("Error reading command body")
		return Command{}, err
	}

	if err = r.Body.Close(); err != nil {
		log.WithField("err", err).Errorln("Error closing command body")
		return Command{}, err
	}

	if err = json.Unmarshal(body, &cmd); err != nil {
		log.WithField("err", err).Errorln("Error unmarshaling command body")
		return Command{}, err
	}

	cmd.Timestamp = time.Now()

	log.WithFields(log.Fields{
		"name":      cmd.Name,
		"target":    cmd.Target,
		"timestamp": cmd.Timestamp,
	}).Debugln("Received command")

	return cmd, nil
}

func executeCommand(cmd Command) {
	switch cmd.Name {
	case "start":
		startCommand(cmd)
	case "stop":
		stopCommand(cmd)
	case "update":
		updateCommand(cmd)
	default:
		log.Errorln("Unrecognized command name: ", cmd.Name)
	}

	log.WithFields(log.Fields{
		"cmd":    cmd.Name,
		"target": cmd.Target,
	}).Debugln("Executed command")
}

func startCommand(cmd Command) {
	switch cmd.Target {
	case "agent":
		startCommunication()
		startAgent()
	case "service":
		name := cmd.Object.(string)
		startService(name)
	default:
		log.WithField("target", cmd.Target).Errorln("Unrecognized target for command start")
	}
}

func startAgent() {
	if !cfg.GetNode().Active {
		go runAgent()
	} else {
		log.Warnln("Node already active")
	}
}

func runAgent() {
	activateNode()
	defer deactivateNode()
	agent.Run()
}

func activateNode() {
	log.Debugln("Activating node")
	cfg.ToggleActiveNode()
	cfg.WriteNodeActive(cfg.GetNodeConfig().Remote, true)
}

func deactivateNode() {
	log.Debugln("Deactivating node")
	cfg.ToggleActiveNode()
	cfg.WriteNodeActive(cfg.GetNodeConfig().Remote, false)
}

func startCommunication() {
	com.Start(
		cfg.GetAgentCommunication().MaxFriends,
		cfg.GetAgentCommunication().LoopTimeInterval,
	)
}

func startService(name string) {
	log.WithField("name", name).Debugln("Starting service")
	toStart, err := service.GetServiceByName(name)
	if err != nil {
		log.WithField("name", name).Debugln("Error starting service")
	}
	ch.SendActionStartMessage(toStart)
}

func stopCommand(cmd Command) {
	switch cmd.Target {
	case "agent":
		//TODO
	case "service":
		name := cmd.Object.(string)
		stopService(name)
	default:
		log.WithField("target", cmd.Target).Errorln("Unrecognized target for command stop")
	}
}

func stopService(name string) {
	log.WithField("name", name).Debugln("Stopping service")
	toStop, err := service.GetServiceByName(name)
	if err != nil {
		log.WithField("name", name).Debugln("Error stopping service")
	}
	if len(toStop.Instances.All) < 1 {
		log.WithField("service", name).Debugln("No active instance to stop")
		return
	}

	ch.SendActionStopMessage(toStop)
}

func updateCommand(cmd Command) {
	log.Debugln("Updating ", cmd.Target)
	switch cmd.Target {
	case "all":
		cluster := cmd.Object.(string)
		updateAll(cluster)
	case "agent":
		cluster := cmd.Object.(string)
		updateAgent(cluster)
	case "services":
		cluster := cmd.Object.(string)
		updateServices(cluster)
	case "tuning":
		cluster := cmd.Object.(string)
		updateTuning(cluster)
	case "node-base-services":
		data := cmd.Object.([]interface{})
		upd := []string{}
		for _, item := range data {
			upd = append(upd, item.(string))
		}
		constraints := cfg.GetNodeConstraints()
		constraints.BaseServices = upd
		cfg.WriteNodeConstraints(cfg.GetNodeConfig().Remote, *constraints)
	case "node-cpumin":
		upd := cmd.Object.(float64)
		constraints := cfg.GetNodeConstraints()
		constraints.CpuMin = upd
		cfg.WriteNodeConstraints(cfg.GetNodeConfig().Remote, *constraints)
	case "node-cpumax":
		upd := cmd.Object.(float64)
		constraints := cfg.GetNodeConstraints()
		constraints.CpuMax = upd
		cfg.WriteNodeConstraints(cfg.GetNodeConfig().Remote, *constraints)
	case "service-mrt":
		name := cmd.Object.(string)
		srv, _ := service.GetServiceByName(name)
		upd := cfg.ReadService(srv.Remote)
		srv.Constraints.MaxRespTime = upd.Constraints.MaxRespTime
	default:
		log.WithField("target", cmd.Target).Errorln("Unrecognized target for command update")
	}
}

func updateAll(cluster string) {
	updateAgent(cluster)
	updateServices(cluster)
	updateTuning(cluster)
}

func updateAgent(cluster string) {
	configPath := c_GRU_REMOTE + cluster + "/" + c_CONFIG_REMOTE
	agentConfig := cfg.Agent{}
	cfg.ReadAgentConfig(configPath, &agentConfig)
	cfg.SetAgent(agentConfig)
	agent.UpdateStrategy()
	log.WithField("agent", agentConfig).Debugln("Agent updated from remote")
}

func updateServices(cluster string) {
	remote := c_GRU_REMOTE + cluster + "/" + c_SERVICES_REMOTE
	services := cfg.ReadServices(remote)
	cfg.SetServices(services)
	log.WithField("services", services).Debugln("Services updated from remote")
}

func updateTuning(cluster string) {
	remote := c_GRU_REMOTE + cluster + "/" + c_TUNING_REMOTE
	tuning := cfg.ReadTuningConfig(remote)
	cfg.SetTuning(tuning)
	log.WithField("tuning", tuning).Debugln("Tuning updated from remote")
}
