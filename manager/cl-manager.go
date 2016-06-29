package manager

// This was inspired by the command line client of influxdb:
// https://github.com/influxdb/influxdb/blob/master/cmd/influx/cli/cli.go

import (
	"fmt"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"text/tabwriter"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/peterh/liner"

	"github.com/elleFlorio/gru/cluster"
	"github.com/elleFlorio/gru/discovery"
	"github.com/elleFlorio/gru/network"
)

const c_GRU_PATH = "/gru/"
const c_NODES_PATH = "nodes/"

type Manager struct {
	Remote      discovery.Discovery
	Cluster     string
	Line        *liner.State
	Quit        chan struct{}
	osSignals   chan os.Signal
	historyFile *os.File
}

func New(etcdAddr string) (*Manager, error) {
	remote, err := discovery.New("etcd", etcdAddr)
	if err != nil {
		return nil, err
	}

	return &Manager{
		Remote:    remote,
		Quit:      make(chan struct{}, 1),
		osSignals: make(chan os.Signal, 1),
	}, nil
}

func (m *Manager) Run() {
	var err error
	signal.Notify(m.osSignals, os.Kill, os.Interrupt, syscall.SIGTERM)

	m.Line = liner.NewLiner()
	defer m.Line.Close()
	m.Line.SetMultiLineMode(true)

	var historyFilePath string
	usr, err := user.Current()
	// Only load/write history if we can get the user
	if err == nil {
		historyFilePath = filepath.Join(usr.HomeDir, ".gru_manager_history")
		if m.historyFile, err = os.OpenFile(historyFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640); err == nil {
			defer m.historyFile.Close()
			m.Line.ReadHistory(m.historyFile)
		}
	}

	fmt.Println("")
	fmt.Println("###################################################")
	fmt.Println("**           GRU COMMAND LINE MANAGER            **")
	fmt.Println("** developed by Luca Florio (github: elleFlorio) **")
	fmt.Println("###################################################")
	fmt.Println("")
	for {
		select {
		case <-m.osSignals:
			close(m.Quit)
		case <-m.Quit:
			m.exit()
		default:
			l, e := m.Line.Prompt("> ")
			if e != nil {
				break
			}
			if m.ParseCommand(l) {
				m.Line.AppendHistory(l)
				_, err := m.Line.WriteHistory(m.historyFile)
				if err != nil {
					fmt.Printf("There was an error writing history file: %s\n", err)
				}
			}
		}
	}
}

func (m *Manager) exit() {
	// write to history file
	_, err := m.Line.WriteHistory(m.historyFile)
	if err != nil {
		fmt.Printf("There was an error writing history file: %s\n", err)
	}
	// release line resources
	m.Line.Close()
	m.Line = nil
	// exit CLI
	os.Exit(0)
}

func (m *Manager) ParseCommand(cmd string) bool {
	lcmd := strings.TrimSpace(strings.ToLower(cmd))
	tokens := strings.Fields(lcmd)

	if len(tokens) > 0 {
		switch tokens[0] {
		case "exit":
			close(m.Quit)
		case "use":
			m.use(cmd)
		case "list":
			m.list(cmd)
		case "set":
			m.set(cmd)
		case "show":
			m.show(cmd)
		case "start":
			m.start(cmd)
		case "stop":
			m.stop(cmd)
		case "update":
			m.update(cmd)
		case "deploy":
			m.deploy()
		case "undeploy":
			m.undeploy()
		default:
			unknown(cmd)
		}

		return true
	}
	return false
}

func (m *Manager) use(cmd string) {
	args := strings.Split(strings.TrimSuffix(strings.TrimSpace(cmd), ";"), " ")
	if len(args) != 2 {
		fmt.Printf("Could not parse cluster name from %q.\n", cmd)
		return
	}
	d := args[1]
	m.Cluster = d
	fmt.Printf("Using cluster %s\n", d)
}

func (m *Manager) list(cmd string) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 1, '\t', 0)
	args := strings.Split(strings.TrimSuffix(strings.TrimSpace(cmd), ";"), " ")
	if len(args) != 2 {
		fmt.Printf("Not enough arguments %q.\n", cmd)
		return
	}

	var names map[string]string
	switch args[1] {
	case "clusters":
		names = cluster.ListClusters()
		fmt.Fprintf(w, "NAME\tUUID\n")
	case "nodes":
		if !m.isClusterSet() {
			return
		}
		names = cluster.ListNodes(m.Cluster, false)
		fmt.Fprintf(w, "NAME\tADDRESS\n")
	case "services":
		if !m.isClusterSet() {
			return
		}
		names = cluster.ListServices(m.Cluster)
		fmt.Fprintf(w, "NAME\tIMAGE\n")
	default:
		fmt.Println("Unrecognized identifier. Please specify clusters/nodes")
		return
	}

	for name, _ := range names {
		fmt.Fprintf(w, "%s\t%s\n", name, names[name])
	}
	w.Flush()
}

func (m *Manager) set(cmd string) {
	args := strings.Split(strings.TrimSuffix(strings.TrimSpace(cmd), ";"), " ")
	where := args[1]
	who := args[2]
	what := args[3]
	to_what := args[4:]

	if !m.isClusterSet() {
		return
	}

	switch where {
	case "node":
		setNode(m.Cluster, who, what, to_what)
	case "service":
		setService(m.Cluster, who, what, to_what)
	default:
		fmt.Println("Unrecognized identifier. Please specify node/service")
	}

}

func setNode(clusterName string, who string, what string, to_what []string) {
	nodes := cluster.ListNodes(clusterName, false)
	dest := []string{}

	if who == "all" {
		for _, address := range nodes {
			dest = append(dest, address)
		}
	} else {
		if address, ok := nodes[who]; !ok {
			fmt.Println("Unrecognized node ", who)
			return
		} else {
			dest = append(dest, address)
		}
	}

	switch what {
	case "base-services":
		services := cluster.ListServices(clusterName)
		names := make([]string, 0, len(services))
		for k, _ := range services {
			names = append(names, k)
		}
		ok, notValid := checkValidServices(to_what, names)
		if !ok {
			fmt.Println("Services are not valid:")
			for _, name := range notValid {
				fmt.Println(name)
			}
			return
		}
		for _, address := range dest {
			err := network.SendUpdateCommand(address, "node-base-services", to_what)
			if err != nil {
				fmt.Println("Error sending update command to ", address)
			}
		}
	case "cpumin":
		cpumin := to_what[0]
		if ok, value := checkValidCpuValue(cpumin); ok {
			for _, address := range dest {
				err := network.SendUpdateCommand(address, "node-cpumin", value)
				if err != nil {
					fmt.Println("Error sending update command to ", address)
				}
			}
		} else {
			fmt.Println("CPU value not valid: it should be a float between 0.0 and 1.0")
		}
	case "cpumax":
		cpumax := to_what[0]
		if ok, value := checkValidCpuValue(cpumax); ok {
			for _, address := range dest {
				err := network.SendUpdateCommand(address, "node-cpumax", value)
				if err != nil {
					fmt.Println("Error sending update command to ", address)
				}
			}
		} else {
			fmt.Println("CPU value not valid: it should be a float between 0.0 and 1.0")
		}
	default:
		fmt.Println("Unrecognized parameter ", what)
	}
}

func checkValidCpuValue(cpuval string) (bool, float64) {
	value, err := strconv.ParseFloat(cpuval, 64)
	if err != nil {
		return false, 0.0
	}
	if value < 0.0 || value > 1.0 {
		return false, 0.0
	}

	return true, value
}

func setService(clusterName string, who string, what string, to_what []string) {
	nodes := cluster.ListNodes(clusterName, false)
	names := []string{}
	for name, _ := range cluster.ListServices(clusterName) {
		names = append(names, name)
	}

	ok, _ := checkValidServices([]string{who}, names)
	if !ok {
		fmt.Println("Unrecognized service ", who)
		return
	}
	service := cluster.GetService(clusterName, who)
	switch what {
	case "TODO":
		cluster.UpdateService(clusterName, who, service)
		for _, address := range nodes {
			network.SendUpdateCommand(address, "TODO", who)
		}
	default:
		fmt.Println("Unrecognized parameter ", what)
	}
}

func checkValidServices(services []string, list []string) (bool, []string) {
	check := true
	notValid := []string{}
	for _, s1 := range services {
		valid := false
		for _, s2 := range list {
			valid = valid || (s1 == s2)
		}
		if valid == false {
			notValid = append(notValid, s1)
		}
		check = check && valid
	}

	return check, notValid
}

func (m *Manager) start(cmd string) {
	args := strings.Split(strings.TrimSuffix(strings.TrimSpace(cmd), ";"), " ")
	if len(args) < 3 {
		fmt.Println("not enough arguments to 'start' command")
		return
	}

	if !m.isClusterSet() {
		return
	}

	who := args[1]

	switch who {
	case "agent":
		m.startAgent(args[2:])
	case "service":
		m.startService(args[2], args[3:])
	default:
		fmt.Println("Unrecognized target ", who)
	}
}

func (m *Manager) startAgent(where []string) {
	var err error
	nodes := cluster.ListNodes(m.Cluster, false)
	if where[0] == "all" {
		for name, addr := range nodes {
			err = network.SendStartAgentCommand(addr)
			if err != nil {
				fmt.Println("Error sending start agent command to node ", name)
			}
		}
	} else {
		for _, node := range where {
			if addr, ok := nodes[node]; ok {
				err = network.SendStartAgentCommand(addr)
				if err != nil {
					fmt.Println("Error sending start agent command to node ", node)
				}
			} else {
				fmt.Println("Cannot get address of node ", node)
			}
		}
	}
}

func (m *Manager) startService(what string, where []string) {
	var err error
	nodes := cluster.ListNodes(m.Cluster, false)
	if where[0] == "all" {
		for name, addr := range nodes {
			err = network.SendStartServiceCommand(addr, what)
			if err != nil {
				fmt.Println("Error sending start service command to node ", name)
			}
		}
	} else {
		for _, node := range where {
			if addr, ok := nodes[node]; ok {
				err = network.SendStartServiceCommand(addr, what)
				if err != nil {
					fmt.Println("Error sending start service command to node ", node)
				}
			} else {
				fmt.Println("Cannot get address of node ", node)
			}
		}
	}
}

func (m *Manager) stop(cmd string) {
	args := strings.Split(strings.TrimSuffix(strings.TrimSpace(cmd), ";"), " ")
	if len(args) < 3 {
		fmt.Println("not enough arguments to 'stop' command")
		return
	}

	if !m.isClusterSet() {
		return
	}

	who := args[1]

	switch who {
	case "agent":
		m.stopAgent(args[2:])
	case "service":
		m.stopService(args[2], args[3:])
	default:
		fmt.Println("Unrecognized target ", who)
	}
}

func (m *Manager) stopAgent(where []string) {
	var err error
	nodes := cluster.ListNodes(m.Cluster, false)
	if where[0] == "all" {
		for name, addr := range nodes {
			err = network.SendStopAgentCommand(addr)
			if err != nil {
				fmt.Println("Error sending stop agent command to node ", name)
			}
		}
	} else {
		for _, node := range where {
			if addr, ok := nodes[node]; ok {
				err = network.SendStopAgentCommand(addr)
				if err != nil {
					fmt.Println("Error sending stop agent command to node ", node)
				}
			} else {
				fmt.Println("Cannot get address of node ", node)
			}
		}
	}
}

func (m *Manager) stopService(what string, where []string) {
	var err error
	nodes := cluster.ListNodes(m.Cluster, false)
	if where[0] == "all" {
		for name, addr := range nodes {
			err = network.SendStopServiceCommand(addr, what)
			if err != nil {
				fmt.Println("Error sending start service command to node ", name)
			}
		}
	} else {
		for _, node := range where {
			if addr, ok := nodes[node]; ok {
				err = network.SendStopServiceCommand(addr, what)
				if err != nil {
					fmt.Println("Error sending start service command to node ", node)
				}
			} else {
				fmt.Println("Cannot get address of node ", node)
			}
		}
	}
}

func (m *Manager) update(cmd string) {
	args := strings.Split(strings.TrimSuffix(strings.TrimSpace(cmd), ";"), " ")
	if len(args) < 2 {
		fmt.Println("not enough arguments to 'update' command")
		return
	}

	if !m.isClusterSet() {
		return
	}

	var err error
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 1, '\t', 0)
	fmt.Fprintf(w, "NODE\tSTATUS\n")

	nodes := cluster.ListNodes(m.Cluster, false)
	for node, address := range nodes {
		status := "done"
		err = network.SendUpdateCommand(address, args[1], m.Cluster)
		if err != nil {
			fmt.Println("Error sending update command to node ", node)
			status = "error"
		}

		fmt.Fprintf(w, "%s\t%s\n",
			node,
			status,
		)
	}

	w.Flush()
}

func (m *Manager) show(cmd string) {
	args := strings.Split(strings.TrimSuffix(strings.TrimSpace(cmd), ";"), " ")
	if len(args) < 4 {
		fmt.Println("not enough arguments to 'show' command")
		return
	}

	where := args[1]
	who := args[2]
	what := args[3]

	switch where {
	case "cluster":
		fmt.Println("TODO")
	case "service":
		fmt.Println("TODO")
	case "node":
		if !m.isClusterSet() {
			return
		}
		showNode(m.Cluster, who, what)
	}
}

func showNode(clusterName string, nodeName string, what string) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 1, '\t', 0)
	nodes := cluster.GetNodes(clusterName, false)
	if nodeName != "all" {
		nodes := cluster.ListNodes(clusterName, false)
		if _, ok := nodes[nodeName]; !ok {
			fmt.Println("Unrecognized node ", nodeName)
			return
		}
	}

	switch what {
	case "config":
		fmt.Fprintf(w, "NAME\tUUID\tADDRESS\tCLUSTER\tREMOTE\n")
		for _, node := range nodes {
			config := node.Configuration
			if nodeName != "all" {
				if node.Configuration.Name == nodeName {
					fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
						config.Name,
						config.UUID,
						config.Address,
						config.Cluster,
						config.Remote,
					)
				}
			} else {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
					config.Name,
					config.UUID,
					config.Address,
					config.Cluster,
					config.Remote,
				)
			}
		}
	case "constraints":
		fmt.Fprintf(w, "NAME\tBASE-SERVICES\tCPU-MIN\tCPU-MAX\n")
		for _, node := range nodes {
			constraints := node.Constraints
			if nodeName != "all" {
				if node.Configuration.Name == nodeName {
					fmt.Fprintf(w, "%s\t%v\t%f\t%f\n",
						node.Configuration.Name,
						constraints.BaseServices,
						constraints.CpuMin,
						constraints.CpuMax,
					)
				}
			} else {
				fmt.Fprintf(w, "%s\t%v\t%f\t%f\n",
					node.Configuration.Name,
					constraints.BaseServices,
					constraints.CpuMin,
					constraints.CpuMax,
				)
			}
		}
	case "resources":
		fmt.Fprintf(w, "NAME\tCORES\tMEMORY\n")
		for _, node := range nodes {
			resources := node.Resources
			if nodeName != "all" {
				if node.Configuration.Name == nodeName {
					fmt.Fprintf(w, "%s\t%d\t%d\n",
						node.Configuration.Name,
						resources.TotalCpus,
						resources.TotalMemory,
					)
				}
			} else {
				fmt.Fprintf(w, "%s\t%d\t%d\n",
					node.Configuration.Name,
					resources.TotalCpus,
					resources.TotalMemory,
				)
			}
		}
	default:
		fmt.Println("Unrecognized node property ", what)
	}

	w.Flush()
}

func (m *Manager) deploy() {
	if !m.isClusterSet() {
		return
	}

	var err error
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 1, '\t', 0)
	fmt.Fprintf(w, "NODE\tSERVICE\tSTATUS\n")

	services := cluster.ListServices(m.Cluster)
	servicesNames := make([]string, 0, len(services))
	for name, _ := range services {
		servicesNames = append(servicesNames, name)
	}

	nodes := cluster.ListNodes(m.Cluster, false)
	nodesNames := make([]string, 0, len(nodes))
	for name, _ := range nodes {
		nodesNames = append(nodesNames, name)
	}

	if len(nodes) < len(services) {
		fmt.Println("Cannot deploy: not enough nodes")
		return
	}

	for i := 0; i < len(services); i++ {
		service := servicesNames[i]
		node := nodesNames[i]
		address := nodes[node]
		status := "done"

		err = network.SendUpdateCommand(address, "node-base-services", []string{service})
		if err != nil {
			fmt.Println("Error sending update command to ", address)
			status = "error"

		} else {
			err = network.SendStartServiceCommand(address, service)
			if err != nil {
				fmt.Println("Error sending start service command to node ", node)
				status = "error"
			}
		}

		fmt.Fprintf(w, "%s\t%s\t%s\n",
			node,
			service,
			status,
		)

	}

	w.Flush()
}

func (m *Manager) undeploy() {
	if !m.isClusterSet() {
		return
	}

	var err error
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 1, '\t', 0)
	fmt.Fprintf(w, "NODE\tSTATUS\n")

	services := cluster.ListServices(m.Cluster)
	nodes := cluster.ListNodes(m.Cluster, false)
	for node, address := range nodes {
		status := "done"
		err = network.SendUpdateCommand(address, "node-base-services", []string{})
		if err != nil {
			fmt.Println("Error sending update command to ", address)
			status = "error"
		}

		for service, _ := range services {
			err = network.SendStopServiceCommand(address, service)
			if err != nil {
				fmt.Println("Error sending stop service command to node ", node)
				status = "error"
			}
		}

		fmt.Fprintf(w, "%s\t%s\n",
			node,
			status,
		)
	}

	w.Flush()
}

func (m *Manager) isClusterSet() bool {
	isSet := m.Cluster != ""
	if !isSet {
		fmt.Println("Cluster not set.")
	}
	return isSet
}

func unknown(cmd string) {
	fmt.Printf("Unknown command %q.\n", cmd)
}
