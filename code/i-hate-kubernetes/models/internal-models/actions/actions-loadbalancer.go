package model_actions

import (
	"strings"

	clientState "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/client-state"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/engine-interface/docker"
	engine_models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/engine-interface/engine-models"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/console"
	models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/internal-models"
)

type UpdateLoadbalancer struct {
	DefaultActionMetadata
	Node                 *models.Node
	Container            *engine_models.Container
	NetworkConfiguration *models.LoadbalancerNetworkConfiguration
}

type CreateNetwork struct {
	DefaultActionMetadata
	Node      *models.Node    //Which node to create the network on
	Service   *models.Service //Which service is connected to this network
	NetworkId *string         //Returned from the engine daemon after the network has been created
}

func CreateLoadbalancerAction(node *models.Node, loadBalancerContainer *engine_models.Container, project *models.Project, containers []engine_models.Container) *UpdateLoadbalancer {
	newConfig := models.LoadbalancerNetworkConfiguration{
		ContainerIdOfLoadbalancerThatHasThisConfig: &loadBalancerContainer.Id,
	}

	services := project.GetLoadbalancedServices()

	serverBlocks := make([]models.Server, 0)
	upstreamBlocks := make([]models.Upstream, 0)
	for _, service := range services {
		upstreamBlock := models.Upstream{
			Name:    service.ServiceName,
			Servers: []models.UpstreamServer{},
		}
		serverBlock := models.Server{
			Location:   []models.ServerLocation{},
			ServerName: service.Domain,
		}
		for _, path := range service.Path {
			//TODO: Loop over path/domain for the service and insert into the location block
			serverBlock.Location = append(serverBlock.Location, models.ServerLocation{
				MatchModifier: "",
				LocationMatch: path,
				ProxyPass:     upstreamBlock.Name,
			})
		}
		for _, container := range containers {
			for _, name := range container.Names {
				if strings.Contains(name, service.Id) {
					//TODO: Handle multiple ports better? Allow the user to select the port atleast
					for _, port := range service.Ports {
						if container.Ip != nil {
							upstreamBlock.Servers = append(upstreamBlock.Servers, models.UpstreamServer{
								Server: *container.GetIp(),
								Port:   port.ContainerPort,
							})
						}
					}
				}
			}
		}
		serverBlocks = append(serverBlocks, serverBlock)
		upstreamBlocks = append(upstreamBlocks, upstreamBlock)
	}
	newConfig.HttpBlock = models.Http{
		Upstream: upstreamBlocks,
		Server:   serverBlocks,
	}

	return &UpdateLoadbalancer{
		Node:                 node,
		Container:            loadBalancerContainer,
		NetworkConfiguration: &newConfig,
	}

}

func (action *DeployNewNode) CreateNetwork(clientState *clientState.ClientState) error {
	return nil
}

func (action *CreateNetwork) Run() (ActionRunResult, error) {
	console.InfoLog.Log("Create network")
	networkId, err := docker.CreateNetwork(*action.Service.Network)
	action.NetworkId = networkId
	return ActionRunResult{IsDone: true}, err
}

func (action *CreateNetwork) Update(actions *[]Action, clientState *clientState.ClientState) (ActionUpdateResult, error) {
	if action.Service.Network.GetName() == nil {
		//TODO: Return error?
		return ActionUpdateResult{IsDone: true}, nil
	}
	//TODO: Do we need to retain the network id?
	console.InfoLog.Log("Network created: ", action.NetworkId)
	if clientState.EngineNetworkToService[*action.Service.Network.GetName()] == nil {
		clientState.EngineNetworkToService[*action.Service.Network.GetName()] = make([]models.Service, 0)
	}
	clientState.EngineNetworkToService[*action.Service.Network.GetName()] = append(clientState.EngineNetworkToService[*action.Service.Network.GetName()], *action.Service)
	return ActionUpdateResult{IsDone: true}, nil
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

func (action *UpdateLoadbalancer) Run() (ActionRunResult, error) {
	console.InfoLog.Log("Update loadbalancer", action.Container.Id)
	err := docker.AddNewNginxConfigurationToContainer(action.NetworkConfiguration.ConfigurationToNginxFile(), *action.Container)
	return ActionRunResult{IsDone: true}, err
}

func (action *UpdateLoadbalancer) Update(actions *[]Action, clientState *clientState.ClientState) (ActionUpdateResult, error) {
	clientState.NetworkConfiguration = *action.NetworkConfiguration
	return ActionUpdateResult{IsDone: true}, nil
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
