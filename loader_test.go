package boongeoppang

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"strings"
	"testing"
	"html/template"
	"os"
	"github.com/gin-gonic/gin"
	"bytes"
	"time"
)

const testTemplateDir = "tests"

func TestBaseLayoutLoad(t *testing.T) {

	container := Load(testTemplateDir)

	notExist := []string{"test", "head", "foot"}
	for _, el := range notExist {
		layout, exist := container.Get(el)
		assert.False(t, exist)
		assert.Nil(t, layout)
	}

	defaultsExpected := []string{"index", "single", "list",  "baseof"}
	for _, el := range defaultsExpected {
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
		assert.True(t, strings.Index(path, testTemplateDir) > -1)
	}

	partialsExpected := []string{"head", "body"}
	for _, el := range partialsExpected {
		path := container.Partials[el]
		assert.NotEmpty(t, path)
		assert.True(t, strings.Index(path, el) > -1)
		assert.True(t, strings.Index(path, ".tmpl") > -1)
		assert.True(t, strings.Index(path, testTemplateDir) > -1)
	}
}

func TestContentSpecifiedLayoutLoad(t *testing.T) {

	container := Load(testTemplateDir)

	defaultsExpected := []string{"product/index", "product/single", "product/list", "product/form"}
	for _, el := range defaultsExpected {

		layout, exist := container.Get(el)
		assert.True(t, exist)
		assert.NotNil(t, layout)

		path := layout.Path

		assert.NotEmpty(t, path)
		assert.True(t, strings.Index(path, filepath.Base(el)) > -1)
		assert.True(t, strings.Index(path, ".tmpl") > -1)
		assert.True(t, strings.Index(path, testTemplateDir) > -1)

		layout.Layout.Execute(os.Stdout, nil)
	}
}

func TestLayoutSetGet(t *testing.T) {

	container := Load(testTemplateDir)
	expected := template.New("Test")
	container.Set("index", expected)

	layout, b := container.Get("index")

	assert.True(t, b)
	assert.Equal(t, expected, layout.Layout)
}

func TestExecute(t *testing.T) {

	container := Load(testTemplateDir)
	expected := "common/login"
	holder, b := container.Get(expected)
	if !b {
		t.Fail()
	}

	buf := bytes.NewBufferString("")
	template := holder.Layout
	err := template.Execute(buf, gin.H{})

	assert.Nil(t, err, err)
	assert.True(t, len(buf.String()) > 0)
}

func TestExecuteFuncMap(t *testing.T) {

	container := Load(testTemplateDir)
	expected := "index"
	holder, b := container.Get(expected)
	if !b {
		t.Fail()
	}
	targetTime, _ := time.Parse("2006-01-02", "2017-01-31")
	buf := bytes.NewBufferString("")
	template := holder.Layout
	err := template.Execute(buf, gin.H{
		"TestDate": targetTime,
	})

	assert.Nil(t, err, err)
	s := buf.String()
	assert.True(t, len(s) > 0)
	assert.True(t, strings.Contains(s,targetTime.Format("2006-01-02")))
}
