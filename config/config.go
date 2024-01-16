package config

import (
	"fmt"
	"os"
	"path"
)

var gitCdHomePath, databaseFilePath, dirChangerPath, projectRootPath string

func Init() {
	homeDir, _ := os.UserHomeDir()
	lookupEnv, exists := os.LookupEnv("GITCD_PROJECT_HOME")

	if exists {
		projectRootPath = lookupEnv
	} else {
		projectRootPath = homeDir
	}

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
