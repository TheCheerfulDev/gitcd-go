package config

import (
	"errors"
	"os"
	"path/filepath"
)

type Config struct {
	GitCdHomePath, DatabaseFilePath, DirChangerPath, ProjectRootPath string
	CaseSensitive                                                    bool
}

var cfg Config

func Default() Config {
	c := Config{}
	homeDir, _ := os.UserHomeDir()
	lookupEnv, exists := os.LookupEnv("GITCD_PROJECT_HOME")

	if exists {
		c.ProjectRootPath = lookupEnv
	} else {
		c.ProjectRootPath = homeDir
	}

	lookupEnv, exists = os.LookupEnv("GITCD_CASE_SENSITIVE")
	if exists {
		c.CaseSensitive = lookupEnv == "true"
	} else {
		c.CaseSensitive = false
	}

	c.GitCdHomePath = filepath.Join(homeDir, ".config", "gitcd")
	c.DatabaseFilePath = filepath.Join(c.GitCdHomePath, "gitcd.db")
	c.DirChangerPath = filepath.Join(c.GitCdHomePath, "change_dir.sh")

	return c
}

func Set(c Config) {
	cfg = c
}

func Get() Config {
	return cfg
}

func Init(c Config) error {
	if _, err := os.Stat(c.GitCdHomePath); os.IsNotExist(err) {
		err := os.MkdirAll(c.GitCdHomePath, 0755)
		if err != nil {
			return errors.New("unable to create gitcd configuration directory " + c.GitCdHomePath)
		}
	}
	Set(c)
	return nil
}
