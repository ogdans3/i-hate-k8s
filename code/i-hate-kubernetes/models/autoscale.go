package models

import external_models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/external-models"

type Autoscale struct {
	Initial   int8
	Autoscale bool
}

func ParseAutoscale(autoscale external_models.Autoscale) Autoscale {
	return Autoscale{
		Initial:   autoscale.Initial,
		Autoscale: autoscale.Autoscale,
	}

}
