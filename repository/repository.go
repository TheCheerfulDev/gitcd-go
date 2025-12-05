package repository

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/thecheerfuldev/gitcd-go/config"
)

var database = map[string]Project{}
var isModified = false
var cfg config.Config

var countSorter = func(c1, c2 *Project) bool {
	return c1.CallCounter > c2.CallCounter
}
var pathSorter = func(c1, c2 *Project) bool {
	return c1.Path > c2.Path
}

type Project struct {
	Path        string
	CallCounter int
}

func (project *Project) UpdateCounter() {
	project.CallCounter += 1
	SaveProject(*project)
}

func (project *Project) saveString() string {
	return fmt.Sprintf("%v;%v", project.Path, project.CallCounter)
}

func addProjectFromDb(path string, callCount int) {
	project := Project{
		Path:        path,
		CallCounter: callCount,
	}

	database[path] = project
}

func AddProject(path string) {
	project, exists := database[path]

	if exists {
		return
	}

	project = Project{
		Path:        path,
		CallCounter: 0,
	}

	database[path] = project
	isModified = true
}

func GetAllProjects() []string {
	result := make([]string, len(database))

	index := 0
	for key := range database {
		result[index] = key
		index++
	}

	return result
}

func RemoveProject(key string) {
	delete(database, key)
	isModified = true
}

func GetProjectsRegex(input string) ([]string, error) {
	projects := make([]Project, 0)
	if caseInsensitive() {
		input = strings.ToLower(input)
	}
	compile, err := regexp.Compile(input)
	if err != nil {
		return nil, errors.New("Invalid regular expression")
	}

	for key, project := range database {
		if caseInsensitive() {
			key = strings.ToLower(key)
		}
		if compile.MatchString(key) {
			projects = append(projects, project)
		}
	}

	sort.Slice(projects, func(i, j int) bool {
		if projects[i].CallCounter != projects[j].CallCounter {
			return projects[i].CallCounter > projects[j].CallCounter
		}
		// CallCounters are equal, sort by Path, alphabetically, hence the "<"
		return projects[i].Path < projects[j].Path
	})

	result := make([]string, 0)

	for _, value := range projects {
		result = append(result, value.Path)
	}
	return result, nil
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
	dbFile, _ := os.Open(cfg.DatabaseFilePath)
	defer dbFile.Close()

	var lines []string
	scanner := bufio.NewScanner(dbFile)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	for _, projectText := range lines {

		if projectText == "" {
			return
		}

		split := strings.Split(projectText, ";")
		i, _ := strconv.ParseInt(split[1], 10, 0)
		callCount := int(i)
		if split[1] == "" {
			continue
		}
		addProjectFromDb(split[0], callCount)
	}
}

func Init(config config.Config) error {
	cfg = config
	if _, err := os.Stat(cfg.DatabaseFilePath); os.IsNotExist(err) {
		create, err := os.Create(cfg.DatabaseFilePath)
		if err != nil {
			fmt.Println("Something went wrong while creating the gitcd database file: ", err)
			return err
		}

		err = create.Close()
		if err != nil {
			fmt.Println("Something went wrong while creating the gitcd database file: ", err)
			return err
		}
	}
	readDatabase()
	return nil
}

func WriteChangesToDatabase() {
	if !isModified {
		return
	}

	databaseFile, err := os.OpenFile(cfg.DatabaseFilePath, os.O_TRUNC|os.O_WRONLY, 0644)
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

func GiveTopTen() []string {
	projects := make([]Project, 0)
	for _, project := range database {
		projects = append(projects, project)
	}

	sort.Slice(projects, func(i, j int) bool {
		if projects[i].CallCounter != projects[j].CallCounter {
			return projects[i].CallCounter > projects[j].CallCounter
		}
		// CallCounters are equal, sort by Path, alphabetically, hence the "<"
		return projects[i].Path < projects[j].Path
	})

	maxSize := 10

	if len(projects) < maxSize {
		maxSize = len(projects)
	}

	topTenProjects := make([]string, maxSize)

	for i, project := range projects[:maxSize] {
		topTenProjects[i] = project.Path
	}

	return topTenProjects

}

func caseInsensitive() bool {
	return !cfg.CaseSensitive
}

func ResetDatabase() {
	database = map[string]Project{}
	isModified = true
}
