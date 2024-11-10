package yaml

import (
	"log"
	"os"

	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/console"
	external_models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/external-models"
	models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/internal-models"
	"gopkg.in/yaml.v3"
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
	console.Log("Before insert project: ", project)
	project.InsertDefaults()
	console.Log("After insert project: ", project)
	parsedProject := models.ParseProject(project)

	//TODO: Validate the file, check that the project is there, etc.
	return parsedProject
}
