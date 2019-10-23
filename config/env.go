package config

import "os"

type environment struct {
	goroot     string
	ownProject string
}

var env *environment

func Setup(ownProject string) {
	goroot := os.Getenv("GOROOT")

	env = &environment{
		goroot:     goroot,
		ownProject: ownProject,
	}
}

func GetEnv() *environment {
	return env
}

func (e *environment) GetGoRoot() string {
	return e.goroot
}

func (e *environment) GetOwnProject() string {
	return e.ownProject
}
