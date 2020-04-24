package echo

import (
	"github.com/ezaurum/boongeoppang"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"path"
)

var _ echo.Renderer = &Template{}

type Template struct {
	templateContainer *boongeoppang.TemplateContainer
}

func NewDebug(templateDir string, funcMap template.FuncMap) *Template {
	b, c := boongeoppang.LoadDebug(templateDir, funcMap)
	t := &Template{
		templateContainer: b,
	}
	go func() {
		for {
			tc := <-c
			t.templateContainer = tc
		}
	}()
	return t
}

// New instance
func New(templateDir string, funcMap template.FuncMap) *Template {
	return &Template{
		templateContainer: boongeoppang.Load(templateDir, funcMap),
	}
}

func (t *Template) Render(w io.Writer, name string,
	data interface{}, c echo.Context) error {
	layoutName := path.Base(name)
	layout, isExist := t.templateContainer.Get(name)
	if !isExist {
		layout, isExist = t.templateContainer.Get(path.Join("common", layoutName))
	}
	if !isExist {
		layout, isExist = t.templateContainer.Get(path.Join("_default", layoutName))
	}
	return layout.Layout.ExecuteTemplate(w, "baseof.tmpl", data)
}
