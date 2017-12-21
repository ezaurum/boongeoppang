package boongeoppang

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"path"
	"log"
	"html/template"
)

const (
	baseOf = "baseof"
	defaultDir = "_default"
	partialsDir = "_partials"

)
var ( EmptyLayoutHolder = LayoutHolder{})

type LayoutHolder struct {
	Path   string
	Layout interface{}
	Name   string
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
			Name: name,
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
		templateKey := contentName + "/" + layoutName
		switch contentName {
		case "":
			return fmt.Errorf("file name is empty %v, %v", path, info)
		case partialsDir:
			partials[layoutName] = path
			break
		case defaultDir:
			defaults[layoutName] = path
			templateKey = layoutName
			fallthrough
		default:
			containers.M[templateKey] = &LayoutHolder{
				Name: layoutName,
				Path: path,
			}
			break
		}

		return err
	})

	var partialsFileNames []string
	for _, v := range partials {
		partialsFileNames = append(partialsFileNames, v)
	}

	// _default/baseof 먼제 체크
	base, isBase := defaults["baseof"]

	for key, value := range containers.M {
		layoutName := value.Name

		// 목록의 제일 처음이 기본 템플릿이 된다.
		var files []string

		// 1. baseof
		if isBase {
			files = append(files, base)
		}

		// 2. _default/layout
		if layoutName != key {
			ln, e := defaults[layoutName]
			if e && len(ln) > 0 {
				files = append(files, ln)
			}
		}

		// 3. domain/layout - path
		files = append(files, value.Path)

		// partials added after first object
		if len(files) > 1 {
			files = append(files[:1], append(partialsFileNames, files[1:]...)...)
		} else {
			files = append(files[:1], partialsFileNames...)
		}

		fmt.Println(files)

		must := template.Must(template.ParseFiles(files...))

		value.Layout = must
	}

	return containers
}
