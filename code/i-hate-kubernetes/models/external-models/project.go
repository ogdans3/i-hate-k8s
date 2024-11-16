package external_models

type Project struct {
	Project string
	Engine  string
	Pwd     string

	Logging      bool
	Registry     bool
	Dashboard    bool
	Analytics    bool
	Loadbalancer bool
	Cicd         bool
	Autoupdate   bool

	Settings Settings
	Services map[string]*Service `yaml:",inline"`
}

func (project *Project) InsertDefaults(directory string) {
	project.Pwd = directory
	for serviceName, service := range project.Services {
		service.InsertDefaults(serviceName)
	}
}
