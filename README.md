# Gru: automatic management of Docker containers
[![Build Status](https://travis-ci.org/elleFlorio/gru.svg?branch=master)](https://travis-ci.org/elleFlorio/gru)
[![Coverage Status](https://coveralls.io/repos/elleFlorio/gru/badge.svg?branch=master&service=github)](https://coveralls.io/github/elleFlorio/gru?branch=master)

Gru brings self-adaptation to Docker containers: it manages Docker containers distributed in a set of nodes scaling them and actuating autonomic actions to keep everything up and running. Gru is designed to help developers build distributed applications based on microservices running in Docker containers.

Gru is part of my PhD research at Politecnico di Milano: "Decentralized Self-Adaptation in Large-Scale Distributed Systems" and it is currenlty under development.

## VERY IMPORTANT
This is a PhD project: this means that by now it is more like a "proof of concept".
I'm working to make Gru as close as possible to a "real" software, but please keep in mind that I'm the only person actively working on it and I'm learning the technology through the development process. So I am open to any suggestion/contribution.
If you like my work and want to help me in some way, you can contact me at **luca[dot]florio[at]polimi[dot]it**:
* if you are a student at Politecnico di Milano, we can discuss a thesis;
* if you are just curious about my work/want to give me a suggestion, I'm happy to have a chat;
* if you are a company and want to offer me an internship/job... Well, you make me very happy! :-)

## Goals
**Distributed**
Gru will be able to automagically manage containers distributed in a huge number of nodes

**Decentralized**
Gru will use a decentralized approach based on the idea of self-organizing multiagent system: Gru Agents are deployed in every node and communicate with the Docker daemon and with other agents to monitor the system and plan the best action to actuate. In this way there is no single point of failure.

**Plug & Play**
Develop your contenerized application with no worries: Gru will integrate seamlessly with you system based on containers. No need to do something strange, just start the Gru Agents in every node and let them manage your distributed application!

## Current status
The project is at an early stage of development.
I don't suggest to try it by now, however below you can find the documentation to run it in your local machine or cluster.
Currently Gru can work in a cluster of nodes, automagically scaling services instances according to traffic load.
Gru needs a running instance of an [etcd](https://github.com/coreos/etcd) for agents discovery and [InfluxDB](https://influxdb.com/) to store metrics data.

## Documentation (Work in Progress)
These are the steps you need to follow to run the current version of Gru in your system. This is not a documentation of all the aspects involved, but it should be enough to understand a little better the project and maybe test it.
Please remember that currently Gru is able only to autoscale your services containers and it's under active development.

### Requirements
Gru needs some external components and some environment variables to run. Gru has been developed and tested in linux, so the best thing is if you use the same environment.

#### External components
* [Docker](https://www.docker.com/)
* [etcd](https://github.com/coreos/etcd)
* [InfluxDB](https://www.influxdata.com/time-series-platform/influxdb/) (optional)

**Docker** is the obvious requirement and can be used also to run an instance of the Gru Agent (more on this later).

**etcd** can be run in a Docker container. This is the command to run a local instance of etcd to test Gru:
```
docker run -d -p 4001:4001 -p 2380:2380 -p 2379:2379 --name etcd quay.io/coreos/etcd \
-name etcd0 \
-advertise-client-urls http://${HostIP}:2379,http://${HostIP}:4001 \
-listen-client-urls http://0.0.0.0:2379,http://0.0.0.0:4001
```

**InfluxDB** is not required, but recommended. It doesn't requires a particular configuration, you just need to create a database that will be used to store the metrics sent by Gru Agents.

#### Environment Variables
Gru needs some environment variables to be set in order to work properly.
* HostIP: the IP of the host machine. Here we assume that etcd and InfluxDB are running inside the same host. A smart trick is to set it automatically using this command:
```
export HostIP=$(ifconfig enp0s3 | grep "inet addr" | awk -F: '{print $2}' | awk '{print $1}')
```
* ETCD_ADDR: the address of the running instance of etcd.
```
export ETCD_ADDR=http://$HostIP:4001
```
* Variables related to InfluxDB, such as:
```
export INFLUX_ADDR=http://$HostIP:8086
export INFLUX_USER=<db_user>
export INFLUX_PWD=<db_pwd>
```

### Run Gru
Gru can be run in two ways:
* Downloading and compiling the sources in this repo.
* Using the Docker image `elleflorio/gru` (use `dev` tag for the development branch image) **[suggested method]**.

In the latter case, Docker should be configured to listen on a specific port (usually :2375) for commands.
To configure properly Docker in linux, follow these instruction (credits to [docker-ui](https://github.com/felixgborrego/docker-ui-chrome-app/wiki/linux)).

#### How to enable Docker Remote API on Linux
We'll need to enable the Docker Remote API, but first make sure Docker daemon is up an running using ```docker info```.

#####Linux with systemd (Ubuntu 15.04, Debian 8,...)

Using systemd, we'll need to enable a systemd socket to access the Docker remote API:

* Create a new systemd config file called ```/etc/systemd/system/docker-tcp.socket``` to make docker available on a TCP socket on port 2375.
```
[Unit]
Description=Docker HTTP Socket for the API

[Socket]
ListenStream=2375
BindIPv6Only=both
Service=docker.service

[Install]
WantedBy=sockets.target
```

* Register the new systemd http socket and restart docker
```
systemctl enable docker-tcp.socket
systemctl stop docker
systemctl start docker-tcp.socket
```

* Open your browser and verify you can connect to http://localhost:2375/_ping

##### Linux  without systemd

You need to enable Docker Remote API. To do so you need to:
* Edit ```/etc/default/docker``` to allow connections adding:
```DOCKER_OPTS='-H tcp://0.0.0.0:2375 -H unix:///var/run/docker.sock'```
* Restart the Docker service using:
```sudo service docker restart```
* Open your browser and verify you can connect to http://localhost:2375/_ping

#### Command Line Client
Gru is provided with a command line client that allows the user to interact with Gru and to execute basic operations.

#### Create a Cluster
The first operation to do is to create a Cluster (if it has not been created yet). The command to create it is the following:
```
gru create <cluster_name>
```
This command will create the correct folder tree in etcd to store the configuration of the Cluster:

#### Join a Cluster
Gru agents can join a cluster using the join command followed by the name of the cluster:
```
gru join <cluster_name> [--name <node_name>]
```
with the flag `--name` (or `-n`) it is possible to set a specific name for the node.

#### Manage the Cluster
Using the command `gru manage` it is possible to manage a cluster of nodes. Here I present some basic commands that can be used in the command line client to deploy the services of the application and start the Gru Agents to manage them.
* `use <cluster_name>`: chose the cluster to manage
* `list nodes`: list the current active nodes in the cluster
* `deploy`: deploy the services of the application to manage (one in each node) according to the ones provided in the configuration
* `set node <node_name> base-services <services_names>`: set the base-service property in the node. If a service is the base-services list, the Gru Agent ensure that a instance of that service will always be running in the node.
* `start service <service_name> node <node_name>`: start an instance of the service in the node
* `start agent <node_name>`: start the agent in the node. To start all the agent use the command `start agent all`

### Configuration
Gru needs some parameters to be configured correctly. The configuration is stored inside etcd in the form of strings as represented in the following tree:
```
/gru/ 
	|---<cluster>/ 
		|---uuid:string (Cluster ID) 
		|---nodes/ 
			|---<node> 
				|---configuration:string 
				|---constraints:string 
				|---resources:string 
				|---active:string 
		|---config (Agent configuration) 
		|---services/ 
			|---<service>:string (Services configuration) 
		|---policy (Policies configuration) 
		|--- analytics/ 
			|---<analytic>:string (Analytics configuration)
```
The configuration is expressed in the JSON format and can be set into etcd using a web interface called [**Gru Configuration Manager**](https://github.com/elleFlorio/gruConfigurationManager).
A Docker image is available to easily run the configuration manager: `elleflorio/grucm`.

The configuration is composed of several files that are described here:

#### Agent Configuration
this file contains the configuration of the Gru Agent. This is an example of the agent configuration.
```
{
	"Docker": {
		"DaemonUrl":"local:2375", //to run Gru inside a Docker container
		"DaemonTimeout":10
	},
	"Autonomic": {
		"LoopTimeInterval":60,
		"PlannerStrategy":"probdelta",
		"EnableLogReading": true
	},
	"Communication":{
		"LoopTimeInterval":55,
		"MaxFriends":5
		},
	"Storage": {
		"StorageService":"internal"
	},
	"Metric": {
		"MetricService":"influxdb",
		"Configuration": {
			"Url":"http://<influx_address>:<influx_port>",
			"DbName": "<gru_db_name>",
			"Username": "<influx_user>",
			"Password": "<influx_pwd<"
		}
	},
	"Discovery": {
		"AppRoot":"app/<app_name>/services",
		"TTL":5
	}
}
```

#### Analytics
The user can provide some analytics that should be computed by Gru Agents for the services. The user should provide an equation that will be evaluated as a value between 0 and 1 that involves the use of some metrics/constraints. The user should create a specific configuration for each analytic, that needs to be composed as follows.
```
{
	"Name": "<analytic_name>",
	"Expr": "<analytic_equation>",
	"Metrics": [<list_of_metrics_involved>],
	"Constraints": [<list_of_constraints_involved>]
}
```

This is an example of a possible `response_time_ratio` analytic, used to understand if a service has a response time that is too high.
```
{
	"Name": "resp_time_ratio",
	"Expr": "execution_time / MAX_RESP_TIME",
	"Metrics": [
		"execution_time"
	],
	"Constraints": [
		"MAX_RESP_TIME"
	]
}
```

#### Policy configuration
This configuration file allows to set the parameters of the policies implemented in the system, as well as enable or disable them.
```
{
	"Scalein": {
		"Enable": true,
		"Threshold": 0.35,
		"Metrics": [<list_of_metrics_involved>],
		"Analytics": [<list_of_analytics_involved>]
	},
	"Scaleout": {
		"Enable": true,
		"Threshold": 0.75,
		"Metrics": [<list_of_metrics_involved>],
		"Analytics": [<list_of_analytics_involved>]
	},
	"Swap": {
		"Enable": true,
		"Threshold": 0.6,
		"Metrics": [<list_of_metrics_involved>],
		"Analytics": [<list_of_analytics_involved>]
	}
}
```

#### Services Descriptors
For each service to be managed, Gru Agents needs to know some information related to the service. This means that it is required a different service descriptor for each service composing the application that we want to manage. The following is an example of a service descriptor.
```
{
	"Name":"<service_name>",
	"Type":"<service_type>",
	"Image":"<service_image>",
	"Remote":"/gru/<cluster_name>/services/<service_name>",
	"DiscoveryPort":"<service_port_for_discovery>",
	"Analytics": [<list_of_analytics>],
	"Constraints":{<key_value_contraints_map>},
	"Configuration":{<docker_configuration>}
}
```
This is an example of the configuration of a service called `service1` in Cluster "myCluster".
```
{
	"Name":"service1",
	"Type":"service1",
	"Image":"elleflorio/service1",
	"Remote":"/gru/myCluster/services/service1",
	"DiscoveryPort":"50000",
	"Analytics": [
		"analytic1",
		"analytic2"
	],
	"Constraints":{
		"MAX_RESP_TIME":1000 
	},
	"Configuration":{
		"cpunumber":1,
		"StopTimeout":30,
		"Env": {
            "ETCD_ADDR":"",
            "HostIP":"",
            "INFLUX_USER":"myUser",
            "INFLUX_PWD":"myPwd",
            "INFLUX_ADDR":"http://192.168.1.1:8080"
        },
		"Ports":{
			"50000":"50000-50004"
		},
		"Cmd":[
			"start",
			"service1"
		]
	}
}
```

### Example Deployment
This is an example of a deployment process. The assumption is that the requirements are met (external tools up and running, env vars set, etc.).
Our cluster is composed of 5 working-nodes and 1 main-node. The external components (etcd, InfluxDB) are deployed in the main node. The working-nodes will be used to deploy our application and will be the hosts of our Gru Agents. The tool Gru is available in all the nodes.

We start connecting to a node and we create our Cluster:
```
gru create myCluster
```

We create our configurations and upload it using the Gru Configuration Manager, and the result is the following Cluster configuration:
```
/gru/ 
	|---myCluster/ 
		|---uuid:<myCluster_ID> 
		|---nodes/ 
			|---<empty>
		|---config 
		|---services/ 
			|---service1
			|---service2
			|---service3 
		|---policy 
		|--- analytics/ 
			|---resp_time_ratio
```

Once our Cluster is ready and the configuration is uploaded, we can register the Gru Agents in every working-node using the `join` command. We suppose that every node has a name that is "node<node_number>" (e.g., "node1", "node2", etc.), so we connect to every working-node and we run the following command:
```
gru join myCluster --name node<node_number>
```
The result is that the list of nodes in the Cluster will be updated with the information of every node:
```
/gru/ 
	|­­­---myCluster/ 
		|­­­---uuid:<myCluster_ID> 
		|---nodes/ 
			|---node1
				|---...
			|---node2
				|---...
			|---node3
				|---...
			|---node4
				|---...
			|---node5
				|---...
		|---config 
		|---services/ 
			|---service1
			|---service2
			|---service3 
		|---policy 
		|--- analytics/ 
			|---resp_time_ratio
```
Obviously we can configure a script to run the `join` command when a working-node starts.

We are ready to deploy our application and to start the management of services. We need to run the command line client connecting to a host and typing the command:
```
gru manage
```
Now we are in the command line client. First of all we need to chose the Cluster:
```
> use myCluster
```
Once the Cluster is set, we can list the active nodes:
```
> list nodes
```
The system will print the list of nodes, as a table name-ip_address and the total number of active nodes. Now we can deploy our application
```
> deploy
```
This command will start a container running every service present in the configuration in different nodes. In our case, suppose the system starts service1 in node1, service2 in node2 and service3 in node3. The `deploy` command also set the base-service properties for every service in the node where it is started, so we will be sure that there will be always a running instances of our services inside the Cluster.
We are ready to start the Gru Agents that will manage our application:
```
> start agent all
```
The application is up and running and our agents are managing it! Enjoy! ;-)
