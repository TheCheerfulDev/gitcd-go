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

func TestValidateChoice_ValidFirstOption(t *testing.T) {
	index, valid := validateChoice("1", 5)
	assert.True(t, valid)
	assert.Equal(t, 0, index)
}

func TestValidateChoice_ValidLastOption(t *testing.T) {
	index, valid := validateChoice("5", 5)
	assert.True(t, valid)
	assert.Equal(t, 4, index)
}

func TestValidateChoice_ValidMiddleOption(t *testing.T) {
	index, valid := validateChoice("3", 5)
	assert.True(t, valid)
	assert.Equal(t, 2, index)
}

func TestValidateChoice_ZeroInvalid(t *testing.T) {
	_, valid := validateChoice("0", 5)
	assert.False(t, valid)
}

func TestValidateChoice_NegativeInvalid(t *testing.T) {
	_, valid := validateChoice("-1", 5)
	assert.False(t, valid)
}

func TestValidateChoice_TooHighInvalid(t *testing.T) {
	_, valid := validateChoice("6", 5)
	assert.False(t, valid)
}

func TestValidateChoice_NonNumericInvalid(t *testing.T) {
	_, valid := validateChoice("abc", 5)
	assert.False(t, valid)
}

func TestValidateChoice_EmptyStringInvalid(t *testing.T) {
	_, valid := validateChoice("", 5)
	assert.False(t, valid)
}

func TestValidateChoice_FloatInvalid(t *testing.T) {
	_, valid := validateChoice("2.5", 5)
	assert.False(t, valid)
}

func TestValidateChoice_VeryLargeNumberInvalid(t *testing.T) {
	_, valid := validateChoice("999999999999999999999", 5)
	assert.False(t, valid)
}

func TestValidateChoice_SingleOption(t *testing.T) {
	index, valid := validateChoice("1", 1)
	assert.True(t, valid)
	assert.Equal(t, 0, index)
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
