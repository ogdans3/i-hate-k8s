package model_actions

import (
	clientState "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/client-state"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/engine-interface/docker"
	engine_models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/engine-interface/engine-models"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/console"
	models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/internal-models"
)

func GetDistinctActions(actions []Action) []Action {
	if len(actions) == 0 {
		return actions
	}
	result := []Action{}

	for _, action := range actions {
		seen := false
		for _, seenAction := range result {
			if action.Equals(seenAction) {
				seen = true
				break
			}
		}
		if !seen {
			result = append(result, action)
		}
	}

	return result
}

type Action interface {
	Run() error
	Update(clientState *clientState.ClientState) error
	Equals(action Action) bool
}

type DeployContainerForService struct {
	Node    *models.Node
	Service *models.Service
	Project *models.Project
}

type RestartContainer struct {
	Node      *models.Node
	Container *engine_models.Container
}

type RemoveContainer struct {
	Node      *models.Node
	Container *engine_models.Container
}

type UpdateLoadbalancer struct {
	Node                 *models.Node
	Container            *engine_models.Container
	NetworkConfiguration *models.LoadbalancerNetworkConfiguration
}

type CreateNetwork struct {
	Node      *models.Node    //Which node to create the network on
	Service   *models.Service //Which service is connected to this network
	NetworkId *string         //Returned from the engine daemon after the network has been created
}

type DeployNewNode struct {
}

func CreateDeployContainerForService(service *models.Service, project *models.Project) *DeployContainerForService {
	return &DeployContainerForService{
		Service: service,
		Project: project,
	}
}

func CreateRestartContainer(container *engine_models.Container, node *models.Node) *RestartContainer {
	return &RestartContainer{
		Container: container,
		Node:      node,
	}
}

func (action *DeployNewNode) CreateNetwork(clientState *clientState.ClientState) error {
	return nil
}

func (action *DeployContainerForService) Run() error {
	console.Log("Deploy container", action.Service.Id)
	if action.Service.Build {
		docker.BuildService(*action.Service, *action.Project)
	}
	_, err := docker.CreateContainerFromService(*action.Service, action.Project)
	return err
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
	err := docker.AddNewNginxConfigurationToContainer(action.NetworkConfiguration.ConfigurationToNginxFile(), *action.Container)
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

func (action *UpdateLoadbalancer) Update(clientState *clientState.ClientState) error {
	clientState.NetworkConfiguration = *action.NetworkConfiguration
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

func (action *DeployContainerForService) Equals(otherAction Action) bool {
	return false
	/*
		other, ok := otherAction.(*DeployContainerForService)
		if !ok {
			return false
		}

		return action.Node == other.Node &&
			action.Service == other.Service &&
			action.Project == other.Project
	*/
}

func (action *RestartContainer) Equals(otherAction Action) bool {
	other, ok := otherAction.(*RestartContainer)
	if !ok {
		return false
	}

	return action.Node == other.Node &&
		action.Container == other.Container
}

func (action *RemoveContainer) Equals(otherAction Action) bool {
	other, ok := otherAction.(*RemoveContainer)
	if !ok {
		return false
	}

	return action.Node == other.Node &&
		action.Container == other.Container
}

func (action *DeployNewNode) Equals(otherAction *DeployNewNode) bool {
	return true
}

func (action *UpdateLoadbalancer) Equals(otherAction Action) bool {
	other, ok := otherAction.(*UpdateLoadbalancer)
	if !ok {
		return false
	}

	return action.Node == other.Node &&
		action.Container == other.Container &&
		action.NetworkConfiguration == other.NetworkConfiguration
}

func (action *CreateNetwork) Equals(otherAction Action) bool {
	other, ok := otherAction.(*CreateNetwork)
	if !ok {
		return false
	}

	return action.Node == other.Node &&
		action.Service == other.Service &&
		action.NetworkId == other.NetworkId
}
