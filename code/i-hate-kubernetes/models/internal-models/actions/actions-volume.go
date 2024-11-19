package model_actions

import (
	clientState "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/client-state"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/engine-interface/docker"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/console"
	models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/internal-models"
)

type CreateVolume struct {
	DefaultActionMetadata
	Id     string
	Node   *models.Node
	Volume *models.Volume
}

func (action *CreateVolume) Run() (ActionRunResult, error) {
	console.Log("Action to create volume: ", action.Id, action.Volume.Name)
	_, err := docker.CreateVolume(action.Volume.Name)
	if err != nil {
		console.InfoLog.Error("Action to create volume failed: ", action.Id)
		return ActionRunResult{IsDone: false}, err
	}

	return ActionRunResult{IsDone: true}, nil
}

func (action *CreateVolume) Update(actions *[]Action, clientState *clientState.ClientState) (ActionUpdateResult, error) {
	return ActionUpdateResult{IsDone: true}, nil
}

func (action *CreateVolume) Equals(otherAction Action) bool {
	other, ok := otherAction.(*CreateVolume)
	if !ok {
		return false
	}

	return action.Node == other.Node &&
		action.Volume == other.Volume
}
