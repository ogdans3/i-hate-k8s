package models

import external_models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/external-models"

type Project struct {
	Project      string
	Engine       string
	Logging      bool
	Dashboard    bool
	Analytics    bool
	Loadbalancer bool
	Settings     Settings
	Services     map[string]Service `yaml:",inline"`
}

func ParseProject(project external_models.Project) Project {
	return Project{
		Project:      project.Project,
		Engine:       project.Engine,
		Logging:      project.Logging,
		Dashboard:    project.Dashboard,
		Analytics:    project.Analytics,
		Loadbalancer: project.Loadbalancer,
		Settings:     ParseSettings(project.Settings),
		Services:     ParseServices(project.Services),
	}
}
