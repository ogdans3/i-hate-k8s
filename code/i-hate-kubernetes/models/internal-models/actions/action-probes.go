package model_actions

import (
	"time"

	clientState "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/client-state"
	engine_models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/engine-interface/engine-models"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/probes"
	models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/internal-models"
)

type CheckReadinesProbe struct {
	Node      *models.Node
	Probe     *models.Probe
	Service   *models.Service
	Container *engine_models.Container
	Result    bool
}

func CreateReadinessProbe(node *models.Node, service *models.Service, container *engine_models.Container, probe *models.Probe, containerMetadata clientState.ContainerMetadata) *CheckReadinesProbe {
	now := time.Now().Unix()
	if now-containerMetadata.ProbesMetadata.Readiness.LastCheck < int64(probe.Interval) {
		return nil
	}

	return &CheckReadinesProbe{
		Node:      node,
		Probe:     probe,
		Service:   service,
		Container: container,
	}
}

func (action *CheckReadinesProbe) Run() error {
	result := probes.ProbeHttpGetMethod(action.Service, action.Container, *action.Probe)
	if result {
		action.Result = result
	}

	return nil
	//return err
}

func (action *CheckReadinesProbe) Update(clientState *clientState.ClientState) error {
	clientState.ContainerMetadata[action.Container.Id].ProbesMetadata.Readiness.LastCheck = time.Now().Unix()
	clientState.ContainerMetadata[action.Container.Id].ProbesMetadata.Readiness.ResultOfLastCheck = action.Result
	return nil
}

func (action *CheckReadinesProbe) Equals(otherAction Action) bool {
	other, ok := otherAction.(*CheckReadinesProbe)
	if !ok {
		return false
	}

	return action.Node == other.Node &&
		action.Container == other.Container &&
		action.Probe == other.Probe
}
