package model_actions

import (
	clientState "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/client-state"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/engine-interface/docker"
	engine_models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/engine-interface/engine-models"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/console"
	models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/internal-models"
)

type DeployContainerForService struct {
	DefaultActionMetadata
	Node    *models.Node
	Service *models.Service
	Project *models.Project
}

type RestartContainer struct {
	DefaultActionMetadata
	Node      *models.Node
	Container *engine_models.Container
}

type RemoveContainer struct {
	DefaultActionMetadata
	Node      *models.Node
	Container *engine_models.Container
}

// A container image has been updated. It is quite complex to update a container
// We need to first deploy the new container
// Then update the loadbalancer configuration to redirect traffic to the new containers
// Then gradually remove existing containers, gracefully
// This could take a long time, so it needs to be split into multiple actions.
// But we also need to keep it a single action, in case it fails and we need to rollback the changes?
type ContainerImageUpdated struct {
	CompositeAction
}

func CreateContainerImageUpdated(clientState *clientState.ClientState, node *models.Node, project *models.Project, service *models.Service) *ContainerImageUpdated {
	containers := clientState.GetContainersForService(service)
	actions := make([]Action, 0)
	for _, container := range containers {
		actions = append(actions,
			CreateDeployContainerForService(service, project),
			&RemoveContainer{
				Node:      node,
				Container: &container,
			},
		)
	}
	//TODO: Update loadbalancer (this is done automatically, so we probably dont need to do this)
	//TODO: Remove old containers gracefully
	//TODO: Rollback?
	//TODO: After the entire process is complete, we probably want to do a second pass over the containers to ensure that all containers have been removed
	// If we dont do a second pass, we could risk that the scheduler deploys more containers with the old image while we remove the old ones
	return &ContainerImageUpdated{
		CompositeAction{
			Actions: actions,
		},
	}
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

func (action *DeployContainerForService) Run() (ActionRunResult, error) {
	console.Log("Deploy container", action.Service.Id)
	if action.Service.Build {
		docker.BuildService(*action.Service, *action.Project)
	}
	_, err := docker.CreateContainerFromService(*action.Service, action.Project)
	return ActionRunResult{IsDone: true}, err
}

func (action *RestartContainer) Run() (ActionRunResult, error) {
	console.Log("Restart container", action.Container.Id)
	docker.StartContainer(action.Container.Id)
	return ActionRunResult{IsDone: true}, nil
}

func (action *RemoveContainer) Run() (ActionRunResult, error) {
	console.Log("Remove container", action.Container.Id)
	err := docker.StopAndRemoveContainer(*action.Container)
	if err != nil {
		return ActionRunResult{IsDone: false}, err
	}
	return ActionRunResult{IsDone: true}, nil
}

func (action *DeployContainerForService) Update(actions *[]Action, clientState *clientState.ClientState) (ActionUpdateResult, error) {
	return ActionUpdateResult{IsDone: true}, nil
}

func (action *RestartContainer) Update(actions *[]Action, clientState *clientState.ClientState) (ActionUpdateResult, error) {
	return ActionUpdateResult{IsDone: true}, nil
}

func (action *RemoveContainer) Update(actions *[]Action, clientState *clientState.ClientState) (ActionUpdateResult, error) {
	return ActionUpdateResult{IsDone: true}, nil
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
