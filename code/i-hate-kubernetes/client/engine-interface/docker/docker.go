package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/go-connections/nat"
	engine_models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/engine-interface/engine-models"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/console"
	models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/internal-models"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/util"
)

func createDockerClient() *client.Client {
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		//TODO: Try a few times before erroring?
		panic(err)
	}
	return apiClient
}

func ListAllContainers() ([]types.Container, error) {
	apiClient := createDockerClient()
	defer apiClient.Close()

	containers, err := apiClient.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		//TODO: Try a few times before erroring?
		console.Error("Unable to list containers", err)
		return nil, err
	}

	return containers, nil
}

func StopAllContainers() {
	apiClient := createDockerClient()
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

func StopContainer(timeout int, ctr engine_models.Container) error {
	apiClient := createDockerClient()
	defer apiClient.Close()
	ctx := context.Background()

	err := apiClient.ContainerStop(ctx, ctr.Id, container.StopOptions{
		Timeout: &timeout,
	})
	if err != nil {
		console.Error("Failed to stop container: %s, %s", ctr.Id, err)
		return err
		//TODO: Send error back
	}
	return nil
}

func StopAndRemoveContainer(ctr engine_models.Container) error {
	apiClient := createDockerClient()
	defer apiClient.Close()
	ctx := context.Background()

	err := StopContainer(10, ctr)
	if err != nil {
		return err
	}

	err = apiClient.ContainerRemove(ctx, ctr.Id, container.RemoveOptions{
		RemoveLinks:   false, //TODO: If this is true, then i get this error: Error response from daemon: Conflict, cannot remove the default link name of the container
		RemoveVolumes: true,  //TODO: Check if persistent volume?
		Force:         true,  //TODO: Should probably not force the first time
	})
	if err != nil {
		fmt.Println(err)
		console.Error("Failed to remove container, volumes, or links: %s, %s", ctr.Id, err)
		return err
		//TODO: Handle error?
	}

	//TODO: Remove container
	//TODO: Remove volumes
	//TODO: Remove solo networks
	//TODO: Remove images
	return nil
}

func StopAndRemoveAllContainersAndNetworks() {
	apiClient := createDockerClient()
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
				fmt.Println(err)
				console.Error("Failed to remove container, volumes, or links: %s, %s", ctr.ID, err)
				//TODO: Handle error?
			}
		}()
	}
	wg.Wait()

	networks, err := apiClient.NetworkList(ctx, network.ListOptions{})
	for _, n := range networks {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err = apiClient.NetworkRemove(ctx, n.ID)
			if err != nil {
				fmt.Println(err)
				console.Error("Failed to remove network: %s, %s", n.ID, err)
				//TODO: Handle error?
			}
		}()
	}
	wg.Wait()

	images, err := apiClient.ImageList(ctx, image.ListOptions{All: true})
	for _, i := range images {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := apiClient.ImageRemove(ctx, i.ID, image.RemoveOptions{})
			if err != nil {
				fmt.Println(err)
				console.Error("Failed to remove network: %s, %s", i.ID, err)
				//TODO: Handle error?
			}
		}()
	}
	wg.Wait()
}

func BuildService(service models.Service, project models.Project) {
	apiClient := createDockerClient()
	defer apiClient.Close()
	ctx := context.Background()

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	defer tw.Close()

	tarContext, err := archive.Tar(filepath.Join(service.Directory), 0)
	if err != nil {
		log.Fatalf("Error creating tar context: %v", err)
	}

	imageName := service.Image
	imageTag := imageName
	if service.Build && project.Registry != nil {
		//TODO: Get this from the actual registry service or container?
		imageTag = "localhost:5000/" + imageName
		imageName = "localhost:5000/" + imageName
	}

	response, err := apiClient.ImageBuild(ctx, tarContext, types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Remove:     true,
		Tags:       []string{imageName},
		PullParent: true,
	})
	if err != nil {
		fmt.Println(err)
		console.Fatal(err)
	}
	defer response.Body.Close()
	console.Copy(response.Body)

	imagePushResponse, err := apiClient.ImagePush(ctx, imageTag, image.PushOptions{
		All:          true,
		RegistryAuth: "TODO2", //TODO: The registry auth must be there, but the value does not matter
	})
	if err != nil {
		fmt.Println(err)
		console.Fatal(err)
	}
	defer imagePushResponse.Close()
	console.Copy(imagePushResponse)
	console.Log("Image built and pushed")
}

func CreateContainerFromService(service models.Service, project *models.Project) (*string, error) {
	apiClient := createDockerClient()
	defer apiClient.Close()
	ctx := context.Background()

	imageName := service.Image

	fmt.Println(project)
	if service.Build && project.Registry != nil {
		//TODO: Get this from the actual registry service or container?
		imageName = "localhost:5000/" + imageName
	} else {
		imageName = "docker.io/library/" + imageName
	}

	console.Log("Image name: ", imageName)
	reader, err := apiClient.ImagePull(ctx, imageName, image.PullOptions{
		All:          false,   //Very specifically set to false! If this is true we start pulling every image in the world
		RegistryAuth: "TODO2", //TODO: The registry auth must be there, but the value does not matter
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer reader.Close()
	console.Copy(reader)

	networkName := service.Network.GetName()
	var networkConfig *network.NetworkingConfig
	if networkName != nil {
		networkConfig = &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				*networkName: {},
			},
		}
	}

	createdContainer, err := apiClient.ContainerCreate(
		ctx,
		&container.Config{
			Image: imageName,
		},
		&container.HostConfig{
			PortBindings: portToPortBinding(service.Ports),
		},
		networkConfig,
		nil,
		service.Id+"-"+service.ContainerName+"-"+util.RandStringBytesMaskImpr(5),
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = StartContainer(createdContainer.ID)

	return &createdContainer.ID, err
}

func StartContainer(containerId string) error {
	apiClient := createDockerClient()
	defer apiClient.Close()
	ctx := context.Background()

	if err := apiClient.ContainerStart(ctx, containerId, container.StartOptions{}); err != nil {
		console.Error("Unable to start container")
		return err
	}
	return nil
}

func CreateNetwork(n models.Network) (*string, error) {
	apiClient := createDockerClient()
	defer apiClient.Close()
	ctx := context.Background()

	if n.GetName() == nil {
		return nil, errors.New("no network name specified")
	}

	networkCreateResponse, err := apiClient.NetworkCreate(ctx, *n.GetName(), network.CreateOptions{
		Driver: "bridge", //TODO isnt this wrong. Shouldnt these networks be isolated from the internet?
	})
	if err != nil {
		console.Error("Unable to create network")
		return nil, err
	}

	return &networkCreateResponse.ID, nil
}

func ListNetworks() ([]network.Inspect, error) {
	apiClient := createDockerClient()
	defer apiClient.Close()
	ctx := context.Background()

	response, err := apiClient.NetworkList(ctx, network.ListOptions{})
	if err != nil {
		console.Error("Unable to list networks")
		return nil, err
	}

	return response, nil
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
