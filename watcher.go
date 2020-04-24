package boongeoppang

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path/filepath"
)

func WatchDir(targetDir string, run func(watcher *fsnotify.Watcher)) *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()
	if nil != err {
		log.Fatal(err)
	}

	filepath.Walk(targetDir, func(path string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			//fmt.Printf("add watch path %v\n", path)
			err = watcher.Add(path)
			if err != nil {
				log.Fatal(err)
			}
		}
		return nil
	})

	go run(watcher)

	return watcher
}
