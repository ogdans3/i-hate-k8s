package external_models

import (
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/definitions"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/util"
)

type Service struct {
	ServiceName   string
	Image         string
	Build         string
	Dev           string
	Watch         string
	ContainerName string `yaml:"container_name"`
	FullName      string `yaml:"full_name"`

	Www       bool
	Https     bool
	Ports     []string
	Autoscale Autoscale
}

func (service *Service) InsertDefaults(serviceName string) {
	var containerName = service.ContainerName
	if service.ContainerName == "" {
		containerName = service.Image
	}
	service.FullName = definitions.CONTAINER_KEY + "_" + containerName + "_" + util.RandStringBytesMaskImpr(3)
	service.ServiceName = serviceName
}
