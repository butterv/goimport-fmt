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

// GetGoRoot gets a GOROOT
func GetGoRoot() string {
	return env.goroot
}

// GetFilePath gets a target file path
func GetFilePath() string {
	return env.filePath
}

// GetOwnProject gets a own project
func GetOwnProject() string {
	return env.ownProject
}
