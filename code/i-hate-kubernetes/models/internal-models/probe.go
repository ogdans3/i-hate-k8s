package models

import (
	external_models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/external-models"
)

type seconds uint8
type httpMethod string

// TOOD: Fix these types later
type Probe struct {
	Command          *string     //If you want to run say cat to check a file or some random shit
	Method           *httpMethod //The http method to use
	Path             *string     //The http path to use, e.g. /ready
	Port             *Port       //The port to check on, if not specified then we default to the container port and fail if multiple ports
	Interval         seconds     //How often to run the probe
	InitialWait      seconds     //How long to wait before first probe
	FailureThreshold seconds     //Not sure what this is?
}

type Probes struct {
	Ready    *Probe
	Started  *Probe
	Liveness *Probe
}

func ParseProbe(probe external_models.Probe) *Probe {
	probeString := string(probe)
	method := httpMethod("GET")
	return &Probe{
		Interval:         5,
		InitialWait:      5,
		FailureThreshold: 30,
		Command:          nil,
		Port:             nil,
		Method:           &method,
		Path:             &probeString,
	}
}

func ParseProbes(probes *external_models.Probes) *Probes {
	if probes == nil {
		return nil
	}
	return &Probes{
		Ready:    ParseProbe(probes.Ready),
		Started:  ParseProbe(probes.Started),
		Liveness: ParseProbe(probes.Liveness),
	}
}
