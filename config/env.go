package config

import (
	"flag"
	"os"
)

type environment struct {
	goroot     string
	filePath   string
	ownProject string
}

var env *environment

func init() {
	filePathPtr := flag.String("filepath", "", "file path")
	ownProjectPtr := flag.String("ownproject", "", "own project")
	flag.Parse()

	filePath := *filePathPtr
	if filePath == "" {
		panic("file path not found")
	}
	ownProject := *ownProjectPtr
	if ownProject == "" {
		panic("own project not found")
	}

	goroot := os.Getenv("GOROOT")

	env = &environment{
		goroot:     goroot,
		filePath:   filePath,
		ownProject: ownProject,
	}
}

func GetEnv() *environment {
	return env
}

func (e *environment) GetGoRoot() string {
	return e.goroot
}

func (e *environment) GetFilePath() string {
	return e.filePath
}

func (e *environment) GetOwnProject() string {
	return e.ownProject
}
