package client

import (
	"strings"
	"sync"
	"time"

	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/engine-interface/docker"
	engine_models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/engine-interface/engine-models"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/console"
	models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/internal-models"
	model_actions "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/internal-models/actions"
)

type Client struct {
	containers            []engine_models.Container
	projects              []models.Project
	containerToServiceMap map[string]string
	Node                  models.Node
}

func CreateClient() Client {
	console.Clear()
	client := Client{
		containers:            make([]engine_models.Container, 0),
		projects:              make([]models.Project, 0),
		containerToServiceMap: make(map[string]string, 0),
		Node: models.Node{
			Ip:       "127.0.0.1",
			Name:     "me",
			HostName: "127.0.0.1",
			Role:     models.ControlPlane,
		},
	}
	client.Update()
	return client
}

func (client *Client) Loop() {
	i := 0
	for {
		console.Spinner("Waiting for something to happen")
		client.Update()
		client.MoveTowardsDesiredState()
		time.Sleep(200 * time.Millisecond)
		i++
	}
}

func (client *Client) MoveTowardsDesiredState() {
	actions := client.CalculateActions()
	var wg sync.WaitGroup
	for _, action := range actions {
		wg.Add(1)
		go func() {
			defer wg.Done()
			action.Run()
		}()
	}
	wg.Wait()
}

func (client *Client) CalculateActions() []model_actions.Action {
	actions := make([]model_actions.Action, 0)

	for _, project := range client.projects {
		for _, service := range project.Services {
			actionsForThisService := addActionsForService(client.Node, service, client.containers)
			actions = append(actions, actionsForThisService...)
		}

		if project.Loadbalancer != nil {
			actionsForThisService := addActionsForService(client.Node, project.Loadbalancer.Service, client.containers)
			actions = append(actions, actionsForThisService...)
		}
	}

	return actions
}

func addActionsForService(node models.Node, service *models.Service, containers []engine_models.Container) []model_actions.Action {
	actions := make([]model_actions.Action, 0)
	foundContainer := false
	for _, container := range containers {
		for _, name := range container.Names {
			if strings.Contains(name, service.Id) {
				switch container.State {
				case "created":
					fallthrough
				case "restarting":
					fallthrough
				case "running":
					//The container is either about to start, or is running. So we do nothing here
					break
				case "paused":
					actions = append(actions, model_actions.CreateDeployContainerForService(*service))
				case "exited":
					fallthrough
				case "dead":
					//TODO: Should we try to restart the container? Maybe the user stopped them on purpose?
					actions = append(actions, model_actions.CreateRestartContainer(container, node))
				}
				foundContainer = true
			}
		}
	}
	if !foundContainer {
		actions = append(actions, model_actions.CreateDeployContainerForService(*service))
	}
	return actions
}

func (client *Client) AddContainerToService(containerId string, service models.Service) {
	client.containerToServiceMap[containerId] = service.Id
}

func (client *Client) Update() {
	containers := docker.ListAllContainers()

	if len(containers) == 0 {
		client.containers = make([]engine_models.Container, 0)
		return
	}
	client.containers = make([]engine_models.Container, 0) //TODO: Dont reset here, wasted resources
	for _, ctr := range containers {
		client.containers = append(client.containers, engine_models.Container{
			Id:      ctr.ID,
			Image:   ctr.Image,
			Command: ctr.Command,
			Status:  ctr.Status,
			State:   ctr.State,
			Names:   ctr.Names,

			Node: client.Node,
		})
	}
}

func (client *Client) AddProject(project models.Project) {
	client.projects = append(client.projects, project)
	/*
		if project.Loadbalancer != nil {
			docker.CreateContainerFromService(project.Loadbalancer.Service)
		}
		for _, service := range project.Services {
			docker.CreateContainerFromService(service)
		}
	*/
}

func (client *Client) GetContainers() []engine_models.Container {
	return client.containers
}

func (client *Client) Nuke() {
	docker.StopAndRemoveAllContainers()
}

func (client *Client) StopProject(project models.Project) {
	//TODO: Use the project specification to remove containers
	docker.StopAllContainers()
}
