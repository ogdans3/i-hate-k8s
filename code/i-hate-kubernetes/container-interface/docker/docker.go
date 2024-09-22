package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models"
	"io"
	"log"
	"os"
)

func ListAllContainers() {
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer apiClient.Close()

	containers, err := apiClient.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		panic(err)
	}

	if len(containers) == 0 {
		fmt.Printf("No containers\n")
	}
	for _, ctr := range containers {
		fmt.Printf("%s %s (status: %s)\n", ctr.ID, ctr.Image, ctr.Status)
	}
}

func StopProjectContainers(project *models.Project) {
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer apiClient.Close()

	ctx := context.Background()
	containers, err := apiClient.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		panic(err)
	}

	if len(containers) == 0 {
		fmt.Printf("No containers\n")
	}
	timeout := 10
	for _, ctr := range containers {
		fmt.Printf("%s %s (status: %s)\n", ctr.ID, ctr.Image, ctr.Status)
		log.Printf("Stopping container %s", ctr.ID)
		err := apiClient.ContainerStop(ctx, ctr.ID, container.StopOptions{
			Timeout: &timeout,
		})
		if err != nil {
			log.Printf("Failed to stop container: %s, %s", ctr.ID, err)
			//TODO: Handle error?
			continue
		}
		log.Printf("%s stopped", ctr.ID)
		log.Printf("Removing container, volumes, and links for %s", ctr.ID)
		err = apiClient.ContainerRemove(ctx, ctr.ID, container.RemoveOptions{
			RemoveLinks:   false, //TODO: If this is true, then i get this error: Error response from daemon: Conflict, cannot remove the default link name of the container
			RemoveVolumes: true,
			Force:         true,
		})
		if err != nil {
			log.Printf("Failed to remove container, volumes, or links: %s, %s", ctr.ID, err)
			//TODO: Handle error?
			continue
		}
		log.Printf("%s removed", ctr.ID)
	}
}

func CreateContainerFromService(service models.Service) {
	ctx := context.Background()

	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer apiClient.Close()

	imageName := service.Image

	reader, err := apiClient.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	io.Copy(os.Stdout, reader)

	fmt.Printf("%v\n", reader)

	createdContainer, err := apiClient.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	},
		&container.HostConfig{
			PortBindings: portToPortBinding(service.Ports),
		}, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := apiClient.ContainerStart(ctx, createdContainer.ID, container.StartOptions{}); err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", createdContainer)
}

func portToPortBinding(ports []models.Port) map[nat.Port][]nat.PortBinding {
	portMap := map[nat.Port][]nat.PortBinding{}
	for _, port := range ports {
		//TODO: Handle multiple from ports? Probably because of multiple protocols over the same port?
		portMap[nat.Port(fmt.Sprintf("%s/%s", port.ContainerPort, port.Protocol))] = []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: port.HostPort,
			},
		}
	}
	return portMap
}
