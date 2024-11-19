package model_actions

import (
	"errors"
	"strings"

	clientState "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/client-state"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/engine-interface/docker"
	engine_models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/engine-interface/engine-models"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/console"
	models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/internal-models"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/util"
)

type certificateState int

const (
	deploy_container certificateState = iota
	container_deployed
)

type CertificateMasterJob struct {
	CompositeAction
	id           string
	handler      *models.CertificateHandler
	state        certificateState
	deployAction *DeployContainerForService
}

type IssueTlsCertificate struct {
	DefaultActionMetadata
	CertificationBlock *models.CertificateBlock
	Container          *engine_models.Container
}

type WaitForTlsCertificate struct {
	DefaultActionMetadata
	CertificationBlock *models.CertificateBlock
	Container          *engine_models.Container
}

func CreateCertificateMasterJob(node *models.Node, certificationHandler *models.CertificateHandler, project *models.Project) *CertificateMasterJob {
	service := certificationHandler.ServiceJob.Service
	deployAction := CreateDeployContainerForService(service.ServiceId, service, project)
	return &CertificateMasterJob{
		id:           util.RandStringBytesMaskImpr(5),
		handler:      certificationHandler,
		state:        deploy_container,
		deployAction: deployAction,
		CompositeAction: CompositeAction{
			Actions: &ActionBloc{
				Actions: []Action{
					&ImageBuild{
						Node:    node,
						Service: service,
						Project: project,
					},
					deployAction,
				},
			},
			Finally: &ActionBloc{},
		},
	}
}

func (action *CertificateMasterJob) Run() (ActionRunResult, error) {
	console.InfoLog.Log("Master certificate job: ", action.id)
	result, err := action.CompositeAction.Run()
	if action.state == deploy_container {
		return ActionRunResult{IsDone: false, NeedsUpdate: true}, err
	}

	//TODO: Oof, this should be handled by the client kind of. It should call a method on this action when it is being removed from consideration
	if action.GetMetadata().GetRetries() > 1 || (result.IsDone && action.State == Action_stage) {
		action.State = Finally_stage
		action.GetMetadata().ResetRetries()
		return ActionRunResult{IsDone: false, NeedsUpdate: false}, err
	}
	return result, err
}

func (action *CertificateMasterJob) Update(actions *[]Action, clientState *clientState.ClientState) (ActionUpdateResult, error) {
	if action.state == deploy_container {
		//Since we are running the update, we should expect that the deploy action has completed
		//We are not done yet, because we have to run more actions

		//Update the image build and the container create
		if len(action.Actions.Actions) > int(*action.Actions.indexOfNextActionToUpdate) {
			action.CompositeAction.Update(actions, clientState)
		}
		if *action.Actions.indexOfNextActionToUpdate == 0 {
			return ActionUpdateResult{IsDone: false}, nil
		}
		if action.deployAction.ContainerId == nil {
			return ActionUpdateResult{IsDone: false}, nil
		}

		ctr := clientState.GetContainerFromContainerId(*action.deployAction.ContainerId)
		loadbalancerContainer := clientState.GetSingleContainerForService(action.deployAction.Project.Loadbalancer.Service)
		//No tls container yet found, wait for the discovery to discover this container
		if ctr == nil || ctr.GetIp() == nil {
			return ActionUpdateResult{IsDone: false}, nil
		}
		//No loadbalancer container yet found, wait for the discovery to discover this container
		if loadbalancerContainer == nil {
			return ActionUpdateResult{IsDone: false}, nil
		}
		for _, certificateBlock := range action.handler.Blocks {
			action.Actions.Actions = append(action.Actions.Actions,
				CreateLoadbalancerAction(
					action.deployAction.Node,
					loadbalancerContainer, //TOOD: This must be the loadbalancer container
					action.deployAction.Project,
					clientState.Containers,
				),
				&CheckReadinessProbe{
					MustSucceed: true,
					Node:        action.deployAction.Node,
					Probe:       action.handler.ServiceJob.Service.Probes.Ready,
					Service:     action.handler.ServiceJob.Service,
					Container:   ctr,
				},
				&IssueTlsCertificate{CertificationBlock: &certificateBlock, Container: ctr},
				&WaitForTlsCertificate{CertificationBlock: &certificateBlock, Container: ctr},
			)
		}
		action.Finally.Actions = append(action.Finally.Actions, &RemoveContainer{Node: action.deployAction.Node, Container: ctr})

		action.state = container_deployed
		return ActionUpdateResult{IsDone: false}, nil
	}
	return ActionUpdateResult{IsDone: true}, nil
}

func (action *CertificateMasterJob) Equals(otherAction Action) bool {
	other, ok := otherAction.(*CertificateMasterJob)
	if !ok {
		return false
	}

	return action.handler.Id == other.handler.Id
}

func (action *IssueTlsCertificate) Run() (ActionRunResult, error) {
	console.InfoLog.Log("Issue Tls Certificate", action)

	commandList := []string{
		"certbot",
		"certonly",
		"--webroot",
		"-w", "/var/www/certbot",
		"--agree-tos",
		"--non-interactive",
		"--dry-run",
	}

	for _, domain := range action.CertificationBlock.Domains {
		commandList = append(commandList, "-d", domain)
	}

	if len(action.CertificationBlock.Emails) > 0 {
		commandList = append(commandList, "--email", strings.Join(action.CertificationBlock.Emails, ","))
	}

	exitCode, err := docker.SendCommandToContainer(
		commandList,
		action.Container.Id,
	)
	if err != nil {
		console.InfoLog.Error("Failed to issue TLS certificate: ", err)
		return ActionRunResult{IsDone: false}, err
	}
	if exitCode != 0 {
		console.InfoLog.Error("Failed to issue TLS certificate: ", exitCode)
		return ActionRunResult{IsDone: false}, errors.New("failed to issue TLS certificate")
	}
	return ActionRunResult{IsDone: true}, nil
}

func (action *IssueTlsCertificate) Update(actions *[]Action, clientState *clientState.ClientState) (ActionUpdateResult, error) {
	return ActionUpdateResult{IsDone: true}, nil
}

func (action *IssueTlsCertificate) Equals(otherAction Action) bool {
	return false
}

func (action *WaitForTlsCertificate) Run() (ActionRunResult, error) {
	console.InfoLog.Log("Issue Tls Certificate", action)
	return ActionRunResult{IsDone: true}, nil
}

func (action *WaitForTlsCertificate) Update(actions *[]Action, clientState *clientState.ClientState) (ActionUpdateResult, error) {
	return ActionUpdateResult{IsDone: true}, nil
}

func (action *WaitForTlsCertificate) Equals(otherAction Action) bool {
	return false
}
