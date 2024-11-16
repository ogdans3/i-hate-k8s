package yaml

import (
	"os"
	"path/filepath"

	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/console"
	external_models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/external-models"
	models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/internal-models"
	"gopkg.in/yaml.v3"
)

func ReadFile(file string) models.Project {
	data, err := os.ReadFile(file)
	if err != nil {
		console.Error("error: %v", err)
		panic("Unable to read file ")
	}

	//TODO: Set default values
	project := external_models.Project{}
	err2 := yaml.Unmarshal([]byte(data), &project)
	if err2 != nil {
		// TODO: Handle the error properly, inspect the yaml file and give proper errors
		console.Error("error: %v", err2)
	}
	project.InsertDefaults(filepath.Dir(file))

	cleanPath := filepath.Clean(file)
	firstDir, _ := filepath.Split(cleanPath)

	_, err = os.Getwd()
	if err != nil {
		console.InfoLog.Error(err)
		panic(err)
	}

	parsedProject := models.ParseProject(project, firstDir)

	//TODO: Validate the file, check that the project is there, etc.
	return parsedProject
}
