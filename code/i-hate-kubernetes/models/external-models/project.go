package external_models

type Project struct {
	Project      string
	Engine       string
	Logging      bool
	Registry     bool
	Dashboard    bool
	Analytics    bool
	Loadbalancer bool
	Settings     Settings
	Services     map[string]Service `yaml:",inline"`
}

func (project *Project) InsertDefaults() {
	for _, service := range project.Services {
		service.InsertDefaults()
	}
}
