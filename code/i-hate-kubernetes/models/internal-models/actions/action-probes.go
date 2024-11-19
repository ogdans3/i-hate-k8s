package model_actions

import (
	"time"

	clientState "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/client-state"
	engine_models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/engine-interface/engine-models"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/probes"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/console"
	models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/internal-models"
)

type CreateReadinessProbeFromDeployAction struct {
	DefaultActionMetadata
	DeployAction *DeployContainerForService
	Probe        *models.Probe
	ProbeAction  *CheckReadinessProbe
}

type CheckLivenessProbe struct {
	DefaultActionMetadata
	Node      *models.Node
	Probe     *models.Probe
	Service   *models.Service
	Container *engine_models.Container
	Result    bool
}

type CheckReadinessProbe struct {
	DefaultActionMetadata
	Node        *models.Node
	Probe       *models.Probe
	Service     *models.Service
	Container   *engine_models.Container
	MustSucceed bool
	Result      bool
}

func CreateReadinessProbe(node *models.Node, service *models.Service, container *engine_models.Container, probe *models.Probe, containerMetadata clientState.ContainerMetadata) *CheckReadinessProbe {
	now := time.Now().Unix()
	if now-containerMetadata.ProbesMetadata.Readiness.LastCheck < int64(probe.Interval) {
		return nil
	}

	return &CheckReadinessProbe{
		Node:        node,
		Probe:       probe,
		Service:     service,
		Container:   container,
		MustSucceed: false,
	}
}

func CreateLivenessProbe(node *models.Node, service *models.Service, container *engine_models.Container, probe *models.Probe, containerMetadata clientState.ContainerMetadata) *CheckLivenessProbe {
	now := time.Now().Unix()
	if now-containerMetadata.ProbesMetadata.Readiness.LastCheck < int64(probe.Interval) {
		return nil
	}

	return &CheckLivenessProbe{
		Node:      node,
		Probe:     probe,
		Service:   service,
		Container: container,
	}
}

func (action *CheckReadinessProbe) Run() (ActionRunResult, error) {
	console.InfoLog.Debug("Check readiness: ", action.GetId())
	result := probes.ProbeHttpGetMethod(action.Service, action.Container, *action.Probe)
	if result {
		action.Result = result
	}
	if action.MustSucceed && !action.Result {
		return ActionRunResult{IsDone: false}, nil
	}
	return ActionRunResult{IsDone: true}, nil
}

func (action *CheckReadinessProbe) Update(actions *[]Action, clientState *clientState.ClientState) (ActionUpdateResult, error) {
	clientState.ContainerMetadata[action.Container.Id].ProbesMetadata.Readiness.LastCheck = time.Now().Unix()
	clientState.ContainerMetadata[action.Container.Id].ProbesMetadata.Readiness.ResultOfLastCheck = action.Result
	return ActionUpdateResult{IsDone: true}, nil
}

func (action *CheckReadinessProbe) Equals(otherAction Action) bool {
	other, ok := otherAction.(*CheckReadinessProbe)
	if !ok {
		return false
	}

	return action.Node == other.Node &&
		action.Container == other.Container &&
		action.Probe == other.Probe
}

func (action *CheckLivenessProbe) Run() (ActionRunResult, error) {
	result := probes.ProbeHttpGetMethod(action.Service, action.Container, *action.Probe)
	if result {
		action.Result = result
	}

	return ActionRunResult{IsDone: true}, nil
}

func (action *CheckLivenessProbe) Update(actions *[]Action, clientState *clientState.ClientState) (ActionUpdateResult, error) {
	clientState.ContainerMetadata[action.Container.Id].ProbesMetadata.Liveness.LastCheck = time.Now().Unix()
	clientState.ContainerMetadata[action.Container.Id].ProbesMetadata.Liveness.ResultOfLastCheck = action.Result
	return ActionUpdateResult{IsDone: true}, nil
}

func (action *CheckLivenessProbe) Equals(otherAction Action) bool {
	other, ok := otherAction.(*CheckLivenessProbe)
	if !ok {
		return false
	}

	return action.Node == other.Node &&
		action.Container == other.Container &&
		action.Probe == other.Probe
}

func (action *CreateReadinessProbeFromDeployAction) Run() (ActionRunResult, error) {
	return ActionRunResult{IsDone: true}, nil
}

func (action *CreateReadinessProbeFromDeployAction) Update(actions *[]Action, clientState *clientState.ClientState) (ActionUpdateResult, error) {
	ctr := clientState.GetContainerFromContainerId(*action.DeployAction.ContainerId)
	action.ProbeAction = CreateReadinessProbe(action.DeployAction.Node, action.DeployAction.Service, ctr, action.Probe, clientState.ContainerMetadata[ctr.Id])
	return ActionUpdateResult{IsDone: true}, nil
}

func (action *CreateReadinessProbeFromDeployAction) Equals(otherAction Action) bool {
	other, ok := otherAction.(*CreateReadinessProbeFromDeployAction)
	if !ok {
		return false
	}

	return action.DeployAction == other.DeployAction &&
		action.Probe == other.Probe
}
