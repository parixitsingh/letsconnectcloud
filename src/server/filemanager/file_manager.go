package filemanager

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

type FileManager interface {
	ListFiles(*http.Request) (interface{}, error)
	AddFiles(*http.Request) (interface{}, error)
	UpdateFiles(*http.Request) (interface{}, error)
	RemoveFile(*http.Request) (interface{}, error)
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
		return nil, fmt.Errorf("filepath walk failed with %v", err)
	}
	return files, nil
}

func (fm *fileManager) AddFiles(r *http.Request) (interface{}, error) {
	var files []file
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&files); err != nil {
		return nil, fmt.Errorf("add files request body decoding failed with %v", err)
	}

	for _, file := range files {
		if err := checkFileStatus(file); err != nil {
			return nil, err
		}
	}

	for _, file := range files {
		if err := createFile(file); err != nil {
			return nil, err
		}
	}
	return files, nil
}

func (fm *fileManager) UpdateFiles(r *http.Request) (interface{}, error) {
	var files []file
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&files); err != nil {
		return nil, fmt.Errorf("add files request body decoding failed with %v", err)
	}

	for _, file := range files {
		if err := updateFile(file); err != nil {
			return nil, err
		}
	}
	return files, nil
}

func (fm *fileManager) RemoveFile(r *http.Request) (interface{}, error) {
	var fileDetail file
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&fileDetail); err != nil {
		return nil, fmt.Errorf("add files request body decoding failed with %v", err)
	}

	if err := removeFile(fileDetail); err != nil {
		return nil, err
	}
	return nil, nil
}

func getFilePath(fileName string) (string, error) {
	filePath, err := filepath.Abs("../files/" + fileName)
	if err != nil {
		return "", fmt.Errorf("error while creating file path %v", err)
	}
	return filePath, nil
}

func checkFileStatus(fileDetail file) error {
	filePath, err := getFilePath(fileDetail.Name)
	if err != nil {
		return err
	}

	_, err = os.Stat(filePath)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return nil
	}

	return fmt.Errorf("file already exists %v", fileDetail.Name)
}

func createFile(fileDetail file) error {
	filePath, err := getFilePath(fileDetail.Name)
	if err != nil {
		return nil
	}

	f, err := os.OpenFile(filePath, os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("create file failed with error %v", err)
	}

	defer f.Close()
	_, err = f.Write(fileDetail.Content)
	if err != nil {
		return fmt.Errorf("write file failed with error %v", err)
	}
	return nil
}

func updateFile(fileDetail file) error {
	filePath, err := getFilePath(fileDetail.Name)
	if err != nil {
		return nil
	}

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("open file failed with error %v", err)
	}

	defer f.Close()
	_, err = f.Write(fileDetail.Content)
	if err != nil {
		return fmt.Errorf("write file failed with error %v", err)
	}
	return nil
}

func removeFile(fileDetail file) error {
	filePath, err := getFilePath(fileDetail.Name)
	if err != nil {
		return nil
	}

	err = os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("delete file failed with error %v", err)
	}
	return nil
}
