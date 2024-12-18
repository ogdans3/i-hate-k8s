package clientState

import (
	"strings"

	engine_models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/engine-interface/engine-models"
	models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/internal-models"
)

func (cs *ClientState) GetSingleContainerForService(service *models.Service) *engine_models.Container {
	for _, container := range cs.Containers {
		for _, name := range container.Names {
			if strings.Contains(name, service.Id) {
				return &container
			}
		}
	}
	return nil
}

func (cs *ClientState) GetContainersForService(service *models.Service) []engine_models.Container {
	containers := make([]engine_models.Container, 0)
	for _, container := range cs.Containers {
		for _, name := range container.Names {
			if strings.Contains(name, service.Id) {
				containers = append(containers, container)
			}
		}
	}
	return containers
}

func (cs *ClientState) GetContainerFromContainerId(containerId string) *engine_models.Container {
	for _, ctr := range cs.Containers {
		if ctr.Id == containerId {
			return &ctr
		}
	}
	return nil
}

func (cs *ClientState) GetVolumeFromName(volumeName string) *engine_models.Volume {
	for _, vol := range cs.Volumes {
		if vol.Name == volumeName {
			return &vol
		}
	}
	return nil
}
