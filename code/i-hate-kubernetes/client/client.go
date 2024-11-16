package client

import (
	"context"
	"fmt"
	"math"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/api"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/api/webhooks"
	clientState "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/client-state"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/engine-interface/docker"
	engine_models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/engine-interface/engine-models"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/stats"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/console"
	models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/internal-models"
	model_actions "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/internal-models/actions"
)

type Client struct {
	ctx      context.Context
	state    clientState.ClientState
	actions  []model_actions.Action
	channels Channels
	settings clientState.ClientSettings
}

type Channels struct {
	webhookChannel chan webhooks.WebhookPayload
}

func CreateClient() Client {
	console.Clear()
	client := Client{
		channels: Channels{
			webhookChannel: make(chan webhooks.WebhookPayload),
		},
		ctx: context.Background(),
		settings: clientState.ClientSettings{
			ApiPort: "6444", //1 better than kubernetes (which is 6443)
		},
		actions: make([]model_actions.Action, 0),
		state: clientState.ClientState{
			Containers:             make([]engine_models.Container, 0),
			Networks:               make([]engine_models.Network, 0),
			Projects:               make([]models.Project, 0),
			NetworkConfiguration:   models.LoadbalancerNetworkConfiguration{},
			EngineNetworkToService: make(map[string][]models.Service, 0),
			ContainerMetadata:      make(map[string]clientState.ContainerMetadata, 0),
			CicdJobs:               make([]models.Cicd, 0),
			Node: models.Node{
				Ip:       "127.0.0.1",
				Name:     "me",
				HostName: "127.0.0.1",
				Role:     models.ControlPlane,
			},
		},
	}
	return client
}

const LOOP_DELAY = 1000
const PRINT_LOOP_DELAY = 100

func (client *Client) Loop() {
	go api.StartApiServer(client.ctx, client.settings.ApiPort, client.channels.webhookChannel)

	pid := os.Getpid()
	procStats := &stats.Stat{}
	var iterations = math.Floor(LOOP_DELAY / PRINT_LOOP_DELAY)

	i := 0
	for {
		for s := 0; s < int(iterations); s++ {
			info := fmt.Sprintf("Waiting for something to happen [GOR: %d] [MEM: %s]", procStats.Goroutines, procStats.MemoryPretty)
			if procStats.CpuUsage != nil {
				info = fmt.Sprintf("%s [CPU: %s%%]", info, procStats.CpuUsagePercentage)
			}
			console.Spinner(info)
			time.Sleep(PRINT_LOOP_DELAY * time.Millisecond)
		}

		//Handle channel updates
		select {
		case payload := <-client.channels.webhookChannel:
			console.InfoLog.Info("Received payload from handler: ", i, payload, "\n")
			client.AddCicdJobSpecFromWebhook(payload)
		default:
		}

		client.Update()
		client.MoveTowardsDesiredState()
		i++

		//Read some stats for fun
		stats.GetProcessStats(pid, procStats)
		console.StatLog.Info(procStats)
	}
}

func (client *Client) AddCicdJobSpecFromWebhook(payload webhooks.WebhookPayload) {
	refParts := strings.Split(payload.Ref, "/")
	branch := strings.Join(refParts[2:], "/")

	for _, project := range client.state.Projects {
		for _, cicdSpec := range project.Cicd {
			if branch != cicdSpec.Branch {
				continue
			}
			if cicdSpec.Url == payload.Repository.GitUrl ||
				cicdSpec.Url == payload.Repository.HtmlUrl ||
				cicdSpec.Url == payload.Repository.CloneUrl ||
				cicdSpec.Url == payload.Repository.SshUrl {
				client.state.CicdJobs = append(client.state.CicdJobs, cicdSpec)
				return
			}
		}
	}
}

func (client *Client) MoveTowardsDesiredState() {
	client.CalculateActions(&client.actions)
	model_actions.GetDistinctActions(&client.actions)
	remainingActions := make([]model_actions.Action, 0)
	var wg sync.WaitGroup
	var mu sync.Mutex //TODO Will using a mutex here make it too slow?
	for _, action := range client.actions {
		wg.Add(1)
		go func() {
			defer wg.Done()
			actionRunResult, err := action.Run()
			if err != nil {
				console.Log(err)
				return
			}

			mu.Lock()
			if actionRunResult.IsDone {
				action.Update(&remainingActions, &client.state)
			} else {
				remainingActions = append(remainingActions, action)
			}
			mu.Unlock()
		}()
	}
	wg.Wait()
	client.actions = remainingActions
}

func (client *Client) CalculateActions(actions *[]model_actions.Action) {
	for _, project := range client.state.Projects {
		//Add actions to create networks
		for _, service := range project.Services {
			if service.Network.GetName() == nil {
				continue
			}
			foundNetworkForService := false
			for _, network := range client.state.Networks {
				if network.Name == service.Network.Name {
					foundNetworkForService = true
					break
				}
			}
			if !foundNetworkForService {
				*actions = append(*actions, &model_actions.CreateNetwork{
					Node:    &client.state.Node,
					Service: service,
				})
			}
		}

		//Add actions for registry (e.g. deploy the registry)
		if project.Registry != nil {
			actionsForThisService := addActionsForService(client, client.state.Node, project.Registry.Service, client.state.Containers, project)
			if len(actionsForThisService) > 0 {
				*actions = append(*actions, actionsForThisService...)
				//The registry must be available for us to continue
				return
			} else if !isRegistryAvailable(client, client.state.Node, project.Registry.Service, client.state.Containers) {
				//The registry must be available for us to continue
				// We could be waiting for a probe, for example.
				return
			}
		}

		//Add actions for services (e.g. deploy new container, restart container)
		for _, service := range project.Services {
			actionsForThisService := addActionsForService(client, client.state.Node, service, client.state.Containers, project)
			*actions = append(*actions, actionsForThisService...)
		}

		//Add actions for loadbalancer (e.g. deploy the loadbalancer)
		if project.Loadbalancer != nil {
			actionsForThisService := addActionsForService(client, client.state.Node, project.Loadbalancer.Service, client.state.Containers, project)
			*actions = append(*actions, actionsForThisService...)
		}

		//Add actions for network config for loadbalancer (e.g. a container changed port, a new container has been deployed)
		if project.Loadbalancer != nil {
			actionsForThisService := addLoadbalancerActions(client, client.state.Node, client.state.NetworkConfiguration, client.state.Containers, project.Services, *project.Loadbalancer.Service)
			*actions = append(*actions, actionsForThisService...)
		}

		//TODO: Seems inefficient to loop the projects for these jobs just to separate between different type of jobs.
		//Add actions for cicd jobs (e.g. build new image for container, update i-hate-kubernetes)
		for _, cicdJob := range client.state.CicdJobs {
			if project.Autoupdate != nil && cicdJob.Id == project.Autoupdate.Id {
				*actions = append(*actions, model_actions.CreateCicdUpdateIHateKubernetes(
					&client.state.Node,
					&cicdJob,
				))
			} else {
				*actions = append(*actions, model_actions.CreateCicdJob(
					&client.state.Node,
					&cicdJob,
					cicdJob.Service,
					&project,
				))
			}
		}
	}
	client.state.CicdJobs = make([]models.Cicd, 0)
}

func addLoadbalancerActions(
	client *Client,
	node models.Node,
	currentLoadbalancerNetworkConfiguration models.LoadbalancerNetworkConfiguration,
	containers []engine_models.Container,
	services map[string]*models.Service,
	loadBalancerService models.Service,
) []model_actions.Action {
	actions := make([]model_actions.Action, 0)

	var loadBalancerContainer *engine_models.Container = nil
	for _, container := range containers {
		for _, name := range container.Names {
			if strings.Contains(name, loadBalancerService.Id) && container.State == "running" {
				loadBalancerContainer = &container
			}
		}
	}
	//The loadbalancer container is not yet available, so we return here
	if loadBalancerContainer == nil {
		return actions
	}

	newConfig := models.LoadbalancerNetworkConfiguration{
		ContainerIdOfLoadbalancerThatHasThisConfig: &loadBalancerContainer.Id,
	}
	for _, service := range services {
		upstreamBlock := models.Upstream{
			Name:    service.ServiceName,
			Servers: []models.UpstreamServer{},
		}
		serverBlock := models.Server{
			Location: []models.ServerLocation{},
		}
		//TODO: Loop over path/domain for the service and insert into the location block
		serverBlock.Location = append(serverBlock.Location, models.ServerLocation{
			MatchModifier: "",
			LocationMatch: "/",
			ProxyPass:     upstreamBlock.Name,
		})

		for _, container := range containers {
			for _, name := range container.Names {
				if strings.Contains(name, service.Id) {
					if container.Ip != nil {
						upstreamBlock.Servers = append(upstreamBlock.Servers, models.UpstreamServer{
							Server: *container.GetIp(),
						})
					}
				}
			}
		}
		newConfig.HttpBlocks = append(newConfig.HttpBlocks, models.Http{
			Upstream: []models.Upstream{
				upstreamBlock,
			},
			Server: []models.Server{
				serverBlock,
			},
		})
	}
	if newConfig.ConfigurationToNginxFile() == currentLoadbalancerNetworkConfiguration.ConfigurationToNginxFile() &&
		currentLoadbalancerNetworkConfiguration.ContainerIdOfLoadbalancerThatHasThisConfig != nil &&
		*currentLoadbalancerNetworkConfiguration.ContainerIdOfLoadbalancerThatHasThisConfig == loadBalancerContainer.Id {
		return actions
	}

	console.InfoLog.Debug(newConfig.ConfigurationToNginxFile())
	actions = append(actions, &model_actions.UpdateLoadbalancer{
		Node:                 &node,
		Container:            loadBalancerContainer,
		NetworkConfiguration: &newConfig,
	})

	return actions
}

func isRegistryAvailable(client *Client, node models.Node, service *models.Service, containers []engine_models.Container) bool {
	for _, container := range containers {
		for _, name := range container.Names {
			if strings.Contains(name, service.Id) {
				switch container.State {
				case "created":
					fallthrough
				case "restarting":
					return false
				case "running":
					//We dont need to check if the probe should run, because we should never enter this code if the probe should run
					if service.Probes != nil && service.Probes.Ready != nil {
						return client.state.ContainerMetadata[container.Id].ProbesMetadata.Readiness.ResultOfLastCheck
					}
					return true
				case "paused":
					return false
				case "exited":
					fallthrough
				case "dead":
					return false
				}
			}
		}
	}
	return false
}

func addActionsForService(c *Client, node models.Node, service *models.Service, containers []engine_models.Container, project models.Project) []model_actions.Action {
	actions := make([]model_actions.Action, 0)
	containersFound := 0
	for _, container := range containers {
		for _, name := range container.Names {
			if strings.Contains(name, service.Id) {
				switch container.State {
				case "created":
					fallthrough
				case "restarting":
					//The container is starting
					//TODO: Add probes
					break
				case "running":
					//The container is either about to start, or is running. So we do nothing here
					if service.Probes != nil && service.Probes.Ready != nil {
						action := model_actions.CreateReadinessProbe(
							&node,
							service,
							&container,
							service.Probes.Ready,
							c.state.ContainerMetadata[container.Id], //TODO: Handler pointers and stufs?
						)
						if action != nil {
							actions = append(actions, action)
						}
					}
				case "paused":
					actions = append(actions, model_actions.CreateDeployContainerForService(service, &project))
				case "exited":
					fallthrough
				case "dead":
					//TODO: Should we try to restart the container? Maybe the user stopped them on purpose?
					actions = append(actions, model_actions.CreateRestartContainer(&container, &node))
				}
				containersFound++
			}
		}
	}
	if containersFound < int(service.Autoscale.Initial) {
		for i := containersFound; i < int(service.Autoscale.Initial); i++ {
			actions = append(actions, model_actions.CreateDeployContainerForService(service, &project))
			console.Log("Add action to create deploy container: ", service.ServiceName)
		}
	}
	return actions
}

func (client *Client) Update() {
	containers, err := docker.ListAllContainers()
	if err != nil {
		return
	}

	client.state.Containers = make([]engine_models.Container, 0) //TODO: Dont reset here, wasted resources
	for _, ctr := range containers {
		var project *models.Project
		var service *models.Service
		for _, p := range client.state.Projects {
			for _, name := range ctr.Names {
				if strings.Contains(name, p.Project) {
					project = &p
				}
			}
		}
		if project != nil {
			for _, name := range ctr.Names {
				if strings.Contains(name, project.Loadbalancer.Service.Id) {
					service = project.Loadbalancer.Service
					break
				}
				if strings.Contains(name, project.Registry.Service.Id) {
					service = project.Registry.Service
					break
				}
				for _, s := range project.Services {
					if strings.Contains(name, s.Id) {
						service = s
					}
				}
			}
		}

		var ip *string
		if service != nil && service.Network.GetName() != nil {
			//TODO: Handle multiple networks, or atleast use a default network for the loadbalancer
			if ctr.NetworkSettings.Networks[*service.Network.GetName()] != nil {
				ip = &ctr.NetworkSettings.Networks[*service.Network.GetName()].IPAddress
			}
		}

		client.state.Containers = append(client.state.Containers, engine_models.Container{
			Id:      ctr.ID,
			Image:   ctr.Image,
			Command: ctr.Command,
			Status:  ctr.Status,
			State:   ctr.State,
			Names:   ctr.Names,
			Ip:      ip,

			ProjectIdentifier: project.GetId(),
			ServiceIdentifier: service.GetId(),

			Node: client.state.Node,
		})
		if _, ok := client.state.ContainerMetadata[ctr.ID]; !ok {
			client.state.ContainerMetadata[ctr.ID] = clientState.ContainerMetadata{
				ProbesMetadata: &clientState.ProbesMetadata{
					Started:   &clientState.ProbeMetadata{},
					Liveness:  &clientState.ProbeMetadata{},
					Readiness: &clientState.ProbeMetadata{LastCheck: time.Now().Unix()},
				},
			}
		}
	}

	networkSummaries, err := docker.ListNetworks()
	if err != nil {
		//TODO: Handle error
		console.Log(err)
		return
	}

	client.state.Networks = make([]engine_models.Network, 0) //TODO: Dont reset here, wasted resources
	for _, networkSummary := range networkSummaries {
		client.state.Networks = append(client.state.Networks, engine_models.Network{
			Id:   networkSummary.ID,
			Name: networkSummary.Name,
		})
	}
}

func (client *Client) AddProject(project models.Project) {
	client.state.Projects = append(client.state.Projects, project)
}

func (client *Client) GetContainers() []engine_models.Container {
	return client.state.Containers
}

func (client *Client) Nuke() {
	docker.StopAndRemoveAllContainersAndNetworks()
}

func (client *Client) StopProject(project models.Project) {
	//TODO: Use the project specification to remove containers
	docker.StopAllContainers()
}
