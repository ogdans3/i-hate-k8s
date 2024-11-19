package models

import (
	"strings"

	external_models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/external-models"
)

type EngineType uint8

const (
	nginx EngineType = iota
	haproxy
)

type LoadBalancer struct {
	Type    EngineType //Which underlying load balancer to use, e.g: nginx, haproxy //TODO: Find a better name
	Service *Service   //The service instance of the load balancer
}

func ParseLoadBalancer(loadbalancer bool, project Project) *LoadBalancer {
	//TODO: Fix hardcoding
	return &LoadBalancer{
		Type: nginx,
		Service: ParseService(&external_models.Service{
			ServiceName:   "loadbalancer",
			Image:         "nginx:1.27.2-alpine",
			Autoscale:     external_models.Autoscale{Initial: 1, Autoscale: false},
			ContainerName: "LB",
			Ports: []string{
				"80:80",
				"443:443",
			},
			Https: true,
			Www:   true,
			Volume: []*external_models.Volume{
				{
					Name:     strings.Join([]string{"hive", project.Project, "certs"}, "-"),
					Target:   "/etc/letsencrypt/live",
					Readonly: true,
				},
			},
		}, project),
	}
}
