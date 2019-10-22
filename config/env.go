package config

import "os"

type environment struct {
	goroot      string
	gopath      string
	go111module string
}

var env *environment

func init() {
	goroot := os.Getenv("GOROOT")
	gopath := os.Getenv("GOPATH")
	go111module := os.Getenv("GO111MODULE")

	env = &environment{
		goroot:      goroot,
		gopath:      gopath,
		go111module: go111module,
	}
}

func GetEnv() *environment {
	return env
}

func (e *environment) GetGoRoot() string {
	return e.goroot
}

func (e *environment) GetGoPath() string {
	return e.gopath
}

func (e *environment) GetGo111Module() string {
	return e.go111module
}
