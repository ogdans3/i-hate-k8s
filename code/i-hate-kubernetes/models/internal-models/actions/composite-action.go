package model_actions

import (
	clientState "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/client-state"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/console"
)

type CompositeActionState int

const (
	Action_stage CompositeActionState = iota
	Finally_stage
)

type ActionBloc struct {
	Actions                   []Action
	indexOfNextActionToRun    uint8
	indexOfNextActionToUpdate *uint8
}

type CompositeAction struct {
	DefaultActionMetadata
	Actions *ActionBloc
	Finally *ActionBloc
	State   CompositeActionState
}

func (a *CompositeAction) Run() (ActionRunResult, error) {
	if a.State == Action_stage {
		return RunActionBlock(a.Actions)
	} else if a.State == Finally_stage {
		return RunActionBlock(a.Finally)
	}
	panic("Eieieiei")
}

func (a *CompositeAction) Update(actions *[]Action, clientState *clientState.ClientState) (ActionUpdateResult, error) {
	if a.State == Action_stage {
		return UpdateActionBlock(a.Actions, actions, clientState)
	} else if a.State == Finally_stage {
		return UpdateActionBlock(a.Finally, actions, clientState)
	}
	panic("Eieieiei")
}

func RunActionBlock(a *ActionBloc) (ActionRunResult, error) {
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

func UpdateActionBlock(a *ActionBloc, actions *[]Action, clientState *clientState.ClientState) (ActionUpdateResult, error) {
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
		var value uint8 = *a.indexOfNextActionToUpdate + 1
		a.indexOfNextActionToUpdate = &value
	}
	if len(a.Actions) < int(*a.indexOfNextActionToUpdate) {
		return ActionUpdateResult{IsDone: false}, nil
	}
	return ActionUpdateResult{IsDone: true}, nil
}

func (action *CompositeAction) Equals(otherAction Action) bool {
	return false
}
