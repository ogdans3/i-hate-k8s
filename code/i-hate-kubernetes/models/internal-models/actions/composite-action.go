package model_actions

import (
	clientState "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/client-state"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/console"
)

type CompositeAction struct {
	DefaultActionMetadata
	Actions                   []Action
	indexOfNextActionToRun    uint8
	indexOfNextActionToUpdate *uint8
}

func (a *CompositeAction) Run() (ActionRunResult, error) {
	if len(a.Actions) <= int(a.indexOfNextActionToRun) {
		return ActionRunResult{IsDone: true}, nil
	}

	action := a.Actions[a.indexOfNextActionToRun]
	result, err := action.Run()
	if err != nil {
		console.InfoLog.Error("Action inside composite action failed to run: ", err)
		return ActionRunResult{IsDone: false}, err
	}
	if result.IsDone {
		a.indexOfNextActionToRun++
	}
	if a.indexOfNextActionToUpdate == nil {
		var value uint8 = 0
		a.indexOfNextActionToUpdate = &value
	}
	if len(a.Actions) > int(a.indexOfNextActionToRun) {
		return ActionRunResult{IsDone: false}, nil
	}
	return ActionRunResult{IsDone: true}, nil
}

func (a *CompositeAction) Update(actions *[]Action, clientState *clientState.ClientState) (ActionUpdateResult, error) {
	if a.indexOfNextActionToUpdate == nil {
		//No actions have run yet, so we do not update the state
		return ActionUpdateResult{IsDone: false}, nil
	}

	action := a.Actions[*a.indexOfNextActionToUpdate]
	result, err := action.Update(actions, clientState)
	if err != nil {
		console.InfoLog.Error("Action inside composite action failed to update the state: ", err)
		return ActionUpdateResult{IsDone: false}, err
	}
	if result.IsDone {
		a.indexOfNextActionToRun++
	}
	if len(a.Actions) < int(a.indexOfNextActionToRun) {
		return ActionUpdateResult{IsDone: false}, nil
	}
	return ActionUpdateResult{IsDone: true}, nil
}

func (action *CompositeAction) Equals(otherAction Action) bool {
	return false
}
