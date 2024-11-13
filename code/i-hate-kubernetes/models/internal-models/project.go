package models

import external_models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/external-models"

type Project struct {
	Project      string              //An identifier for this project used for grouping pods
	Engine       string              //Which container engine to use, e.g. docker, podman, etc
	Logging      bool                //Enable logging
	Registry     bool                //Deploy an internal registry, can be reachable from the outside in order to not being forced to use docker image registry
	Dashboard    bool                //Enable dashboard
	Analytics    bool                //Enable analytics
	Loadbalancer *LoadBalancer       //Which loadbalancer to use
	Settings     Settings            //Settings?
	Services     map[string]*Service `yaml:",inline"` //A list of the services to deploy
}

func ParseProject(project external_models.Project) Project {
	return Project{
		Project:      project.Project,
		Engine:       project.Engine,
		Logging:      project.Logging,
		Registry:     project.Registry,
		Dashboard:    project.Dashboard,
		Analytics:    project.Analytics,
		Loadbalancer: ParseLoadBalancer(project.Loadbalancer, project.Project),
		Settings:     ParseSettings(project.Settings),
		Services:     ParseServices(project.Services, project.Project),
	}
}

func (project *Project) GetId() *string {
	if project == nil {
		return nil
	}

	return &project.Project
}
