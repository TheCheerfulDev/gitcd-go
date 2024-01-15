package config

import (
	"fmt"
	"os"
	"path"
)

var gitCdHomePath, databaseFilePath, dirChangerPath, projectRootPath string

func Init() {

	lookupEnv, exists := os.LookupEnv("GITCD_PROJECT_HOME")

	if !exists {
		fmt.Println("GITCD_PROJECT_HOME environment variable is missing")
		os.Exit(1)
	}

	projectRootPath = lookupEnv
	homeDir, _ := os.UserHomeDir()
	gitCdHomePath = path.Join(homeDir, ".config", "gitcd")
	databaseFilePath = path.Join(gitCdHomePath, "gitcd.db")
	dirChangerPath = path.Join(gitCdHomePath, "change_dir.sh")

	if _, err := os.Stat(gitCdHomePath); os.IsNotExist(err) {
		err := os.MkdirAll(gitCdHomePath, 0755)
		if err != nil {
			fmt.Println("Unable to create gitcd configuration directory: ", gitCdHomePath)
		}
	}
}

func GetProjectRootPath() string {
	return projectRootPath
}

func GetDatabaseFilePath() string {
	return databaseFilePath
}

func GetDirChangerPath() string {
	return dirChangerPath
}
