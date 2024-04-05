package repository

import (
	"github.com/stretchr/testify/assert"
	"github.com/thecheerfuldev/gitcd-go/config"
	"os"
	"path/filepath"
	"testing"
)

func TestAddProject(t *testing.T) {
	initRepositoryTest(t)
	path := "/test/path/to/project"

	AddProject(path)

	assert.Len(t, database, 1, "Expected database to have 1 entry")
	assert.Equal(t, 0, database[path].CallCounter, "Expected call counter to be 0")
	assert.Equal(t, path, database[path].Path, "Expected path to be '%s'", path)

	AddProject(path)
	assert.Len(t, database, 1, "Expected database to have 1 entry, since the project already exists")
	assert.True(t, isModified, "Expected isModified to be true")
}

func TestUpdateCounter(t *testing.T) {
	initRepositoryTest(t)
	path := "/test/path/to/project"

	AddProject(path)

	project := database[path]

	project.UpdateCounter()

	assert.Equal(t, 1, project.CallCounter, "Expected call counter to be 1")
	assert.Equal(t, project.CallCounter, database[path].CallCounter, "Expected call counter to be 1")
	assert.True(t, isModified, "Expected isModified to be true")
}

func TestSaveProject(t *testing.T) {
	initRepositoryTest(t)
	path := "/test/path/to/project"

	AddProject(path)

	project := database[path]

	project.CallCounter = 42
	SaveProject(project)

	assert.Equal(t, 42, database[path].CallCounter, "Expected call counter to be 42")
	assert.True(t, isModified, "Expected isModified to be true")
}

func TestSaveString(t *testing.T) {
	initRepositoryTest(t)
	path := "/test/path/to/project"

	AddProject(path)

	project := database[path]

	assert.Equal(t, "/test/path/to/project;0", project.saveString(), "Expected save string to be '/test/path/to/project;0'")
}

func TestAddProjectFromDb(t *testing.T) {
	initRepositoryTest(t)
	path := "/test/path/to/project"

	addProjectFromDb(path, 42)

	assert.Equal(t, 42, database[path].CallCounter, "Expected call counter to be 42")
	assert.False(t, isModified, "Expected isModified to be false")
}

func TestGetAllProjects(t *testing.T) {
	initRepositoryTest(t)
	path := "/test/path/to/project"
	path2 := "/test/path/to/another/project"

	AddProject(path)
	AddProject(path2)

	projects := GetAllProjects()

	assert.Len(t, projects, 2, "Expected to have 2 project")
}

func TestRemoveProject(t *testing.T) {
	initRepositoryTest(t)
	path := "/test/path/to/project"
	path2 := "/test/path/to/another/project"

	AddProject(path)
	AddProject(path2)

	RemoveProject(path)

	assert.Len(t, database, 1, "Expected database to have 1 entry")
	assert.Equal(t, path2, GetAllProjects()[0], "Expected path to be '%s'", path2)
	assert.True(t, isModified, "Expected isModified to be true")
}

func TestGetProjectsRegex(t *testing.T) {
	initRepositoryTest(t)
	path := "/test/path/to/project"
	path2 := "/test/path/to/another/project"

	AddProject(path)
	AddProject(path2)

	projects, _ := GetProjectsRegex(".*another.*")

	assert.Len(t, projects, 1, "Expected to have 1 project")
	assert.Equal(t, path2, projects[0], "Expected path to be '%s'", path2)
}

func TestGetProjectsRegexEmpty(t *testing.T) {
	initRepositoryTest(t)
	path := "/test/path/to/project"
	path2 := "/test/path/to/another/project"

	AddProject(path)
	AddProject(path2)

	projects, _ := GetProjectsRegex(".*notfound.*")

	assert.Len(t, projects, 0, "Expected to have 0 project")
}

func TestGetProjectsRegexInvalidRegex(t *testing.T) {
	initRepositoryTest(t)
	path := "/test/path/to/project"
	path2 := "/test/path/to/another/project"

	AddProject(path)
	AddProject(path2)

	_, err := GetProjectsRegex(".*(")

	assert.EqualError(t, err, "Invalid regular expression")

}

func TestGetProject(t *testing.T) {
	initRepositoryTest(t)
	path := "/test/path/to/project"

	AddProject(path)

	project := GetProject(path)

	assert.Equal(t, path, project.Path, "Expected path to be '%s'", path)
}

func TestReadDatabase(t *testing.T) {
	initRepositoryTest(t)

	project1 := Project{
		Path:        "/test/path/to/project",
		CallCounter: 42,
	}
	project2 := Project{
		Path:        "/test/path/to/another/project",
		CallCounter: 23,
	}

	_ = os.WriteFile(config.Get().DatabaseFilePath, []byte(project1.saveString()+"\n"+project2.saveString()), 0644)

	readDatabase()

	assert.Len(t, database, 2, "Expected database to have 2 entries")

}

func TestWriteChangesToDatabase(t *testing.T) {
	initRepositoryTest(t)

	path1 := "/test/path/to/project"
	path2 := "/test/path/to/another/project"

	project1 := Project{
		Path:        path1,
		CallCounter: 42,
	}
	project2 := Project{
		Path:        path2,
		CallCounter: 23,
	}

	AddProject(path1)
	AddProject(path2)

	database = map[string]Project{
		project1.Path: project1,
		project2.Path: project2,
	}

	WriteChangesToDatabase()

	content, _ := os.ReadFile(config.Get().DatabaseFilePath)
	assert.Contains(t, string(content), project1.saveString(), "Expected database to contain '%s'", project1.saveString())
	assert.Contains(t, string(content), project2.saveString(), "Expected database to contain '%s'", project2.saveString())

}

func TestGiveTopTen(t *testing.T) {
	initRepositoryTest(t)

	project1 := Project{
		Path:        "/test/path/to/project",
		CallCounter: 42,
	}
	project2 := Project{
		Path:        "/test/path/to/another/project",
		CallCounter: 42,
	}
	project3 := Project{
		Path:        "/test/path/to/yet/another/project",
		CallCounter: 13,
	}
	project4 := Project{
		Path:        "/test/path/to/one/more/project",
		CallCounter: 3,
	}
	project5 := Project{
		Path:        "/test/path/to/last/project",
		CallCounter: 1,
	}

	database = map[string]Project{
		project1.Path: project1,
		project2.Path: project2,
		project3.Path: project3,
		project4.Path: project4,
		project5.Path: project5,
	}

	projects := GiveTopTen()

	assert.Len(t, projects, 5, "Expected to have 5 projects")
	assert.Equal(t, project1.Path, projects[1], "Expected path to be '%s'", project1.Path)
	assert.Equal(t, project2.Path, projects[0], "Expected path to be '%s'", project2.Path)
	assert.Equal(t, project3.Path, projects[2], "Expected path to be '%s'", project3.Path)
	assert.Equal(t, project4.Path, projects[3], "Expected path to be '%s'", project4.Path)
	assert.Equal(t, project5.Path, projects[4], "Expected path to be '%s'", project5.Path)

}

func TestResetDatabase(t *testing.T) {
	initRepositoryTest(t)

	project1 := Project{
		Path:        "/test/path/to/project",
		CallCounter: 42,
	}
	project2 := Project{
		Path:        "/test/path/to/another/project",
		CallCounter: 42,
	}
	project3 := Project{
		Path:        "/test/path/to/yet/another/project",
		CallCounter: 13,
	}
	project4 := Project{
		Path:        "/test/path/to/one/more/project",
		CallCounter: 3,
	}
	project5 := Project{
		Path:        "/test/path/to/last/project",
		CallCounter: 1,
	}

	database = map[string]Project{
		project1.Path: project1,
		project2.Path: project2,
		project3.Path: project3,
		project4.Path: project4,
		project5.Path: project5,
	}

	ResetDatabase()

	assert.Len(t, database, 0, "Expected database to be empty")
	assert.True(t, isModified, "Expected isModified to be true")

}

func TestCaseInsensitive(t *testing.T) {
	initRepositoryTest(t)
	path := "/test/path/to/project"

	AddProject(path)

	projects, _ := GetProjectsRegex(".*PROJECT.*")

	assert.Len(t, projects, 1, "Expected to have 1 project")
	assert.Equal(t, path, projects[0], "Expected path to be '%s'", path)
}

func initRepositoryTest(t *testing.T) {
	tempDir := t.TempDir()
	c := config.Config{
		GitCdHomePath:    tempDir,
		DatabaseFilePath: filepath.Join(tempDir, "gitcd.db"),
		DirChangerPath:   filepath.Join(tempDir, "change_dir.sh"),
	}
	_ = config.Init(c)
	Init(c)
	database = make(map[string]Project)
	isModified = false
}
