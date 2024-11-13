package clientState

import (
	engine_models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/engine-interface/engine-models"
	models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/internal-models"
)

type ClientState struct {
	Containers             []engine_models.Container
	Networks               []engine_models.Network
	Projects               []models.Project
	Node                   models.Node
	NetworkConfiguration   models.LoadbalancerNetworkConfiguration
	EngineNetworkToService map[string][]models.Service
}
