package whitewalker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"path"
	"log"
)

const (
	baseOf = "baseof"
	defaultDir = "_default"
	partialsDir = "partials"

)
var ( EmptyLayoutHolder = LayoutHolder{})

type LayoutHolder struct {
	Path   string
	Layout *interface{}
	Key    string
}

type TemplateContainer struct {
	M        map[string]*LayoutHolder
	Partials map[string]string
	Defaults map[string]string
}

func ( t TemplateContainer)  Set(name string, layout interface{})  {
	get, _ := t.Get(name)
	get.Layout = &layout
}

func ( t TemplateContainer)  Get(name string) ( *LayoutHolder, bool ) {
	if r, b :=t.M[name] ; b {
		return r, true
	}

	baseName := path.Base(name)

	if mm, b := t.Defaults[baseName] ; b && baseName != baseOf {
		t.M[name] = &LayoutHolder{
			Key:  name,
			Path: mm,
		}
		return t.M[name],true
	}

	return nil, false
}

func Load(rootDir string) *TemplateContainer {
	partials := make(map[string]string)
	defaults := make(map[string]string)

	containers := &TemplateContainer{
		Partials: partials,
		Defaults: defaults,
		M:        make(map[string]*LayoutHolder),
	}

	filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if nil != err {
			log.Printf("err before %v, %v",path, err)
		}

		// 디렉토리는 패스
		if info.IsDir() {
			return nil
		}

		filename := info.Name()
		layoutName := strings.TrimSuffix(filename, filepath.Ext(filename))
		if layoutName == "" {
			return fmt.Errorf("file name is empty %v, %v", path, info)
		}

		contentName := filepath.Base(filepath.Dir(path))
		switch contentName {
		case "":
			return fmt.Errorf("file name is empty %v, %v", path, info)
		case partialsDir:
			partials[layoutName] = path
			break
		case defaultDir:
			defaults[layoutName] = path
			if baseOf != layoutName {
				containers.M[layoutName] = &LayoutHolder{
					//Partials: partials,
					Key:  layoutName,
					Path: path,
				}
			}
			break
		default:
			containers.M[contentName+"/"+layoutName] = &LayoutHolder{
				Key:  layoutName,
				Path: path,
			}
			break
		}

		return err
	})

	return containers
}
