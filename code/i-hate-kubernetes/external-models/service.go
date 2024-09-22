package external_models

import (
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/definitions"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/util"
)

type Service struct {
	Image         string
	Build         string
	Dev           string
	Watch         string
	ContainerName string `yaml:"container_name"`
	FullName      string

	Www       bool
	Https     bool
	Ports     []string
	Autoscale Autoscale
}

func (service *Service) InsertDefaults() {
	var containerName = service.ContainerName
	if service.ContainerName == "" {
		containerName = service.Image
	}
	service.FullName = definitions.CONTAINER_KEY + "_" + containerName + "_" + util.RandStringBytesMaskImpr(3)
}
