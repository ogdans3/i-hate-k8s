package models

import (
	external_models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/external-models"
)

type DeploymentStrategy int

const (
	Gracefully DeploymentStrategy = iota
	Forcefully
)

type Cicd struct {
	Id                      string
	Url                     string
	Remote                  string
	Branch                  string
	Directory               string
	Service                 *Service
	UpdateContainers        bool
	UpdateContainerStrategy DeploymentStrategy
}

func ParseCicds(services map[string]*Service, externalServices map[string]*external_models.Service, wd string) []Cicd {
	parsed := []Cicd{}

	for _, externalService := range externalServices {
		if externalService.Cicd == nil {
			continue
		}
		parsed = append(parsed,
			Cicd{
				Id:        externalService.Cicd.Id,
				Url:       externalService.Cicd.Url,
				Branch:    externalService.Cicd.Branch,
				Remote:    "origin",
				Directory: wd,
				Service:   services[externalService.ServiceName],
			},
		)
	}
	return parsed
}

func ParseAutoupdate(autoupdate bool, wd string) *Cicd {
	if !autoupdate {
		return nil
	}
	validUrls := []string{"https://github.com/ogdans3/i-hate-k8s.git", "git@github.com:ogdans3/i-hate-k8s.git"}
	return &Cicd{
		Url:                     validUrls[0],
		Branch:                  "master",
		Directory:               wd,
		UpdateContainers:        true,
		UpdateContainerStrategy: Gracefully,
	}
}
