package model_actions

import (
	clientState "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/client-state"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/engine-interface/docker"
	engine_models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/engine-interface/engine-models"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/console"
	models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/internal-models"
)

type Action interface {
	Run() error
	Update(clientState *clientState.ClientState) error
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

type CreateNetwork struct {
	Node      models.Node     //Which node to create the network on
	Service   *models.Service //Which service is connected to this network
	NetworkId *string         //Returned from the engine daemon after the network has been created
}

type DeployNewNode struct {
}

func CreateDeployContainerForService(service models.Service) *DeployContainerForService {
	return &DeployContainerForService{
		Service: service,
	}
}

func CreateRestartContainer(container engine_models.Container, node models.Node) *RestartContainer {
	return &RestartContainer{
		Container: container,
		Node:      node,
	}
}

func (action *DeployContainerForService) Run() error {
	console.Log("Deploy container", action.Service.Id)
	docker.CreateContainerFromService(action.Service)
	return nil
}

func (action *RestartContainer) Run() error {
	console.Log("Restart container", action.Container.Id)
	docker.StartContainer(action.Container.Id)
	return nil
}

func (action *RemoveContainer) Run() error {
	console.Log("Remove container", action.Container.Id)
	return nil
}

func (action *DeployNewNode) Run() error {
	console.Log("Deploy new node", action)
	return nil
}

func (action *CreateNetwork) Run() error {
	console.Log("Create network")
	networkId, err := docker.CreateNetwork(*action.Service.Network)
	action.NetworkId = networkId
	return err
}

func (action *UpdateLoadbalancer) Run() error {
	console.Log("Update loadbalancer", action.Container.Id)
	err := docker.AddNewNginxConfigurationToContainer(action.NetworkConfiguration.ConfigurationToNginxFile(), action.Container)
	return err
}

func (action *DeployContainerForService) Update(clientState *clientState.ClientState) error {
	return nil
}

func (action *RestartContainer) Update(clientState *clientState.ClientState) error {
	return nil
}

func (action *RemoveContainer) Update(clientState *clientState.ClientState) error {
	return nil
}

func (action *DeployNewNode) Update(clientState *clientState.ClientState) error {
	return nil
}

func (action *DeployNewNode) CreateNetwork(clientState *clientState.ClientState) error {
	return nil
}

func (action *UpdateLoadbalancer) Update(clientState *clientState.ClientState) error {
	clientState.NetworkConfiguration = action.NetworkConfiguration
	return nil
}

func (action *CreateNetwork) Update(clientState *clientState.ClientState) error {
	if action.Service.Network.GetName() == nil {
		//TODO: Return error?
		return nil
	}
	//TODO: Do we need to retain the network id?
	console.Log("Network created: ", action.NetworkId)
	if clientState.EngineNetworkToService[*action.Service.Network.GetName()] == nil {
		clientState.EngineNetworkToService[*action.Service.Network.GetName()] = make([]models.Service, 0)
	}
	clientState.EngineNetworkToService[*action.Service.Network.GetName()] = append(clientState.EngineNetworkToService[*action.Service.Network.GetName()], *action.Service)
	return nil
}
