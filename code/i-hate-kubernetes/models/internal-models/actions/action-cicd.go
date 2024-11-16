package model_actions

import (
	"os"
	"os/exec"
	"syscall"

	clientState "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/client-state"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/engine-interface/docker"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/console"
	models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/internal-models"
)

type CicdUpdateImage struct {
	DefaultActionMetadata
	Node    *models.Node
	Cicd    *models.Cicd
	Service *models.Service
	Project *models.Project
	//Container *engine_models.Container //Add a container to run the job in?
}

type CicdUpdateIHateKubernetes struct {
	DefaultActionMetadata
	Node *models.Node
	Cicd *models.Cicd
}

func CreateCicdJob(node *models.Node, cicd *models.Cicd, service *models.Service, project *models.Project) *CicdUpdateImage {
	return &CicdUpdateImage{
		Node:    node,
		Cicd:    cicd,
		Service: service,
		Project: project,
	}
}

func (action *CicdUpdateImage) Run() (ActionRunResult, error) {
	console.Log("Run cicd update image job: ", action.Cicd)
	cmd := exec.Command("git", "pull", "origin", "master")
	cmd.Dir = action.Cicd.Directory

	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	console.InfoLog.Info(string(output))

	docker.BuildService(*action.Service, *action.Project)
	return ActionRunResult{IsDone: true}, nil
}

func (action *CicdUpdateImage) Update(actions *[]Action, clientState *clientState.ClientState) (ActionUpdateResult, error) {
	*actions = append(*actions, CreateContainerImageUpdated(clientState, action.Node, action.Project, action.Service))
	return ActionUpdateResult{IsDone: true}, nil
}

func (action *CicdUpdateImage) Equals(otherAction Action) bool {
	other, ok := otherAction.(*CicdUpdateImage)
	if !ok {
		return false
	}

	return action.Node == other.Node &&
		action.Cicd == other.Cicd
}

func CreateCicdUpdateIHateKubernetes(node *models.Node, cicd *models.Cicd) *CicdUpdateIHateKubernetes {
	return &CicdUpdateIHateKubernetes{
		Node: node,
		Cicd: cicd,
	}
}

// TODO: Handle errors properly
func (action *CicdUpdateIHateKubernetes) Run() (ActionRunResult, error) {
	console.Log("Update main program: ", action.Cicd)
	console.Log("Run cicd update i-hate-kubernetes job: ", action.Cicd)
	cmd := exec.Command("git", "pull", "origin", "master")
	cmd.Dir = action.Cicd.Directory

	output, err := cmd.CombinedOutput()
	if err != nil {
		console.InfoLog.Error("Error pulling application code:", err)
		return ActionRunResult{IsDone: true}, err
	}
	console.InfoLog.Info(string(output))

	cmd = exec.Command("go", "build")
	if action.Cicd.Directory != "" {
		cmd.Dir = action.Cicd.Directory
	}

	output, err = cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	console.InfoLog.Info(string(output))

	// Get the current program's arguments
	args := os.Args

	// Restart the program with the same arguments
	cmd = exec.Command("./i-hate-kubernetes", args[1:]...) // Exclude the program name from args
	if action.Cicd.Directory != "" {
		cmd.Dir = action.Cicd.Directory
	}
	cmd.Env = os.Environ() // Retain the environment variables

	// To detach the new process, you can set the SysProcAttr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

	err = cmd.Start()
	if err != nil {
		console.InfoLog.Error("Error starting the new process:", err)
		return ActionRunResult{IsDone: true}, nil
	}

	// Exit the current program (optional, as it will be restarted)
	console.InfoLog.Info("Program restarted without error. Exiting this instance")
	os.Exit(0)
	return ActionRunResult{IsDone: true}, nil
}

func (action *CicdUpdateIHateKubernetes) Update(actions *[]Action, clientState *clientState.ClientState) (ActionUpdateResult, error) {
	return ActionUpdateResult{IsDone: true}, nil
}

func (action *CicdUpdateIHateKubernetes) Equals(otherAction Action) bool {
	other, ok := otherAction.(*CicdUpdateIHateKubernetes)
	if !ok {
		return false
	}

	return action.Node == other.Node &&
		action.Cicd == other.Cicd
}
