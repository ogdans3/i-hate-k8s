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
	Port           string
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
	ServerName    string
	ProxyPass     string
}

type Server struct {
	Location   []ServerLocation
	ServerName []string
}

type Http struct {
	Upstream []Upstream
	Server   []Server
}

type LoadbalancerNetworkConfiguration struct {
	ContainerIdOfLoadbalancerThatHasThisConfig *string
	HttpBlock                                  Http
}

const lineEnding string = "\r\n"
const tab string = "\t"

func (configuration *LoadbalancerNetworkConfiguration) ConfigurationToNginxFile() string {
	var builder strings.Builder
	//builder.WriteString("error_log /var/log/nginx/error.log debug;" + lineEnding)
	builder.WriteString("events {}" + lineEnding)
	builder.WriteString("http {" + lineEnding)
	builder.WriteString(tab + `log_format custom '$remote_addr -> $upstream_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent" "$http_x_forwarded_for"';` + lineEnding)
	builder.WriteString(tab + "access_log /var/log/nginx/access.log custom;" + lineEnding)
	builder.WriteString(tab + "error_log /var/log/nginx/error.log warn;" + lineEnding)
	for _, upstreamBlock := range configuration.HttpBlock.Upstream {
		builder.WriteString(tab + "upstream " + upstreamBlock.Name + " {" + lineEnding)
		for _, server := range upstreamBlock.Servers {
			//TODO: Add all the other shit, down, weight, route, etc
			builder.WriteString(tab + tab + "server " + server.Server + ":" + server.Port + ";" + lineEnding)
			//builder.WriteString(tab + tab + "server " + server.Server + ";" + lineEnding)
		}
		builder.WriteString(tab + "}" + lineEnding)
	}
	builder.WriteString(tab + "server {" + lineEnding)
	//TODO: Fix this listen statement
	builder.WriteString(tab + tab + "listen 80;" + lineEnding)
	for _, serverBlock := range configuration.HttpBlock.Server {
		builder.WriteString(tab + tab + "server_name ")
		for _, serverName := range serverBlock.ServerName {
			if serverName != "" {
				builder.WriteString(serverName + " ")
			}
		}
		builder.WriteString(" ;" + lineEnding)
		for _, locationBlock := range serverBlock.Location {
			builder.WriteString(tab + tab + "location " + locationBlock.MatchModifier + " " + locationBlock.LocationMatch + " {" + lineEnding)
			builder.WriteString(tab + tab + tab + "proxy_pass http://" + locationBlock.ProxyPass + ";" + lineEnding)
			builder.WriteString(tab + tab + "}" + lineEnding)
		}
	}
	builder.WriteString(tab + "}" + lineEnding)
	builder.WriteString("}" + lineEnding)
	return builder.String()
}
