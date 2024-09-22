package models

import (
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/external-models"
)

type Service struct {
	Image         string
	Build         string
	Dev           string
	Watch         string
	ContainerName string
	FullName      string

	Www       bool
	Https     bool
	Ports     []Port
	Autoscale Autoscale
}

func ParseService(service external_models.Service) Service {
	return Service{
		Image:         service.Image,
		Build:         service.Build,
		Dev:           service.Dev,
		Watch:         service.Watch,
		ContainerName: service.ContainerName,
		FullName:      service.FullName,

		Www:       service.Www,
		Https:     service.Https,
		Ports:     ParsePorts(service.Ports),
		Autoscale: ParseAutoscale(service.Autoscale),
	}
}

func ParseServices(unparsedServices map[string]external_models.Service) map[string]Service {
	services := map[string]Service{}
	for key, service := range unparsedServices {
		services[key] = ParseService(service)
	}
	return services
}
