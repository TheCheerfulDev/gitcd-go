package repository

import (
	"bufio"
	"fmt"
	"github.com/thecheerfuldev/gitcd-go/config"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
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

func GetProjectsRegex(input string) []string {
	projects := make([]Project, 0)
	compile, err := regexp.Compile(input)
	if err != nil {
		fmt.Println("Invalid regular expression.")
		os.Exit(1)
	}

	for key, project := range database {
		if compile.MatchString(key) {
			projects = append(projects, project)
		}
	}

	OrderedBy(countSorter, pathSorter).Sort(projects)

	result := make([]string, 0)

	for _, value := range projects {
		result = append(result, value.Path)
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

func Init(config config.Config) {
	cfg = config
	if _, err := os.Stat(cfg.DatabaseFilePath); os.IsNotExist(err) {
		create, _ := os.Create(cfg.DatabaseFilePath)
		err := create.Close()
		if err != nil {
			fmt.Println("Something went wrong while creating the gitcd database file: ", err)
			os.Exit(1)
		}
	}
	readDatabase()
}

func WriteChangesToDatabase() {
	if !isModified {
		return
	}

	databaseFile, err := os.OpenFile(cfg.DatabaseFilePath, os.O_TRUNC|os.O_WRONLY, os.ModePerm)
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

type lessFunc func(p1, p2 *Project) bool

type MultiSorter struct {
	projects []Project
	less     []lessFunc
}

func (ms *MultiSorter) Sort(projects []Project) {
	ms.projects = projects
	sort.Sort(ms)
}

func OrderedBy(less ...lessFunc) *MultiSorter {
	return &MultiSorter{
		less: less,
	}
}

// Len is part of sort.Interface.
func (ms *MultiSorter) Len() int {
	return len(ms.projects)
}

// Swap is part of sort.Interface.
func (ms *MultiSorter) Swap(i, j int) {
	ms.projects[i], ms.projects[j] = ms.projects[j], ms.projects[i]
}

func (ms *MultiSorter) Less(i, j int) bool {
	p, q := &ms.projects[i], &ms.projects[j]
	var k int
	for k = 0; k < len(ms.less)-1; k++ {
		less := ms.less[k]
		switch {
		case less(p, q):
			// p < q, so we have a decision.
			return true
		case less(q, p):
			// p > q, so we have a decision.
			return false
		}
	}
	return ms.less[k](q, p)
}

func GiveTopTen() []string {
	projects := make([]Project, 0)
	for _, project := range database {
		projects = append(projects, project)
	}

	OrderedBy(countSorter, pathSorter).Sort(projects)
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

func ResetDatabase() {
	database = map[string]Project{}
}
