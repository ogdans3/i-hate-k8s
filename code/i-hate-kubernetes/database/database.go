package database

import (
	"encoding/json"
	"os"

	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/console"
)

type ContainerTable struct {
}

type Database struct {
	ContainerTable ContainerTable
}

var database *Database = nil

func getFile() string {
	pwd, _ := os.Getwd()
	return pwd + "/.database.json"
}

func read() {
	data, err := os.ReadFile(getFile())
	if err != nil {
		console.Fatal("Cannot find database file")
	}
	err = json.Unmarshal(data, &database)
	if err != nil {
		console.Fatal("Malformed database file")
	}
}

func write() {
	str, err := json.Marshal(&database)
	if err != nil {
		console.Fatal("Unable to marshal the database")
	}
	err = os.WriteFile(getFile(), str, 0644)
	if err != nil {
		console.Fatal("Unable to store the database")
	}
}
