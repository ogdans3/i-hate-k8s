package docker

import (
	"context"
	"fmt"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/console"
	models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/internal-models"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/util"
)

func ListAllContainers() []types.Container {
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer apiClient.Close()

	containers, err := apiClient.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		panic(err)
	}

	return containers
}

func StopAllContainers() {
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

	timeout := 10
	var wg sync.WaitGroup
	for _, ctr := range containers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := apiClient.ContainerStop(ctx, ctr.ID, container.StopOptions{
				Timeout: &timeout,
			})
			if err != nil {
				console.Error("Failed to stop container: %s, %s", ctr.ID, err)
				//TODO: Send error back
			}
		}()
	}
	wg.Wait()
}

func StopAndRemoveAllContainers() {
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
	StopAllContainers()

	var wg sync.WaitGroup
	for _, ctr := range containers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err = apiClient.ContainerRemove(ctx, ctr.ID, container.RemoveOptions{
				RemoveLinks:   false, //TODO: If this is true, then i get this error: Error response from daemon: Conflict, cannot remove the default link name of the container
				RemoveVolumes: true,
				Force:         true,
			})
			if err != nil {
				console.Error("Failed to remove container, volumes, or links: %s, %s", ctr.ID, err)
				//TODO: Handle error?
			}
		}()
	}
	wg.Wait()
}

func CreateContainerFromService(service models.Service) string {
	ctx := context.Background()

	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer apiClient.Close()

	imageName := service.Image

	reader, err := apiClient.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		console.Error(err)
	}
	defer reader.Close()

	createdContainer, err := apiClient.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	},
		&container.HostConfig{
			PortBindings: portToPortBinding(service.Ports),
		}, nil, nil, service.Id+"-"+service.ContainerName+"-"+util.RandStringBytesMaskImpr(5))

	if err != nil {
		panic(err)
	}

	StartContainer(createdContainer.ID)

	return createdContainer.ID
}

func StartContainer(containerId string) {
	ctx := context.Background()

	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer apiClient.Close()

	if err := apiClient.ContainerStart(ctx, containerId, container.StartOptions{}); err != nil {
		panic(err)
	}
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
