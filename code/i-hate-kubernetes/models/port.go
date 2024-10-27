package models

import (
	"strings"
)

type Port struct {
	HostPort      string //The port on the host server
	ContainerPort string //The port on the container
	Protocol      string //The protocol for this port mapping (tcp, udp, etc)
}

func ParsePort(strPort string) Port {
	protocol := "tcp"
	hostPort := "80"
	containerPort := "80"

	parts := strings.Split(strPort, "/")
	// Check if there's a protocol specified
	if len(parts) > 1 {
		protocol = parts[1]
	}

	// Split hostPort and containerPort
	hostPortParts := strings.Split(parts[0], ":")
	if len(hostPortParts) == 1 {
		hostPort = hostPortParts[0]
		containerPort = hostPort
	} else {
		hostPort = hostPortParts[0]
		containerPort = hostPortParts[1]
	}
	return Port{
		HostPort:      hostPort,
		ContainerPort: containerPort,
		Protocol:      protocol,
	}
}

func ParsePorts(strPorts []string) []Port {
	ports := make([]Port, len(strPorts))
	for _, port := range strPorts {
		ports = append(ports, ParsePort(port))
	}
	return ports
}
