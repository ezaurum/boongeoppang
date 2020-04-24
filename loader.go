package boongeoppang

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"html/template"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	baseOf             = "baseof"
	defaultDir         = "_default"
	partialsDir        = "_partials"
	themesDir          = "_themes"
	DefaultTemplateDir = "templates"
	DefaultTemplateExt = ".tmpl"
)

type LayoutHolder struct {
	Path   string
	Layout *template.Template
	Name   string
}

type TemplateContainer struct {
	M           map[string]*LayoutHolder
	Partials    map[string]string
	Themes      map[string]string
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
			"stringInSlice":   stringInSlice,
		},
	}
}

func LoadDefaultDebug() (*TemplateContainer, chan *TemplateContainer) {
	return LoadDebug(DefaultTemplateDir, nil)
}

func LoadDefault() *TemplateContainer {
	return Load(DefaultTemplateDir, nil)
}

func Load(rootDir string, funcMap template.FuncMap) *TemplateContainer {
	return Default().SetFuncMap(funcMap).Load(rootDir)
}

func LoadDebug(rootDir string, funcMap template.FuncMap) (*TemplateContainer, chan *TemplateContainer) {
	container := Default()
	container.debug = true
	container.SetFuncMap(funcMap)
	load := container.Load(rootDir)
	c := load.Watch(funcMap)
	return load, c
}

func (t *TemplateContainer) SetFuncMap(funcMap template.FuncMap) *TemplateContainer {
	if len(funcMap) > 0 {
		for k, v := range funcMap {
			t.FuncMap[k] = v
		}
	}
	return t
}

func (t *TemplateContainer) Watch(funcMap template.FuncMap) chan *TemplateContainer {
	c := make(chan *TemplateContainer)
	WatchDir(t.TemplateDir, func(watcher *fsnotify.Watcher) {
		for {
			select {
			case ev := <-watcher.Events:
				if DefaultTemplateExt == filepath.Ext(ev.Name) {
					tc := Load(t.TemplateDir, funcMap)
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
func walkInner(t *TemplateContainer) filepath.WalkFunc {

	return func(pathString string, info os.FileInfo, err error) error {
		if nil != err {
			log.Printf("err before %v, %v", pathString, err)
			return err
		}

		// 디렉토리는 패스
		if info.IsDir() && path.Base(pathString) != info.Name() {
			filepath.Walk(path.Join(pathString, info.Name()), walkInner(t))
		}

		filename := info.Name()
		ext := filepath.Ext(filename)

		// 템플릿이 아니면 패스
		if ext != DefaultTemplateExt {
			return nil
		}

		layoutName := strings.TrimSuffix(filename, ext)
		if layoutName == "" {
			return fmt.Errorf("file name is empty %v, %v", pathString, info)
		}

		contentName := AfterSecond(filepath.Dir(pathString))
		templateKey := contentName + "/" + layoutName
		switch contentName {
		case "":
			return fmt.Errorf("file name is empty %v, %v", pathString, info)
		case themesDir:
			t.Themes[templateKey] = pathString
		case partialsDir:
			t.Partials[layoutName] = pathString
		case defaultDir:
			t.Defaults[layoutName] = pathString
			templateKey = layoutName
			fallthrough
		default:
			t.M[templateKey] = &LayoutHolder{
				Name: layoutName,
				Path: pathString,
			}
			break
		}

		if t.debug {
			log.Printf("key:%v, name:%v, pathString:%v", templateKey, layoutName, pathString)
		}

		return err
	}
}

// Base returns the last element of path.
// Trailing path separators are removed before extracting the last element.
// If the path is empty, Base returns ".".
// If the path consists entirely of separators, Base returns a single separator.
func AfterSecond(pathString string) string {
	if pathString == "" {
		return "."
	}
	// Strip trailing slashes.
	for len(pathString) > 0 && os.IsPathSeparator(pathString[len(pathString)-1]) {
		pathString = pathString[0 : len(pathString)-1]
	}
	// Throw away volume name
	pathString = pathString[len(filepath.VolumeName(pathString)):]
	// Find the last element
	i := len(pathString)
	if i < 1 {
		return string(filepath.Separator)
	}
	l := 1
	for l < i && !os.IsPathSeparator(pathString[l]) {
		l++
	}
	if l < i {
		pathString = pathString[l+1:]
	}
	// If empty now, it had only slashes.
	if pathString == "" {
		return string(filepath.Separator)
	}
	return pathString
}

func (t *TemplateContainer) Load(rootDir string) *TemplateContainer {

	t.TemplateDir = rootDir

	filepath.Walk(rootDir, walkInner(t))

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
