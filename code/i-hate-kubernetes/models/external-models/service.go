package external_models

import (
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/definitions"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/util"
)

type Cicd struct {
	Id     string
	Url    string
	Branch string
}

type Service struct {
	Id            string
	ServiceName   string
	Image         string
	Build         bool   //True if we should build this service using docker build. cicd must also be true
	Directory     string //Relative directory to the project to use for docker contexts, also used to default cicd directory if not specified in cicd section
	Dockerfile    string //Dockerfile, relative to the specified Directory
	Domain        []string
	Path          []string
	Dev           string
	Watch         string
	ContainerName string `yaml:"container_name"`
	FullName      string `yaml:"full_name"`

	Www       bool
	Https     bool
	Ports     []string
	Autoscale Autoscale
	Probes    *Probes
	Cicd      *Cicd
}

func (service *Service) InsertDefaults(serviceName string) {
	var containerName = service.ContainerName
	if service.ContainerName == "" {
		containerName = service.Image
	}
	service.FullName = definitions.CONTAINER_KEY + "_" + containerName + "_" + util.RandStringBytesMaskImpr(3)
	service.ServiceName = serviceName
	service.Id = util.RandStringBytesMaskImpr(5)
	if service.Cicd != nil {
		service.Cicd.InsertDefaults()
	}
	if len(service.Path) == 0 {
		service.Path = []string{"/"}
	}
	service.Domain = append(service.Domain, "localhost", "127.0.0.1")
	if service.Dockerfile == "" {
		service.Dockerfile = "Dockerfile"
	}
}

func (cicd *Cicd) InsertDefaults() {
	cicd.Id = util.RandStringBytesMaskImpr(5)
	if cicd.Branch == "" {
		cicd.Branch = "master"
	}
	if cicd.Url == "" {
		panic("Cicd must provide a url")
	}
}
