package clientState

import (
	engine_models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/engine-interface/engine-models"
	models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/internal-models"
)

type ProbeMetadata struct {
	LastCheck         int64
	ResultOfLastCheck bool
}

type ProbesMetadata struct {
	Started   *ProbeMetadata
	Liveness  *ProbeMetadata
	Readiness *ProbeMetadata
}

type ContainerMetadata struct {
	ProbesMetadata *ProbesMetadata
}

type ClientState struct {
	Containers             []engine_models.Container
	Networks               []engine_models.Network
	Volumes                []engine_models.Volume
	Projects               []models.Project
	Node                   models.Node
	NetworkConfiguration   models.LoadbalancerNetworkConfiguration
	EngineNetworkToService map[string][]models.Service
	ContainerMetadata      map[string]ContainerMetadata
	CicdJobs               []models.Cicd //Cicd jobs that are queued and waiting to be ran
}

type ClientSettings struct {
	ApiPort string
}
