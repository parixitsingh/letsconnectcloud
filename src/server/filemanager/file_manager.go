package filemanager

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
)

type FileManager interface {
	ListFiles(*http.Request) (interface{}, error)
	AddFiles(*http.Request) (interface{}, error)
	//UpdateFiles(context.Context, *http.Request) (interface{}, error)
	//RemoveFiles(context.Context, *http.Request) (interface{}, error)
	//WordCounts(context.Context, *http.Request) (interface{}, error)
	//WordFrequency(context.Context, *http.Request) (interface{}, error)
}

type file struct {
	Name    string `json:"name"`
	Content []byte `json:"content"`
}

type fileManager struct{}

func NewFileManager() FileManager {
	return &fileManager{}
}

func (fm *fileManager) ListFiles(_ *http.Request) (interface{}, error) {
	return readDir("../files")
}

func readDir(fPath string) ([]string, error) {
	files := []string{}
	walkFunc := func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !fileInfo.IsDir() {
			files = append(files, path)
		}
		return nil
	}
	if err := filepath.Walk(fPath, walkFunc); err != nil {
		return nil, err
	}
	return files, nil
}

func (fm *fileManager) AddFiles(r *http.Request) (interface{}, error) {
	var files []file
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&files); err != nil {
		return nil, err
	}

	for _, file := range files {
		if err := writeFile(file); err != nil {
			return nil, err
		}
	}
	return files, nil
}

func writeFile(fileDetail file) error {
	f, err := os.OpenFile("./files/"+fileDetail.Name, os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(fileDetail.Content)
	return err
}
