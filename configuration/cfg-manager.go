package configuration

const (
	c_AGENT_DOCKER    = "docker"
	c_AGENT_AUTONOMIC = "autonomic"
	c_AGENT_COM       = "communication"
	c_AGENT_STORAGE   = "storage"
	c_AGENT_METRIC    = "metric"
	c_AGENT_DISCOVERY = "discovery"

	c_NODE_CONFIG      = "config"
	c_NODE_CONSTRAINTS = "constraints"
	c_NODE_RESOURCES   = "resources"
	c_NODE_INSTANCES   = "instances"
	c_NODE_ACTIVE      = "active"
)

var (
	agent       Agent
	node        Node
	services    []Service = []Service{}
	tuning      Tuning
	expressions map[string]Expression
)

func init() {
	expressions = make(map[string]Expression)
}

func SetAgent(cfg Agent) {
	agent = cfg
}

func GetAgent() *Agent {
	return &agent
}

func GetAgentDocker() *DockerConfig {
	return getAgentSubConfig(c_AGENT_DOCKER).(*DockerConfig)
}

func GetAgentAutonomic() *AutonomicConfig {
	return getAgentSubConfig(c_AGENT_AUTONOMIC).(*AutonomicConfig)
}

func GetAgentCommunication() *CommunicationConfig {
	return getAgentSubConfig(c_AGENT_COM).(*CommunicationConfig)
}

func GetAgentStorage() *StorageConfig {
	return getAgentSubConfig(c_AGENT_STORAGE).(*StorageConfig)
}

func GetAgentMetric() *MetricConfig {
	return getAgentSubConfig(c_AGENT_METRIC).(*MetricConfig)
}

func GetAgentDiscovery() *DiscoveryConfig {
	return getAgentSubConfig(c_AGENT_DISCOVERY).(*DiscoveryConfig)
}

func getAgentSubConfig(subCfg string) interface{} {
	switch subCfg {
	case c_AGENT_DOCKER:
		return &agent.Docker
	case c_AGENT_AUTONOMIC:
		return &agent.Autonomic
	case c_AGENT_COM:
		return &agent.Communication
	case c_AGENT_STORAGE:
		return &agent.Storage
	case c_AGENT_METRIC:
		return &agent.Metric
	case c_AGENT_DISCOVERY:
		return &agent.Discovery
	}

	return nil
}

func SetNode(cfg Node) {
	node = cfg
}

func GetNode() *Node {
	return &node
}

func GetNodeConfig() *NodeConfig {
	return getNodeSubConfig(c_NODE_CONFIG).(*NodeConfig)
}

func GetNodeConstraints() *NodeConstraints {
	return getNodeSubConfig(c_NODE_CONSTRAINTS).(*NodeConstraints)
}

func GetNodeResources() *NodeResources {
	return getNodeSubConfig(c_NODE_RESOURCES).(*NodeResources)
}

func GetNodeInstances() *ServiceStatus {
	return getNodeSubConfig(c_NODE_INSTANCES).(*ServiceStatus)
}

func getNodeSubConfig(subCfg string) interface{} {
	switch subCfg {
	case c_NODE_CONFIG:
		return &node.Configuration
	case c_NODE_CONSTRAINTS:
		return &node.Constraints
	case c_NODE_RESOURCES:
		return &node.Resources
	case c_NODE_INSTANCES:
		return &node.Instances
	}

	return nil
}

func ToggleActiveNode() {
	node.Active = !node.Active
}

func ClearNodeInstances() {
	node.Instances = ServiceStatus{}
}

func SetServices(cfg []Service) {
	services = cfg
}

func GetServices() []Service {
	return services
}

func AddServices(newServices []Service) {
	services = append(services, newServices...)
}

func RemoveServices(rmServices []string) {
	indexes := make([]int, len(rmServices), len(rmServices))

	for i, rmService := range rmServices {
		for j, service := range services {
			if service.Name == rmService {
				indexes[i] = j
			}
		}
	}

	for _, index := range indexes {
		services = append(services[:index], services[index+1:]...)
	}
}

func CleanServices() {
	services = make([]Service, 0)
}

func SetTuning(cfg Tuning) {
	tuning = cfg
}

func GetTuning() *Tuning {
	return &tuning
}

func SetExpr(cfg map[string]Expression) {
	expressions = cfg
}

func GetExpr() map[string]Expression {
	return expressions
}
