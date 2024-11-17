package console

import (
	"log"
	"os"
)

func CreateLoggingDirectory() string {
	base, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	_, err = os.Stat(base)
	if os.IsNotExist(err) {
		log.Fatal("We do not create the base path for you. We only support operating systems with the "+base+" directory", err)
	}
	base = base + "/.logs/"
	logDir := base + "hive/"

	_, err = os.Stat(logDir)
	if os.IsNotExist(err) {
		err := os.MkdirAll(logDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
		return logDir
	}
	if err != nil {
		panic(err)
	}
	return logDir
}

var InfoLog = NewLog().
	AddFile("log.log", &LogDestination{maximumLogLevel: ERROR}).
	AddStd(&LogDestination{minimumLogLevel: DEBUG, flags: 0}).
	AddFile("error.log", &LogDestination{minimumLogLevel: ERROR})

var StatLog = NewLog().
	AddFile("stats.log", nil)
