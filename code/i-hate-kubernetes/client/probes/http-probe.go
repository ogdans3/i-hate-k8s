package probes

import (
	"net/http"
	"path"

	engine_models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/engine-interface/engine-models"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/console"
	models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/internal-models"
)

func ProbeHttpGetMethod(service *models.Service, container *engine_models.Container, probe models.Probe) bool {
	var port string
	if probe.Port != nil {
		port = probe.Port.ContainerPort
	} else {
		if len(service.Ports) != 1 {
			panic("The service ports must only have one entry, because you did not specify a port in the probe")
		}
		port = service.Ports[0].ContainerPort
	}

	//TODO: This wont work if the container is unreachable on the host? Because of the network configuration
	// Which is what we want, so that it is secure by default. Do we then need a new pod that can command into
	// which then sends the http requests and returns a response to us???
	url := *container.GetIp() + ":" + port
	url = path.Join(url, *probe.Path)
	url = "http://" + url
	resp, err := http.Get(url)
	if err != nil {
		console.Debug(err)
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}
