package model_actions

import (
	clientState "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/client-state"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/util"
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
	IsDone      bool
	NeedsUpdate bool
}

type ActionUpdateResult struct {
	IsDone bool
}

type ActionMetadata interface {
	GetRetries() uint8
	IncreaseRetries()
	ResetRetries()
}

type DefaultActionMetadata struct {
	Id      *string
	Retries uint8
}

type Action interface {
	Run() (ActionRunResult, error)
	Update(actions *[]Action, clientState *clientState.ClientState) (ActionUpdateResult, error)
	Equals(action Action) bool
	GetMetadata() ActionMetadata
}

func (metadata *DefaultActionMetadata) GetMetadata() ActionMetadata {
	return metadata
}

func (action *DefaultActionMetadata) GetRetries() uint8 {
	return action.Retries
}

func (action *DefaultActionMetadata) IncreaseRetries() {
	action.Retries++
}

func (action *DefaultActionMetadata) ResetRetries() {
	action.Retries = 0
}

// TODO: This is NOT a good way to initialise the id, but whatever for now.
func (action *DefaultActionMetadata) GetId() string {
	if action.Id == nil {
		id := util.RandStringBytesMaskImpr(5)
		action.Id = &id
	}
	return *action.Id
}
