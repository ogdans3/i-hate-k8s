package models

type Service struct {
	image string
	build string
	dev   string
	watch string

	www       bool
	https     bool
	ports     []string
	autosacle Autoscale
}
