package render

import (
	"github.com/ezaurum/boongeoppang"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"html/template"
)

//check implementation
var _ render.HTMLRender = Render{}

func Default() Render {
	return New(boongeoppang.DefaultTemplateDir, nil)
}

func NewDebug(templateDir string, funcMap template.FuncMap, engine *gin.Engine) {
	b, c := boongeoppang.LoadDebug(templateDir, funcMap)
	engine.HTMLRender = Render{
		templateContainer: b,
	}
	go func() {
		for {
			tc := <-c
			engine.HTMLRender = Render{
				templateContainer: tc,
			}
		}
	}()
}

// New instance
func New(templateDir string, funcMap template.FuncMap) Render {
	return Render{
		templateContainer: boongeoppang.Load(templateDir, funcMap),
	}
}

type Render struct {
	templateContainer *boongeoppang.TemplateContainer
}

// Instance find by name
func (r Render) Instance(name string, data interface{}) render.Render {
	layout, isExist := r.templateContainer.Get(name)
	if !isExist {
		panic("not exist template " + name)
	}
	return render.HTML{
		Template: layout.Layout,
		Data:     data,
	}
}
