package render

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/ezaurum/boongeoppang"
)

//check implementation
var _ render.HTMLRender = Render{}

func Default() (Render) {
	return New(boongeoppang.DefaultTemplateDir)
}

func NewDebug(templateDir string, engine *gin.Engine) {
	b, c := boongeoppang.LoadDebug(templateDir)
	engine.HTMLRender = Render{
		templateContainer: b,
	}
	go func() {
		tc := <-c
		engine.HTMLRender = Render{
			templateContainer: tc,
		}
	}()
}

// New instance
func New(templateDir string) Render {
	return Render{
		templateContainer: boongeoppang.Load(templateDir),
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
