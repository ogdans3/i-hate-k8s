package client

import (
	"fmt"
	"time"

	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/engine-interface/docker"
	engine_models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/engine-interface/engine-models"
	models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/internal-models"
)

type Client struct {
	containers []engine_models.Container
}

func CreateClient() Client {
	client := Client{
		containers: make([]engine_models.Container, 0),
	}
	client.Update()
	return client
}

func (client *Client) Loop() {
	for {
		fmt.Printf("Loop")
		client.Update()
		time.Sleep(1 * time.Second)
	}
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
		})
	}
}

func (client *Client) AddProject(project models.Project) {
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
	docker.StopAllContainers()
}

func (client *Client) StopProject(project models.Project) {
	//TODO: Use the project specification to remove containers
	docker.StopAllContainers()
}

func (client *Client) PrettyPrint() {
	for _, container := range client.containers {
		fmt.Println(container.Id, container.Status)
	}
}
