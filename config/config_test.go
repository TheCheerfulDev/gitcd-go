package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestDefault(t *testing.T) {
	projectRoot := t.TempDir()
	_ = os.Setenv("GITCD_PROJECT_HOME", projectRoot)

	cfg := Default()
	home, _ := os.UserHomeDir()
	got := cfg.GitCdHomePath
	want := path.Join(home, ".config", "gitcd")
	assert.Equal(t, got, want, "got %v, want %v", got, want)

	got = cfg.DatabaseFilePath
	want = path.Join(cfg.GitCdHomePath, "gitcd.db")
	assert.Equal(t, got, want, "got %v, want %v", got, want)

	got = cfg.ProjectRootPath
	want = projectRoot
	assert.Equal(t, got, want, "got %v, want %v", got, want)

	got = cfg.DirChangerPath
	want = path.Join(cfg.GitCdHomePath, "change_dir.sh")
	assert.Equal(t, got, want, "got %v, want %v", got, want)

}

func TestDefaultWithHomeDir(t *testing.T) {
	home, _ := os.UserHomeDir()
	_ = os.Unsetenv("GITCD_PROJECT_HOME")

	cfg := Default()
	got := cfg.ProjectRootPath
	want := home
	assert.Equal(t, got, want, "got %v, want %v", got, want)

}

func TestSetGet(t *testing.T) {
	cfg := Default()
	Set(cfg)
	got := Get()
	want := cfg
	assert.Equal(t, got, want, "got %v, want %v", got, want)
}

func TestInit(t *testing.T) {
	tempDir := path.Join(os.TempDir(), "gitcd")
	defer os.RemoveAll(tempDir)

	cfg := Config{
		GitCdHomePath: tempDir,
	}
	err := Init(cfg)

	assert.NoError(t, err, "got %v, want %v", err, nil)
	assert.DirExists(t, tempDir, "got %v, want %v", tempDir, "to exist")
}

func TestGet(t *testing.T) {
	cfg := Default()
	Set(cfg)
	got := Get()
	want := cfg
	assert.Equal(t, got, want, "got %v, want %v", got, want)
}
