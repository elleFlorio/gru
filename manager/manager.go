package manager

import (
	"fmt"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
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
			// signal the program to exit
			close(m.Quit)
		case "use":
			m.use(cmd)
		case "list":
			m.list(cmd)
		case "set":
			//TODO
		case "start":
			m.start(cmd)
		default:
			unknown(cmd)
		}

		return true
	}
	return false
}

func (m *Manager) list(cmd string) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 1, '\t', 0)
	args := strings.Split(strings.TrimSuffix(strings.TrimSpace(cmd), ";"), " ")
	if len(args) != 2 {
		fmt.Printf("Please specify clusters/nodes in command %q.\n", cmd)
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
	default:
		fmt.Println("Unrecognized identifier. Please specify clusters/nodes")
		return
	}

	for name, _ := range names {
		fmt.Fprintf(w, "%s\t%s\n", name, names[name])
	}
	w.Flush()
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

func (m *Manager) set(cmd string) {
	// args := strings.Split(strings.TrimSuffix(strings.TrimSpace(cmd), ";"), " ")
	// who := args[1]
	// what := args[2]
	// to_what := args[3:]

	// if m.Cluster == "" {
	// 	fmt.Println("Cluster not set. Please set it with 'use' command")
	// 	return
	// }
	// nodes := cluster.ListNodes(m.Cluster, false)
	// services := cluster.ListServices(m.Cluster)

	// switch what {
	// case "base-services":
	// 	if address, ok := nodes[who]; !ok {
	// 		fmt.Println("Unrecognized node ", who)
	// 	}
	// 	ok, notValid := checkValidServices(to_what, services)
	// 	if !ok {
	// 		fmt.Println("Services are not valid:")
	// 		for _, name := range notValid {
	// 			fmt.Println(name)
	// 		}
	// 		return
	// 	}
	// 	key := c_GRU_PATH + m.Cluster + "/" + c_NODES_PATH + who + "/constraints/baseservices"
	// 	value := to_what
	// 	opt := discovery.Options{}
	// 	discovery.Set(key, value, opt)
	// default:
	// 	fmt.Println("Unrecognized parameter ", what)
	// }

}

// func checkValidServices(services []string, list []string) (bool, []string) {
// 	check := true
// 	notValid := []string{}
// 	for _, s1 := range services {
// 		valid := false
// 		for _, s2 := range list {
// 			valid = valid || (s1 == s2)
// 		}
// 		if valid == false {
// 			notValid = append(notValid, s1)
// 		}
// 		check = check && valid
// 	}

// 	return check, notValid
// }

func (m *Manager) start(cmd string) {
	args := strings.Split(strings.TrimSuffix(strings.TrimSpace(cmd), ";"), " ")
	if len(args) < 2 {
		fmt.Println("not enough arguments to 'start' command")
		return
	}

	if !m.isClusterSet() {
		return
	}

	var err error
	nodes := cluster.ListNodes(m.Cluster, false)
	if args[1] == "all" {
		for name, addr := range nodes {
			err = network.SendStartCommand(addr)
			if err != nil {
				fmt.Println("Error sending start command to node ", name)
			}
		}
	} else {
		for _, node := range args[1:] {
			if addr, ok := nodes[node]; ok {
				err = network.SendStartCommand(addr)
				if err != nil {
					fmt.Println("Error sending start command to node ", node)
				}
			} else {
				fmt.Println("Cannot get address of node ", node)
			}
		}
	}

}

func unknown(cmd string) {
	fmt.Printf("Unknown command %q.\n", cmd)
}

func (m *Manager) isClusterSet() bool {
	isSet := m.Cluster != ""
	if !isSet {
		fmt.Println("Cluster not set.")
	}
	return isSet
}
