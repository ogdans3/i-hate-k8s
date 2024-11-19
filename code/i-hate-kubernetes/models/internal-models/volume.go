package models

import external_models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/external-models"

type Volume struct {
	Id       string
	Name     string
	Target   string
	Type     string
	Readonly bool
}

func ParseVolume(volume *external_models.Volume) *Volume {
	return &Volume{
		Id:       volume.Id,
		Name:     volume.Name,
		Target:   volume.Target,
		Type:     volume.Type,
		Readonly: volume.Readonly,
	}
}

func ParseVolumes(vols []*external_models.Volume) []*Volume {
	volumes := make([]*Volume, 0)
	for _, vol := range vols {
		volumes = append(volumes, ParseVolume(vol))
	}
	return volumes
}
