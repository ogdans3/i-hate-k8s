package models

import external_models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/external-models"

type Autoscale struct {
	Initial   int8 //How many pods to initially deploy, defaults to 1
	Autoscale bool //If this pod should be autoscaled
}

func ParseAutoscale(autoscale external_models.Autoscale) Autoscale {
	return Autoscale{
		Initial:   autoscale.Initial,
		Autoscale: autoscale.Autoscale,
	}
}
