package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"testing"
)

func TestDefault(t *testing.T) {
	projectRoot := t.TempDir()
	_ = os.Setenv("GITCD_PROJECT_HOME", projectRoot)

	cfg := Default()
	home, _ := os.UserHomeDir()
	actual := cfg.GitCdHomePath
	expected := path.Join(home, ".config", "gitcd")
	assert.Equal(t, expected, actual, "actual %v, expected %v", actual, expected)

	actual = cfg.DatabaseFilePath
	expected = path.Join(cfg.GitCdHomePath, "gitcd.db")
	assert.Equal(t, expected, actual, "actual %v, expected %v", actual, expected)

	actual = cfg.ProjectRootPath
	expected = projectRoot
	assert.Equal(t, expected, actual, "actual %v, expected %v", actual, expected)

	actual = cfg.DirChangerPath
	expected = path.Join(cfg.GitCdHomePath, "change_dir.sh")
	assert.Equal(t, expected, actual, "actual %v, expected %v", actual, expected)

}

func TestDefaultWithHomeDir(t *testing.T) {
	home, _ := os.UserHomeDir()
	_ = os.Unsetenv("GITCD_PROJECT_HOME")

	cfg := Default()
	actual := cfg.ProjectRootPath
	expected := home
	assert.Equal(t, expected, actual, "actual %v, expected %v", actual, expected)

}

func TestDefaultWithCaseSensitive(t *testing.T) {
	_ = os.Setenv("GITCD_CASE_SENSITIVE", "true")

	cfg := Default()
	actual := cfg.CaseSensitive
	expected := true
	assert.Equal(t, expected, actual, "actual %v, expected %v", actual, expected)
}

func TestDefaultWithCaseInsensitive(t *testing.T) {
	_ = os.Unsetenv("GITCD_CASE_SENSITIVE")

	cfg := Default()
	actual := cfg.CaseSensitive
	expected := false
	assert.Equal(t, expected, actual, "actual %v, expected %v", actual, expected)
}

func TestSetGet(t *testing.T) {
	cfg := Default()
	Set(cfg)
	actual := Get()
	expected := cfg
	assert.Equal(t, expected, actual, "actual %v, expected %v", actual, expected)
}

func TestInit(t *testing.T) {
	tempDir := path.Join(os.TempDir(), "gitcd")
	defer os.RemoveAll(tempDir)

	cfg := Config{
		GitCdHomePath: tempDir,
	}
	err := Init(cfg)

	require.NoError(t, err, "actual %v, expected %v", err, nil)
	assert.DirExists(t, tempDir, "expected %v to exist", tempDir)
}

func TestGet(t *testing.T) {
	cfg := Default()
	Set(cfg)
	actual := Get()
	expected := cfg
	assert.Equal(t, expected, actual, "actual %v, expected %v", actual, expected)
}
