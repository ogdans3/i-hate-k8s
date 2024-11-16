package model_actions

import (
	clientState "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/client-state"
)

func GetDistinctActions(actions *[]Action) {
	if len(*actions) == 0 {
		return
	}
	result := []Action{}

	for _, action := range *actions {
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

	*actions = result
}

type ActionRunResult struct {
	IsDone bool
}

type ActionUpdateResult struct {
	IsDone bool
}

type ActionMetadata struct {
	Retries uint8
}

type DefaultActionMetadata struct {
	Metadata ActionMetadata
}

type Action interface {
	Run() (ActionRunResult, error)
	Update(actions *[]Action, clientState *clientState.ClientState) (ActionUpdateResult, error)
	Equals(action Action) bool
	GetMetadata() ActionMetadata
}

func (action DefaultActionMetadata) GetMetadata() ActionMetadata {
	return action.Metadata
}
