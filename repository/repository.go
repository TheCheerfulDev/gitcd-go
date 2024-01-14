package repository

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var database = map[string]Project{}
var dbFilePath = "/tmp/gitcd.db"
var isModified = false

type Project struct {
	Name        string
	Path        string
	CallCounter int
}

func (project *Project) UpdateCounter() {
	project.CallCounter += 1
	SaveProject(*project)
}

func (project Project) Display() {
	fmt.Printf("Name: %v Path:%v Called: %v\n", project.Name, project.Path, project.CallCounter)
}

func (project *Project) saveString() string {
	return fmt.Sprintf("%v;%v;%v", project.Name, project.Path, project.CallCounter)
}

func AddProjectFromDb(name, path string, callCount int) *Project {
	project := Project{
		Name:        name,
		Path:        path,
		CallCounter: callCount,
	}

	database[path] = project
	return &project
}

func AddProject(name, path string, callCount int) *Project {
	project, exists := database[path]

	if exists {
		return &project
	}

	project = Project{
		Name:        name,
		Path:        path,
		CallCounter: callCount,
	}

	database[path] = project
	isModified = true
	return &project
}

func GetProjectContaining(input string) []string {
	result := make([]string, 0)

	for key := range database {
		if strings.Contains(key, input) {
			result = append(result, key)
		}
	}
	return result
}

func SaveProject(project Project) {
	database[project.Path] = project
	isModified = true
}

func GetProject(key string) Project {
	project := database[key]
	return project
}

func readDatabase() {
	dbFile, _ := os.Open(dbFilePath)
	defer dbFile.Close()

	var lines []string
	scanner := bufio.NewScanner(dbFile)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	for _, projectText := range lines {
		split := strings.Split(projectText, ";")
		i, _ := strconv.ParseInt(split[2], 10, 0)
		callCount := int(i)
		if split[1] == "" {
			continue
		}
		AddProjectFromDb(split[0], split[1], callCount)
	}
}

func Init() {
	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {
		create, _ := os.Create(dbFilePath)
		create.Close()
	}
	readDatabase()
}

func WriteChangesToDatabase() {
	if !isModified {
		return
	}

	fmt.Println("DB modified, writing to disk")

	databaseFile, err := os.OpenFile(dbFilePath, os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
	defer databaseFile.Close()

	for _, project := range database {
		_, err := databaseFile.WriteString(project.saveString() + "\n")
		if err != nil {
			fmt.Println(err)
		}
	}
}
