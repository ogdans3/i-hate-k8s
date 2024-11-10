package model_actions

import (
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/engine-interface/docker"
	engine_models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/engine-interface/engine-models"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/console"
	models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/internal-models"
)

type Action interface {
	Run()
}

type DeployContainerForService struct {
	Service models.Service
	Node    *models.Node
}

type RestartContainer struct {
	Container engine_models.Container
	Node      models.Node
}

type RemoveContainer struct {
	Container engine_models.Container
	Node      models.Node
}

type DeployNewNode struct {
}

func CreateDeployContainerForService(service models.Service) DeployContainerForService {
	return DeployContainerForService{
		Service: service,
	}
}

func CreateRestartContainer(container engine_models.Container, node models.Node) RestartContainer {
	return RestartContainer{
		Container: container,
		Node:      node,
	}
}

func (action DeployContainerForService) Run() {
	console.Log("Deploy container", action)
	docker.CreateContainerFromService(action.Service)
}

func (action RestartContainer) Run() {
	console.Log("Restart container", action)
	docker.StartContainer(action.Container.Id)
}

func (action RemoveContainer) Run() {
	console.Log("Remove container", action)
}

func (action DeployNewNode) Run() {
	console.Log("Deploy new node", action)
}
