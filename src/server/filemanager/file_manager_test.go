package filemanager

import (
	"bytes"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
)

func getReq(method, url string, body interface{}) *http.Request {
	requestBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, url, bytes.NewReader(requestBody))
	return req
}

func Test_fileManager_AddFiles(t *testing.T) {
	type file struct {
		Name    string `json:"name"`
		Content []byte `json:"content"`
	}

	type args struct {
		r *http.Request
	}

	exampleFile := file{
		Name:    "test.txt",
		Content: []byte("testing"),
	}

	tests := []struct {
		name    string
		fm      *fileManager
		args    args
		want    interface{}
		wantErr bool
	}{
		{name: "success",
			fm: &fileManager{},
			args: args{
				r: getReq(http.MethodGet, "fakeURL", []*file{&exampleFile}),
			},
		},
		{name: "negative",
			fm: &fileManager{},
			args: args{
				r: getReq(http.MethodGet, "fakeURL", []*file{&exampleFile}),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fm.AddFiles(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("fileManager.AddFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fileManager.AddFiles() = %v, want %v", got, tt.want)
			}
		})
	}
	(&fileManager{}).RemoveFile(getReq(http.MethodDelete, "fakeURL", exampleFile))
}

func Test_fileManager_UpdateFiles(t *testing.T) {
	type args struct {
		r *http.Request
	}
	exampleFile := file{
		Name:    "test.txt",
		Content: []byte("testing"),
	}
	tests := []struct {
		name    string
		fm      *fileManager
		args    args
		want    interface{}
		wantErr bool
	}{
		{name: "success",
			fm: &fileManager{},
			args: args{
				r: getReq(http.MethodPut, "fakeURL", []*file{&exampleFile}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fm.UpdateFiles(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("fileManager.UpdateFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fileManager.UpdateFiles() = %v, want %v", got, tt.want)
			}
		})
	}
	(&fileManager{}).RemoveFile(getReq(http.MethodDelete, "fakeURL", exampleFile))
}

func Test_fileManager_RemoveFile(t *testing.T) {
	type args struct {
		r *http.Request
	}
	exampleFile := file{
		Name:    "test.txt",
		Content: []byte("testing"),
	}
	(&fileManager{}).UpdateFiles(getReq(http.MethodPut, "fakeURL", []*file{&exampleFile}))
	tests := []struct {
		name    string
		fm      *fileManager
		args    args
		want    interface{}
		wantErr bool
	}{
		{name: "success",
			fm: &fileManager{},
			args: args{
				r: getReq(http.MethodDelete, "fakeURL", exampleFile),
			},
		},
		{name: "negative",
			fm: &fileManager{},
			args: args{
				r: getReq(http.MethodDelete, "fakeURL", exampleFile),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fm.RemoveFile(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("fileManager.RemoveFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fileManager.RemoveFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
