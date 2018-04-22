package boongeoppang

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"github.com/fsnotify/fsnotify"
)

const (
	baseOf             = "baseof"
	defaultDir         = "_default"
	partialsDir        = "_partials"
	DefaultTemplateDir = "templates"
	DefaultTemplateExt = ".tmpl"
)

var (
	EmptyLayoutHolder = LayoutHolder{}
)

type LayoutHolder struct {
	Path   string
	Layout *template.Template
	Name   string
}

type TemplateContainer struct {
	M           map[string]*LayoutHolder
	Partials    map[string]string
	Defaults    map[string]string
	FuncMap     template.FuncMap
	debug       bool
	TemplateDir string
}

func (t TemplateContainer) Set(name string, layout *template.Template) {
	get, _ := t.Get(name)
	get.Layout = layout
}

func (t TemplateContainer) Get(name string) (*LayoutHolder, bool) {
	if r, b := t.M[name]; b {
		return r, true
	}

	baseName := path.Base(name)

	if mm, b := t.Defaults[baseName]; b && baseName != baseOf {
		t.M[name] = &LayoutHolder{
			Name:   name,
			Path:   mm,
			Layout: t.M[baseName].Layout,
		}
		return t.M[name], true
	}

	return nil, false
}

func Default() *TemplateContainer {
	partials := make(map[string]string)
	defaults := make(map[string]string)
	return &TemplateContainer{
		Partials: partials,
		Defaults: defaults,
		M:        make(map[string]*LayoutHolder),
		FuncMap: template.FuncMap{
			"asDate":          asDate,
			"asDate12HMinute": asDate12HMinute,
			"asDate24HMinute": asDate24HMinute,
			"asTime12H":       asTime12H,
			"asTime24H":       asTime24H,
		},
	}
}

func LoadDefaultDebug() (*TemplateContainer, chan *TemplateContainer) {
	return LoadDebug(DefaultTemplateDir)
}

func LoadDefault() *TemplateContainer {
	return Load(DefaultTemplateDir)
}

func Load(rootDir string) *TemplateContainer {
	return Default().Load(rootDir)
}

func LoadDebug(rootDir string) (*TemplateContainer, chan *TemplateContainer) {
	container := Default()
	container.debug = true
	load := container.Load(rootDir)
	c := load.Watch()
	return load, c
}

func (t *TemplateContainer) Watch() chan *TemplateContainer {
	c := make(chan *TemplateContainer)
	WatchDir(t.TemplateDir, func(watcher *fsnotify.Watcher) {
		for {
			select {
			case ev := <-watcher.Events:
				if ev.Op&fsnotify.Create == fsnotify.Create &&
					DefaultTemplateExt == filepath.Ext(ev.Name) {
					tc := Load(t.TemplateDir)
					tc.debug = t.debug
					c <- tc
				}
			case err := <-watcher.Errors:
				log.Fatal("error:", err)
			}
		}
	})
	return c
}

func (t *TemplateContainer) Load(rootDir string) *TemplateContainer {

	t.TemplateDir = rootDir

	filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if nil != err {
			log.Printf("err before %v, %v", path, err)
		}

		// 디렉토리는 패스
		if info.IsDir() {
			return nil
		}

		filename := info.Name()
		ext := filepath.Ext(filename)

		// 템플릿이 아니면 패스
		if ext != DefaultTemplateExt {
			return nil
		}

		layoutName := strings.TrimSuffix(filename, ext)
		if layoutName == "" {
			return fmt.Errorf("file name is empty %v, %v", path, info)
		}

		contentName := filepath.Base(filepath.Dir(path))
		templateKey := contentName + "/" + layoutName
		switch contentName {
		case "":
			return fmt.Errorf("file name is empty %v, %v", path, info)
		case partialsDir:
			t.Partials[layoutName] = path
			break
		case defaultDir:
			t.Defaults[layoutName] = path
			templateKey = layoutName

			fallthrough
		default:
			t.M[templateKey] = &LayoutHolder{
				Name: layoutName,
				Path: path,
			}
			break
		}

		if t.debug {
			log.Printf("key:%v, name:%v, path:%v", templateKey, layoutName, path)
		}

		return err
	})

	t.initiateTemplates()

	return t
}

// initiate html/template
func (t *TemplateContainer) initiateTemplates() {

	var partialsFileNames []string
	for _, v := range t.Partials {
		partialsFileNames = append(partialsFileNames, v)
	}
	// _default/baseof 먼제 체크
	base, isBase := t.Defaults["baseof"]
	for key, value := range t.M {
		layoutName := value.Name

		// 목록의 제일 처음이 기본 템플릿이 된다.
		var files []string

		// 1. baseof
		if isBase {
			files = append(files, base)
		}

		// 2. _default/layout
		if layoutName != key {
			ln, e := t.Defaults[layoutName]
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

		// baseof.tmpl이 실행되어야 한다.
		// parse를 하면 각 파일 이름별로 하나씩 내부에 템플릿이 만들어진다
		value.Layout =
			template.Must(template.
				New(path.Base(files[0])).
				Funcs(t.FuncMap).ParseFiles(files...))
	}
}
