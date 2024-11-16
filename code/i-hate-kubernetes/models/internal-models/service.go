package models

import (
	"path/filepath"

	external_models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/external-models"
)

type Service struct {
	//TODO: Rename the ID?
	Id            string //An ID generated by this program in order to attach services to containers, includes the project id
	ServiceId     string //An ID generated by this program to uniquely identify this service
	Directory     string //Directory of this service, used for build commands that must happen in the folder. Relative path to the project pwd
	Dockerfile    string //The filename of the dockerfile
	ServiceName   string //The name of this service, specified by the key in yaml
	Image         string //The docker image to use
	Build         bool   //True if we should run docker build to build image for this service
	Dev           string //A dev command?
	Watch         string //No idea?
	ContainerName string //The name for this container, will appear in docker ps
	FullName      string //The name for this container, will appear in docker ps

	Www       bool      //Should requests be redirected from example.com to www.example.com
	Https     bool      //Should https be used
	Ports     []Port    //A list of port mappings
	Autoscale Autoscale //Autoscaling settings for this pod

	Network *Network
	Probes  *Probes
}

func ParseService(service *external_models.Service, project Project) *Service {
	projectId := *project.GetId()
	return &Service{
		Id:            projectId + "-" + service.ServiceName,
		ServiceId:     service.Id,
		Directory:     filepath.Join(project.Pwd, service.Directory),
		Dockerfile:    "Dockerfile", //TODO: Dont hardcode
		ServiceName:   service.ServiceName,
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
		Probes:    ParseProbes(service.Probes),

		Network: &Network{
			Name: projectId,
		},
	}
}

func ParseServices(unparsedServices map[string]*external_models.Service, project Project) map[string]*Service {
	services := map[string]*Service{}
	for key, service := range unparsedServices {
		services[key] = ParseService(service, project)
	}
	return services
}

func (service *Service) GetId() *string {
	if service == nil {
		return nil
	}

	return &service.Id
}
