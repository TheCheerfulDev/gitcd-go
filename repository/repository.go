package repository

import (
	"bufio"
	"fmt"
	"github.com/thecheerfuldev/gitcd-go/config"
	"os"
	"sort"
	"strconv"
	"strings"
)

var database = map[string]Project{}
var dbFilePath string
var isModified = false

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

func AddProjectFromDb(path string, callCount int) *Project {
	project := Project{
		Path:        path,
		CallCounter: callCount,
	}

	database[path] = project
	return &project
}

func AddProject(path string, callCount int) *Project {
	project, exists := database[path]

	if exists {
		return &project
	}

	project = Project{
		Path:        path,
		CallCounter: callCount,
	}

	database[path] = project
	isModified = true
	return &project
}

func GetAllProjects() []string {
	//result := make([]string, 0)
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
		i, _ := strconv.ParseInt(split[1], 10, 0)
		callCount := int(i)
		if split[1] == "" {
			continue
		}
		AddProjectFromDb(split[0], callCount)
	}
}

func Init() {
	dbFilePath = config.GetDatabaseFilePath()
	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {
		create, _ := os.Create(dbFilePath)
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

// Less is part of sort.Interface. It is implemented by looping along the
// less functions until it finds a comparison that discriminates between
// the two items (one is less than the other). Note that it can call the
// less functions twice per call. We could change the functions to return
// -1, 0, 1 and reduce the number of calls for greater efficiency: an
// exercise for the reader.
func (ms *MultiSorter) Less(i, j int) bool {
	p, q := &ms.projects[i], &ms.projects[j]
	// Try all but the last comparison.
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
		// p == q; try the next comparison.
	}
	// All comparisons to here said "equal", so just return whatever
	// the final comparison reports.
	return ms.less[k](p, q)
}

func GiveTopTen() []string {
	projects := make([]Project, 0)
	for _, project := range database {
		projects = append(projects, project)
	}

	// Closures that order the Change structure.
	count := func(c1, c2 *Project) bool {
		return c1.CallCounter > c2.CallCounter
	}
	path := func(c1, c2 *Project) bool {
		return c1.Path < c2.Path
	}

	// Simple use: Sort by count.
	OrderedBy(count, path).Sort(projects)

	topTenProjects := make([]string, 10)
	for i, project := range projects[:10] {
		topTenProjects[i] = project.Path
	}

	return topTenProjects

}
