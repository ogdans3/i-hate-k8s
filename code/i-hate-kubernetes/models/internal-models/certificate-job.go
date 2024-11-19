package models

import (
	"path/filepath"

	external_models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/external-models"
)

type CertificateJob struct {
	Service *Service //A service object for the certification job
}

func ParseCertificateJob(project Project) *CertificateJob {
	job := &CertificateJob{
		Service: ParseService(&external_models.Service{
			ServiceName:   "hive-certs",
			ContainerName: "hive-certs",
			Image:         "hive-certs:latest",
			Ports: []string{
				"80",
			},
			Autoscale: external_models.Autoscale{Initial: 1, Autoscale: false}, //Can certbot autoscale? We probably need to ensure that domains go to the correct instance
			Probes: &external_models.Probes{
				Ready:    "/ready",
				Liveness: "/live",
			},

			Build:      true,
			Dockerfile: "Dockerfile",
		}, project),
	}
	job.Service.Directory = filepath.Join(project.ProgramDirectory, "../services/tls/") //TODO: This is a hack until we start to upload these images to repo
	return job
}
