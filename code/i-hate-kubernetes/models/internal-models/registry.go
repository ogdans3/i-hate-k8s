package models

import external_models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/external-models"

type Registry struct {
	Service *Service //The service instance of the load balancer
}

func ParseRegistry(registry bool, project Project) *Registry {
	return &Registry{
		Service: ParseService(&external_models.Service{
			ServiceName:   "registry",
			ContainerName: "registry",
			Image:         "registry:2.8",
			Ports: []string{
				"5000:5000", //TODO This should probably not be exposed to the internet
			},
			Autoscale: external_models.Autoscale{Initial: 1, Autoscale: false},
			Probes: &external_models.Probes{
				Ready: "/v2/_catalog",
			},
		}, project),
	}
}
