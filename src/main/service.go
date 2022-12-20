package main

import (
	"container/list"
	"strings"
)

var (
	Device  = new(ServiceDefinition)
	Devices = list.New()

	endpoint = "opc.tcp://192.168.1.121:4840"
	ctx      = context.Background()
	nodeID   = flag.String("node", "i=85", "node id for the root node")
)

//Methods------------------------------------------------
// ------------------------------------------------------

func GetAllSensors() {
	Device = &ServiceDefinition{
		ServiceDefinition: "Sensor",
	}

		// Browse for nodes and compare them to the wished for axis
		nodeList, err := functions.Browse(ctx, c.Node(id), "", 0)
		if err != nil {
			log.Fatal(err)
		}

	if strings.Contains(, nodeID) {

	}
}

func GetAllActuators() {
	Device = &ServiceDefinition{
		ServiceDefinition: "Actuator",
	}
}
