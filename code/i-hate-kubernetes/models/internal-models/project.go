package models

import external_models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/external-models"

type Project struct {
	Project string //An identifier for this project used for grouping pods
	Engine  string //Which container engine to use, e.g. docker, podman, etc
	Pwd     string

	Logging      bool          //Enable logging
	Dashboard    bool          //Enable dashboard
	Analytics    bool          //Enable analytics
	Loadbalancer *LoadBalancer //Which loadbalancer to use
	Registry     *Registry     //Deploy an internal registry, can be reachable from the outside in order to not being forced to use docker image registry
	Autoupdate   *Cicd         //Autoupdate i-hate-kubernetes whenever something is pushed to master branch (branch to be updated to latest at one point)
	Cicd         []Cicd        //If true then we spin up a CICD (continous integration, continious deployment) pipeline

	Settings Settings            //Settings?
	Services map[string]*Service `yaml:",inline"` //A list of the services to deploy
}

func ParseProject(project external_models.Project, pwd string) Project {
	p := Project{
		Project: project.Project,
		Engine:  project.Engine,
		Pwd:     pwd, //TODO: Insert from the external project if it is there

		Logging:    project.Logging,
		Dashboard:  project.Dashboard,
		Analytics:  project.Analytics,
		Autoupdate: ParseAutoupdate(project.Autoupdate, pwd), //TODO: Fix this, because pwd here only makes sense in development

		Settings: ParseSettings(project.Settings),
	}
	p.Services = ParseServices(project.Services, p)
	p.Registry = ParseRegistry(project.Registry, p)
	p.Loadbalancer = ParseLoadBalancer(project.Loadbalancer, p)
	p.Cicd = ParseCicds(p.Services, project.Services, p.Pwd)

	return p
}

func (project *Project) GetId() *string {
	if project == nil {
		return nil
	}

	return &project.Project
}
