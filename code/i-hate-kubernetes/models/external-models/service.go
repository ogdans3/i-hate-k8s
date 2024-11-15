package external_models

import (
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/definitions"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/util"
)

type Service struct {
	ServiceName   string
	Image         string
	Build         bool //True if we should build this service using docker build. cicd must also be true
	Dev           string
	Watch         string
	ContainerName string `yaml:"container_name"`
	FullName      string `yaml:"full_name"`
	Pwd           string

	Www       bool
	Https     bool
	Ports     []string
	Autoscale Autoscale
	Probes    *Probes
}

func (service *Service) InsertDefaults(serviceName string) {
	var containerName = service.ContainerName
	if service.ContainerName == "" {
		containerName = service.Image
	}
	service.FullName = definitions.CONTAINER_KEY + "_" + containerName + "_" + util.RandStringBytesMaskImpr(3)
	service.ServiceName = serviceName
}
