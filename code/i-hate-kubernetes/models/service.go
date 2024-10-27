package models

import (
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/external-models"
)

type Service struct {
	Image         string //The docker image to use
	Build         string //A build command which will be used during development
	Dev           string //A dev command?
	Watch         string //No idea?
	ContainerName string //The name for this container, will appear in docker ps
	FullName      string //The name for this container, will appear in docker ps

	Www       bool      //Should requests be redirected from example.com to www.example.com
	Https     bool      //Should https be used
	Ports     []Port    //A list of port mappings
	Autoscale Autoscale //Autoscaling settings for this pod
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
