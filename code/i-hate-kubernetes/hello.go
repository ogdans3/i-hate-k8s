package main

import (
	"fmt"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/container-interface/docker"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

func readFile(file string) models.Project {
	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	//TODO: Set default values
	project := models.Project{}
	err2 := yaml.Unmarshal([]byte(data), &project)
	if err2 != nil {
		// TODO: Handle the error properly, inspect the yaml file and give proper errors
		log.Fatalf("error: %v", err2)
	}

	//TODO: Validate the file, check that the project is there, etc.
	return project
}

func main() {
	pwd, _ := os.Getwd()
	spec := readFile(pwd + "/examples/hello-world.yml")

	fmt.Printf("%v", spec.Project)
	fmt.Printf("--- t:\n%v\n\n", spec)

	docker.ListAllContainers()
}
