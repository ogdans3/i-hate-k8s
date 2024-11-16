package model_actions

import (
	clientState "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/client-state"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/console"
)

type DeployNewNode struct {
	DefaultActionMetadata
}

func (action *DeployNewNode) Run() (ActionRunResult, error) {
	console.Log("Deploy new node", action)
	return ActionRunResult{IsDone: true}, nil
}

func (action *DeployNewNode) Update(actions *[]Action, clientState *clientState.ClientState) (ActionUpdateResult, error) {
	return ActionUpdateResult{IsDone: true}, nil
}

func (action *DeployNewNode) Equals(otherAction *DeployNewNode) bool {
	return true
}
