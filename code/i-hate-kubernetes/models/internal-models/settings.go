package models

import external_models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/external-models"

type Settings struct {
	Interval int8
}

func ParseSettings(settings external_models.Settings) Settings {
	return Settings{
		Interval: settings.Interval,
	}
}
