package external_models

// TOOD: Fix these types later
type Probe string

type Probes struct {
	Liveness Probe
	Ready    Probe
	Started  Probe
}
