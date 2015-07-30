# Gru: automatic management of Docker containers
[![Build Status](https://travis-ci.org/elleFlorio/gru.svg?branch=master)](https://travis-ci.org/elleFlorio/gru)

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
Gru will be able to automatic manage containers distributed in a huge number of nodes

**Decentralized**
Gru will use a decentralized approach based on the idea of self-organizing multi-agets system: Gru Agents are deployed in every node and communicate with the Docker daemon and with other agents to monitor the system and plan the best action to actuate. In this way there is no single point of failure.

**Plug & Play**
Gru will integrate seamlessly with you system based on containers: no need to do something strange, just start the Gru Agents in every node and let them manage your distributed application!

## Current status
The project is at an early stage of development.
That translates in something like version 0.0.0.0.0.0.0.0.0.0.0.0.1-pre_development_alpha.
I don't suggest to try it by now, however I will post all the instructions needed to run it ASAP. The same goes for godoc documentation.
Currently Gru is working on a single node, scaling containers in order to balance the resource of that node between services according to the workload.

## Instruction
These are the steps you need to follow to run the current version of Gru in your system (Linux only). Please remember that currently Gru is not able to really manage your containers and it's under active development.

###### Get Gru
`go get` this repo, or download it and compile/install using the go compiler

###### Create configuration files
* Create the folder `gru/config` in your home directory
* Inside the configuration folder you need to create 2 json files: `gruagentconfig.json` and `nodeconfig.json`. These files are needed to configure the agent and to provide information about the node. This is an example of the configuration files:
```json
//gruagentconfig.json
{
	"DaemonUrl":"unix:///var/run/docker.sock",
	"DaemonTimeout":10,
	"LoopTimeInterval":5,
	"ServiceConfigFolder":"/home/gru/config/services"
}
```
```json
//nodeconfig.json
{
	"Name":"node1",
	"Constraints":{
		"CpuMax":0.8,
		"CpuMin":0.2,
		"MaxInstances":8
	}
}
```
* Create the `services` folder at the location specified in `gruagentconfig.json`. Inside it create a `.json` file for each service you want to manage. Each service is bound to a Docker Image. This is an example of a service configuration file:
```json
//example.json
{
	"Name":"service1",
	"Type":"service1",
	"Image":"service1",
	"Constraints":{
		"MinActive":1,
		"MaxActive":5
	},
	"ContainerConfig":{
		//Docker configuration needed to start the container
	}
}
```
###### Run Gru
* Run/start the containers of the services you want to manage
* Run the Gru agent with the command `gru agent`. You can specify the logging level using the flag `-l`: e.g. `gru -l debug agent`
* Enjoy

