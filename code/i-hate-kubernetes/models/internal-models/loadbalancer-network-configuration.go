package models

import (
	"strings"
)

type LoadBalancingMethod uint
type Time string

const (
	ROUND_ROBIN = iota
	LEAST_CONNECTIONS
	IP_HASH
	GENERIC_HASH
	RANDOM
)

type UpstreamServer struct {
	Server         string
	Weight         int
	MaxConnections int
	MaxFails       int
	FailTimeout    Time
	Backup         bool
	Down           bool
	Resolve        bool
	Route          string
	Service        string
	SlowStart      Time
	Drain          bool
}

type Upstream struct {
	Name    string
	Servers []UpstreamServer
	Method  LoadBalancingMethod
}

type ServerLocation struct {
	MatchModifier string //"~" for example for regex
	LocationMatch string
	ProxyPass     string
}

type Server struct {
	Location []ServerLocation
}

type Http struct {
	Upstream []Upstream
	Server   []Server
}

type LoadbalancerNetworkConfiguration struct {
	HttpBlocks []Http
}

const lineEnding string = "\r\n"
const tab string = "\t"

func (configuration *LoadbalancerNetworkConfiguration) ConfigurationToNginxFile() string {
	var builder strings.Builder
	for _, block := range configuration.HttpBlocks {
		builder.WriteString("error_log /var/log/nginx/error.log debug;" + lineEnding)
		builder.WriteString("events {}" + lineEnding)
		builder.WriteString("http {" + lineEnding)
		for _, upstreamBlock := range block.Upstream {
			builder.WriteString(tab + "upstream " + upstreamBlock.Name + " {" + lineEnding)
			for _, server := range upstreamBlock.Servers {
				//TODO: Add all the other shit, down, weight, route, etc
				builder.WriteString(tab + tab + "server " + server.Server + ";" + lineEnding)
			}
			builder.WriteString(tab + "}" + lineEnding)
		}
		for _, serverBlock := range block.Server {
			builder.WriteString(tab + "server {" + lineEnding)
			//TODO: Fix this listen statement
			builder.WriteString(tab + tab + "listen 80;" + lineEnding)
			for _, locationBlock := range serverBlock.Location {
				builder.WriteString(tab + tab + "location " + locationBlock.MatchModifier + " " + locationBlock.LocationMatch + " {" + lineEnding)
				builder.WriteString(tab + tab + tab + "proxy_pass http://" + locationBlock.ProxyPass + ";" + lineEnding)
				builder.WriteString(tab + tab + "}" + lineEnding)
			}
			builder.WriteString(tab + "}" + lineEnding)
		}
		builder.WriteString("}" + lineEnding)
	}
	return builder.String()
}
