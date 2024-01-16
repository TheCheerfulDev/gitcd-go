package config

import (
	"os"
	"path"
	"strings"
	"testing"
)

func TestDefault(t *testing.T) {

	cfg := Default()
	home, _ := os.UserHomeDir()

	got := cfg.gitCdHomePath
	want := path.Join(home, ".config", "gitcd")
	compare := strings.Compare(got, want)
	if compare != 0 {
		t.Errorf("got %v, want %v", got, want)
	}

	got = cfg.DatabaseFilePath
	want = path.Join(cfg.gitCdHomePath, "gitcd.db")
	compare = strings.Compare(got, want)
	if compare != 0 {
		t.Errorf("got %v, want %v", got, want)
	}

	got = cfg.ProjectRootPath

	lookupEnv, exists := os.LookupEnv("GITCD_PROJECT_HOME")
	if exists {
		want = lookupEnv
	} else {
		want = home
	}

	compare = strings.Compare(got, want)
	if compare != 0 {
		t.Errorf("got %v, want %v", got, want)
	}

	got = cfg.DirChangerPath
	want = path.Join(cfg.gitCdHomePath, "change_dir.sh")
	compare = strings.Compare(got, want)
	if compare != 0 {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestSetGet(t *testing.T) {
	cfg := Default()
	Set(cfg)
	got := Get()
	want := cfg
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestInit(t *testing.T) {
	tempDir := path.Join(os.TempDir(), "gitcd")
	defer os.RemoveAll(tempDir)

	cfg := Default()
	cfg.gitCdHomePath = tempDir
	err := Init(cfg)
	if err != nil {
		t.Errorf("got %v, want %v", err, nil)
	}
	if _, err := os.Stat(tempDir); err != nil {
		t.Errorf("got %v, want %v", err, nil)
	}

	got := Get()

	if got != cfg {
		t.Errorf("got %v, want %v", got, cfg)
	}
}
