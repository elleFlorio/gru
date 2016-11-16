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
Gru needs a running instance of an etcd server (https://github.com/coreos/etcd) for agents discovery and influxdb (https://influxdb.com/) to store metrics data.

## Documentation (Work in Progress)
These are the steps you need to follow to run the current version of Gru in your system.
Please remember that currently Gru is able only to autoscale your services containers and it's under active development.

### Requirements
Gru needs some external components and some environment variables to run. Gru has been developed and tested in linux, so the best thing is if you use the same environment.

#### External components
* Docker - https://www.docker.com/
* etcd - https://github.com/coreos/etcd
* InfluxDB (optional) - https://www.influxdata.com/time-series-platform/influxdb/

**Docker** is the obvious requirement and can be used also to run an instance of the Gru Agent (more on this later)

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
* Using the Docker image elleflorio/gru (elleflorio/gru:dev for the development branch image) **[suggested method]**.

In the latter case, Docker should be configured to listen on a specific port (usually 2375) for commands.
To configure properly Docker in linux, follow these instruction (credits to https://github.com/felixgborrego/docker-ui-chrome-app/wiki/linux).

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
/gru/ 
	|­­­---<cluster>/ 
		|­­­---uuid:string (Cluster ID) 
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