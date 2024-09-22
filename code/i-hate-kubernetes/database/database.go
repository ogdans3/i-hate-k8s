package database

import (
	"encoding/json"
	"log"
	"os"
)

type ContainerTable struct {
	container
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
		log.Fatal("Cannot find database file")
	}
	err = json.Unmarshal(data, &database)
	if err != nil {
		log.Fatal("Malformed database file")
	}
}

func write() {
	str, err := json.Marshal(&database)
	if err != nil {
		log.Fatal("Unable to marshal the database")
	}
	err = os.WriteFile(getFile(), str, 0644)
	if err != nil {
		log.Fatal("Unable to store the database")
	}
}
