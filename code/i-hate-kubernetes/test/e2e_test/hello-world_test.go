package e2e_test

import (
	"testing"

	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/engine-interface/docker"
	external_models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/external-models"
	models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/internal-models"
)

func TestEmptyEnvironment(t *testing.T) {
	c := client.CreateClient()
	c.Nuke()

	c.Update() //TODO: Remove this forced update
	containers := c.GetContainers()

	if len(containers) != 0 {
		t.Errorf("Found %q containers, wanted %q containers", len(containers), 0)
	}
}

func TestEnvironmentWithOnlyALoadBalancer(t *testing.T) {
	c := client.CreateClient()
	c.Nuke()
	lb := models.ParseLoadBalancer(true)
	docker.CreateContainerFromService(lb.Service)

	c.Update() //TODO: Remove this forced update
	containers := c.GetContainers()

	if len(containers) != 1 {
		t.Errorf("Found %q containers, wanted %q containers", len(containers), 1)
	}
}

func TestEnvironmentWithOneApplicationContainer(t *testing.T) {
	c := client.CreateClient()
	c.Nuke()
	docker.CreateContainerFromService(models.ParseService(external_models.Service{
		Image: "strm/helloworld-http",
		Autoscale: external_models.Autoscale{
			Initial:   1,
			Autoscale: false,
		},
	}))

	c.Update() //TODO: Remove this forced update
	containers := c.GetContainers()

	if len(containers) != 1 {
		t.Errorf("Found %q containers, wanted %q containers", len(containers), 1)
	}
}

func TestEnvironmentWithTwoApplicationContainers(t *testing.T) {
	c := client.CreateClient()
	c.Nuke()
	docker.CreateContainerFromService(models.ParseService(external_models.Service{
		Image: "strm/helloworld-http",
		Autoscale: external_models.Autoscale{
			Initial:   1,
			Autoscale: false,
		},
	}))
	docker.CreateContainerFromService(models.ParseService(external_models.Service{
		Image: "strm/helloworld-http",
		Autoscale: external_models.Autoscale{
			Initial:   1,
			Autoscale: false,
		},
	}))

	c.Update() //TODO: Remove this forced update
	containers := c.GetContainers()

	if len(containers) != 2 {
		t.Errorf("Found %q containers, wanted %q containers", len(containers), 2)
	}
}

func TestEnvironmentWithTwoApplicationContainersAndLoadbalancer(t *testing.T) {
	c := client.CreateClient()
	c.Nuke()
	docker.CreateContainerFromService(models.ParseService(external_models.Service{
		Image: "strm/helloworld-http",
		Autoscale: external_models.Autoscale{
			Initial:   1,
			Autoscale: false,
		},
	}))
	docker.CreateContainerFromService(models.ParseService(external_models.Service{
		Image: "strm/helloworld-http",
		Autoscale: external_models.Autoscale{
			Initial:   1,
			Autoscale: false,
		},
	}))
	lb := models.ParseLoadBalancer(true)
	docker.CreateContainerFromService(lb.Service)

	c.Update() //TODO: Remove this forced update
	containers := c.GetContainers()

	if len(containers) != 3 {
		t.Errorf("Found %q containers, wanted %q containers", len(containers), 3)
	}
}
