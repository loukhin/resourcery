package main

import (
	"archive/zip"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func generateResourcePack(resourcePath string) string {
	_ = os.RemoveAll("./cache")
	_ = os.Mkdir("./cache", 0755)
	file, err := os.Create("./cache/temporary.zip")
	if err != nil {
		log.Println(err)
	}
	w := zip.NewWriter(file)

	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasPrefix(info.Name(), ".") || strings.Contains(info.Name(), "~") {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func(file *os.File) {
			_ = file.Close()
		}(file)

		f, err := w.Create(path)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		return nil
	}

	err = filepath.Walk(resourcePath, walker)
	if err != nil {
		log.Fatal(err)
	}

	fileName := file.Name()
	_ = w.Close()
	_ = file.Close()

	file, err = os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	h := sha1.New()
	_, err = io.Copy(h, file)
	if err != nil {
		log.Fatal(err)
	}
	_ = file.Close()

	fileHash := hex.EncodeToString(h.Sum(nil))
	fileDir := path.Dir(file.Name())
	newFileName := path.Join(fileDir, fileHash+".zip")

	err = os.Rename(fileName, newFileName)
	if err != nil {
		log.Fatal(err)
	}

	return fileHash
}

func watch(watchPath string, onChange func()) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer func(watcher *fsnotify.Watcher) {
		err := watcher.Close()
		if err != nil {
			log.Fatal()
		}
	}(watcher)

	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}

		err = watcher.Add(path)
		if err != nil {
			log.Fatal(err)
		}

		return nil
	}

	err = filepath.Walk(watchPath, walker)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if strings.Contains(event.Name, "~") || strings.HasPrefix(event.Name, ".") || event.Has(fsnotify.Chmod) {
				continue
			}
			onChange()
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
