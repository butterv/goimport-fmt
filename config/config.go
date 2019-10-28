package config

type config struct {
	goroot     string
	filePath   string
	ownProject string
}

var c *config

// Set sets the setting value
func Set(goroot, filePath, ownProject string) {
	if c == nil {
		c = &config{
			goroot:     goroot,
			filePath:   filePath,
			ownProject: ownProject,
		}
	}
}

// GetGoRoot gets a GOROOT
func GetGoRoot() string {
	return c.goroot
}

// GetFilePath gets a target file path
func GetFilePath() string {
	return c.filePath
}

// GetOwnProject gets a own project
func GetOwnProject() string {
	return c.ownProject
}
