package engine_models

import models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/internal-models"

type Container struct {
	Id      string
	Image   string
	Command string
	Status  string
	State   string //One of created, restarting, running, removing, paused, exited, or dead
	Names   []string
	Ip      string

	Node models.Node
}
