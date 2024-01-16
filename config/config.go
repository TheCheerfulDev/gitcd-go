package config

import (
	"errors"
	"os"
	"path"
)

type Config struct {
	gitCdHomePath, DatabaseFilePath, DirChangerPath, ProjectRootPath string
}

var cfg *Config

func Default() *Config {
	c := &Config{}
	homeDir, _ := os.UserHomeDir()
	lookupEnv, exists := os.LookupEnv("GITCD_PROJECT_HOME")

	if exists {
		c.ProjectRootPath = lookupEnv
	} else {
		c.ProjectRootPath = homeDir
	}

	c.gitCdHomePath = path.Join(homeDir, ".config", "gitcd")
	c.DatabaseFilePath = path.Join(c.gitCdHomePath, "gitcd.db")
	c.DirChangerPath = path.Join(c.gitCdHomePath, "change_dir.sh")

	return c
}

func Set(c *Config) {
	cfg = c
}

func Get() *Config {
	return cfg
}

func Init(c *Config) error {
	if _, err := os.Stat(c.gitCdHomePath); os.IsNotExist(err) {
		err := os.MkdirAll(c.gitCdHomePath, 0755)
		if err != nil {
			//fmt.Println("Unable to create gitcd configuration directory: ", c.gitCdHomePath)
			return errors.New("unable to create gitcd configuration directory " + c.gitCdHomePath)
		}
	}
	Set(c)
	return nil
}
