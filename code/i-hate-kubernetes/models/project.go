package models

type Project struct {
	Project      string
	Engine       string
	Logging      bool
	Dashboard    bool
	Analytics    bool
	Loadbalancer bool
	Services     map[string]Service `yaml:",inline"`
}
