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
Gru will integrate seamlessly with you system based on containers: no need to do something strange, just start the Gru Agents in every node and let them manage your distributed application!

## Current status
The project is at an early stage of development.
That translates in something like version 0.0.0.0.0.0.0.0.0.0.0.0.1-pre_development_alpha.
I don't suggest to try it by now, however below you can find the documentation to run it in your local machine.
Currently Gru is working on a single node, scaling containers in order to balance the resource of that node between services according to the workload. The implementation of the distributed system is in progress; if you want to run Gru on multiple nodes you need an etcd server (https://github.com/coreos/etcd) for agents discovery.

## Documentation
These are the steps you need to follow to run the current version of Gru in your system (Linux only). Please remember that currently Gru is not able to really manage your containers and it's under active development.

###### Get Gru
`go get` this repo, or download it and compile/install using the go compiler

###### Create configuration files
* Create the folder `gru/config` in your home directory
* Inside the configuration folder you need to create 2 json files: `gruagentconfig.json` and `nodeconfig.json`. These files are needed to configure the agent and to provide information about the node. This is an example of the configuration files:
```json
//gruagentconfig.json
{
	"Service": {
		"ServiceConfigFolder":"/gru/config/services"
	},

	"Node": {
		"NodeConfigFile":"/gru/config/nodeconfig.json"
	},

	"Network": {
		"IpAddress":"127.0.0.1",
		"Port":"5000"
	},

	"Docker": {
		"DaemonUrl":"unix:///var/run/docker.sock",
		"DaemonTimeout":10
	},

	"Autonomic": {
		"LoopTimeInterval":5,
		"MaxFriends":5
	},

	"Discovery": {
		"DiscoveryService":"etcd",
		"DiscoveryServiceUri":"http://127.0.0.1:4001"
	},
	
	"Storage": {
		"StorageService":"internal"
	}
}
```
```json
//nodeconfig.json
{
	"Name":"node_name",
	"Constraints":{
		BaseServices:[]
	}
}
```
* Create the `services` folder at the location specified in `gruagentconfig.json`. Inside it create a `.json` file for each service you want to manage. Each service is bound to a Docker Image. This is an example of a service configuration file:
```json
//example.json
{
	"Name":"service_name",
	"Type":"service_type",
	"Image":"service_image",
	"Constraints":{
		"MaxRespTime":0
	},
	"Configuration":{
		"Cmd":[],
		"Volumes":null,
		"Entrypoint":[],
		"Memory":"0Gb",
		"CpuShares": 0,
		"CpusetCpus": "",
		"PortBindings": {},
		"Links":[]
	}
}
```
The field type can be used to specify the type of the service (e.g. Database, Webserver, etc.) and it is not used yet. The configuration field allow to specify some parameters related to the docker container. Please refer to the Docker documentation to an explanation of each configuration field.
###### Run Gru
* Run/start the containers of the services you want to manage
* Run the Gru agent with the command `gru start`. You can specify the logging level using the flag `-l`: e.g. `gru -l debug start`
* Enjoy

