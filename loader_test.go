package boongeoppang

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"strings"
	"path/filepath"
)

func TestBaseLayoutLoad(t *testing.T) {

	rootDir := "tests/full"
	container := Load(rootDir)

	notExist := []string{"test","head","foot"}
	for _,el := range notExist {
		layout, exist := container.Get(el)
		assert.False(t, exist)
		assert.Nil(t, layout)
	}

	defaultsExpected := []string{"index","single","list", "form", "baseof"}
	for _,el := range defaultsExpected {
		path := container.Defaults[el]

		if el != "baseof" {
			layout, exist := container.Get(el)
			assert.True(t, exist)
			assert.NotNil(t, layout)
			assert.Equal(t, layout.Path, path)
		}

		assert.NotEmpty(t, path)
		assert.True(t, strings.Index(path, el) > -1)
		assert.True(t, strings.Index(path, ".tmpl") > -1)
		assert.True(t, strings.Index(path, "tests") > -1)
	}

	partialsExpected := []string{"head","body"}
	for _,el := range partialsExpected {
		path := container.Partials[el]
		assert.NotEmpty(t, path)
		assert.True(t, strings.Index(path, el) > -1)
		assert.True(t, strings.Index(path, ".tmpl") > -1)
		assert.True(t, strings.Index(path, "tests") > -1)
	}
}

func TestContentSpecifiedLayoutLoad(t *testing.T) {

	container := Load("tests/full")

	defaultsExpected := []string{"product/index","product/single","product/list", "product/form" }
	for _,el := range defaultsExpected {

		layout, exist := container.Get(el)
		assert.True(t, exist)
		assert.NotNil(t, layout)

		path := layout.Path

		assert.NotEmpty(t, path)
		assert.True(t, strings.Index(path, filepath.Base(el)) > -1)
		assert.True(t, strings.Index(path, ".tmpl") > -1)
		assert.True(t, strings.Index(path, "tests") > -1)
	}


}

func TestLayoutSetGet(t *testing.T) {

	container := Load("tests/full")
	expected := "IndexLayout"
	container.Set("index", expected)

	layout, b := container.Get("index")

	assert.True(t, b)
	assert.Equal(t, expected, layout.Layout)
}
