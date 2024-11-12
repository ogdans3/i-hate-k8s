package model_actions

import (
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/engine-interface/docker"
	engine_models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/engine-interface/engine-models"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/console"
	models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/internal-models"
)

type Action interface {
	Run() error
}

type DeployContainerForService struct {
	Node    *models.Node
	Service models.Service
}

type RestartContainer struct {
	Node      models.Node
	Container engine_models.Container
}

type RemoveContainer struct {
	Node      models.Node
	Container engine_models.Container
}

type UpdateLoadbalancer struct {
	Node                 models.Node
	Container            engine_models.Container
	NetworkConfiguration models.LoadbalancerNetworkConfiguration
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

func (action DeployContainerForService) Run() error {
	console.Log("Deploy container", action.Service.Id)
	docker.CreateContainerFromService(action.Service)
	return nil
}

func (action RestartContainer) Run() error {
	console.Log("Restart container", action.Container.Id)
	docker.StartContainer(action.Container.Id)
	return nil
}

func (action RemoveContainer) Run() error {
	console.Log("Remove container", action.Container.Id)
	return nil
}

func (action DeployNewNode) Run() error {
	console.Log("Deploy new node", action)
	return nil
}

func (action UpdateLoadbalancer) Run() error {
	console.Log("Update loadbalancer", action.Container.Id)
	err := docker.AddNewNginxConfigurationToContainer(action.NetworkConfiguration.ConfigurationToNginxFile(), action.Container)
	return err
}
