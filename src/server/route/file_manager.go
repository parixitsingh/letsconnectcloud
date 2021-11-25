package route

import (
	"context"
	"io/ioutil"
	"net/http"
)

type FileManager interface {
	ListFiles(context.Context, *http.Request) (interface{}, error)
	//AddFiles(context.Context, *http.Request) (interface{}, error)
	//UpdateFiles(context.Context, *http.Request) (interface{}, error)
	//RemoveFiles(context.Context, *http.Request) (interface{}, error)
	//WordCounts(context.Context, *http.Request) (interface{}, error)
	//WordFrequency(context.Context, *http.Request) (interface{}, error)
}

type fileManager struct{}

func NewFileManager() FileManager {
	return &fileManager{}
}

func (fm *fileManager) ListFiles(ctx context.Context, r *http.Request) (interface{}, error) {
	var files []string
	filesDetails, err := ioutil.ReadDir("../files")
	if err != nil {
		return files, err
	}

	for _, file := range filesDetails {
		files = append(files, file.Name())
	}
	return files, nil
}
