package dtrain

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"path"
)

const (
	baseOf = "baseof"
	defaultDir = "_default"
	partialsDir = "partials"

)

type FilenameHolder struct {
	Filename string
	Key      string
}

type LayoutHolder struct {
	Content  string
	Key      string
}

type mapHolder map[string]LayoutHolder
type TemplateContainer struct {
	M        mapHolder
	Partials map[string]FilenameHolder
	Defaults map[string]FilenameHolder
}

func ( t TemplateContainer)  Get(name string) LayoutHolder {
	if r, b :=t.M[name] ; b {
		return r
	}

	baseName := path.Base(name)

	if mm, b := t.Defaults[baseName] ; b && baseName != baseOf {
		t.M[name] = LayoutHolder{
			Content:mm.Filename,
			Key:name,
		}
		return t.M[name]
	}

	return LayoutHolder{}
}

func Load(rootDir string) *TemplateContainer {
	partials := make(map[string]FilenameHolder)
	defaults := make(map[string]FilenameHolder)

	containers := &TemplateContainer{
		Partials: partials,
		Defaults: defaults,
		M:        make(map[string]LayoutHolder),
	}

	filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		//TODO 중간에 멈춰야 하나? 안 멈추는 게 좋을듯
		if nil != err {
			return err
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
		case defaultDir:
			defaults[layoutName] = FilenameHolder{Filename: filename, Key: layoutName}
			if baseOf != layoutName {
				containers.M[layoutName] = LayoutHolder{
					//Partials: partials,
					Key:     layoutName,
					Content: contentName,
				}
			}
			break
		case partialsDir:
			partials[layoutName] = FilenameHolder{Filename: filename, Key: layoutName}
			break
		default:
			containers.M[contentName+"/"+layoutName] = LayoutHolder{
				Key:     layoutName,
				Content: contentName,
			}
			break
		}

		return err
	})

	return containers
}
