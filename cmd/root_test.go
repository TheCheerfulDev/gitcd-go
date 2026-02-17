package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thecheerfuldev/gitcd-go/config"
	"github.com/thecheerfuldev/gitcd-go/repository"
)

func TestExtractExpression_regex(t *testing.T) {
	input := ".*test.*"
	expected := ".*test.*"
	actual := extractExpression([]string{input})
	assert.Equal(t, expected, actual, "Regex should stay unchanged")
}

func TestExtractExpression_single(t *testing.T) {
	input := "test"
	expected := "test"
	actual := extractExpression([]string{input})
	assert.Equal(t, expected, actual, "Single argument should stay unchanged")
}

func TestExtractExpression_multi(t *testing.T) {
	input := []string{"foo", "bar", "baz"}
	expected := "foo.*bar.*baz"
	actual := extractExpression(input)
	assert.Equal(t, expected, actual, "Regex should stay unchanged")
}

func TestHandleSingleMatch(t *testing.T) {
	initTest(t)
	// write DB to correct path

}

func initTest(t *testing.T) {
	config.Set(config.Config{
		GitCdHomePath:    t.TempDir(),
		DatabaseFilePath: t.TempDir(),
		DirChangerPath:   t.TempDir(),
		ProjectRootPath:  t.TempDir(),
		CaseSensitive:    false,
	})
	_ = repository.Init(config.Get())
}
