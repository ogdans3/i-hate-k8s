package models

import (
	"os"
	"sort"

	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/console"
	external_models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/external-models"
)

type Project struct {
	Project          string //An identifier for this project used for grouping pods
	Engine           string //Which container engine to use, e.g. docker, podman, etc
	Pwd              string
	ProgramDirectory string

	Logging      bool          //Enable logging
	Dashboard    bool          //Enable dashboard
	Analytics    bool          //Enable analytics
	Loadbalancer *LoadBalancer //Which loadbalancer to use
	Registry     *Registry     //Deploy an internal registry, can be reachable from the outside in order to not being forced to use docker image registry
	Autoupdate   *Cicd         //Autoupdate i-hate-kubernetes whenever something is pushed to master branch (branch to be updated to latest at one point)
	Cicd         []Cicd        //If true then we spin up a CICD (continous integration, continious deployment) pipeline

	Settings           Settings            //Settings?
	Services           map[string]*Service `yaml:",inline"` //A list of the services to deploy
	CertificateHandler *CertificateHandler //TODO: Comment
}

func ParseProject(project external_models.Project, pwd string) Project {
	programDirectory, err := os.Getwd()
	if err != nil {
		console.InfoLog.Panic("Unable to find program directory")
	}
	p := Project{
		Project:          project.Project,
		Engine:           project.Engine,
		Pwd:              pwd,              //TODO: Insert from the external project if it is there
		ProgramDirectory: programDirectory, //TODO: Insert from the external project if it is there

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
	p.CertificateHandler = ParseCertificateBlocks(p, p.Services)
	return p
}

func (project *Project) GetLoadbalancedServices() []*Service {
	keys := make([]string, 0, len(project.Services))
	for key := range project.Services {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	services := make([]*Service, 0)
	if project.CertificateHandler != nil && project.CertificateHandler.ServiceJob != nil && project.CertificateHandler.ServiceJob.Service != nil {
		services = append(services, project.CertificateHandler.ServiceJob.Service)
	}
	for _, key := range keys {
		service := project.Services[key]
		services = append(services, service)
	}
	return services
}

func (project *Project) GetId() *string {
	if project == nil {
		return nil
	}

	return &project.Project
}
