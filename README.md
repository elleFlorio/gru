# Gru: automatic management of Docker containers
Gru brings self-adaptation to Docker containers: it manages Docker containers distributed in a set of nodes scaling them and actuating autonomic actions to keep everything up and running. Gru is designed to help developers build distributed applications based on microservices running in Docker containers.

Gru is part of my PhD research at Politecnico di Milano: "Decentralized Self-Adaptation in Large-Scale Distributed Systems" and it is currenlty under development.

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

## VERY IMPORTANT
This is a PhD project: this means that by now it is more like a "proof of concept".
I'm working to make Gru as close as possible to a "real" software, but please keep in mind that I'm the only person actively working on it and I'm learning the technology through the development process. So I am open to any suggestion/contribution.
If you like my work and want to help me in some way, you can contact me at **luca[dot]florio[at]polimi[dot]it**:
* if you are a student at Politecnico di Milano, we can discuss a thesis;
* if you are just curious about my work/want to give me a suggestion, I'm happy to have a chat;
* if you are a company and want to offer me an internship/job... Well, you make me very happy! :-)
