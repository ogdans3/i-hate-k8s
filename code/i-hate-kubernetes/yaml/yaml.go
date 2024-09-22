package yaml

import (
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/external-models"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

func ReadFile(file string) models.Project {
	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	//TODO: Set default values
	project := external_models.Project{}
	err2 := yaml.Unmarshal([]byte(data), &project)
	if err2 != nil {
		// TODO: Handle the error properly, inspect the yaml file and give proper errors
		log.Fatalf("error: %v", err2)
	}
	project.InsertDefaults()
	parsedProject := models.ParseProject(project)

	//TODO: Validate the file, check that the project is there, etc.
	return parsedProject
}
