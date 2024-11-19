package model_actions

import (
	clientState "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/client-state"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/engine-interface/docker"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/console"
	models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/internal-models"
)

type ImageBuild struct {
	DefaultActionMetadata
	Node    *models.Node
	Service *models.Service
	Project *models.Project
}

func (action *ImageBuild) Run() (ActionRunResult, error) {
	console.InfoLog.Log("Build image: ", action.Service.Image)

	docker.BuildService(*action.Service, *action.Project)
	return ActionRunResult{IsDone: true}, nil
}

func (action *ImageBuild) Update(actions *[]Action, clientState *clientState.ClientState) (ActionUpdateResult, error) {
	return ActionUpdateResult{IsDone: true}, nil
}

func (action *ImageBuild) Equals(otherAction Action) bool {
	other, ok := otherAction.(*ImageBuild)
	if !ok {
		return false
	}

	return action.Node == other.Node &&
		action.Service == other.Service
}
