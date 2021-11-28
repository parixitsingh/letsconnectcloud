package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

/*
command structure

store [command_type] [values]

for eg:-
store add filename1 filename2

command_type is add
values are filenames {filename1, filename2}

*/

const (
	LS        string = "ls"
	ADD       string = "add"
	UPDATE    string = "update"
	RM        string = "rm"
	WC        string = "wc"
	FREQWORDS string = "freq-words"
)

const (
	baseURL = "http://localhost:8080/"
)

type file struct {
	Name    string `json:"name"`
	Content []byte `json:"content"`
}

type store struct {
	client  *http.Client
	command string
	options []string
}

func newStore() *store {
	// reading the args
	// skipping the 0 index as it gives program
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("store version is 1.0")
		os.Exit(0)
	}
	return &store{
		client:  &http.Client{},
		command: args[0],
		options: args[1:],
	}
}

func main() {

	// creating new store
	store := newStore()

	switch store.command {
	case LS:
		store.listFiles()
	case ADD:
		store.addFiles()
	case "update":
		store.updateFiles()
	case "rm":
		store.removeFile()
	// case "wc":
	// 	response, err := wordCounts(httpClient)
	// 	if err != nil {
	// 		fmt.Printf("error occured while counting the words %v", err)
	// 		os.Exit(1)
	// 	}
	// 	fmt.Println(response)
	// case "freq-words":
	// 	response, err := frequencyWords(httpClient)
	// 	if err != nil {
	// 		fmt.Printf("error occured while counting the words frequencies %v", err)
	// 		os.Exit(1)
	// 	}
	// 	fmt.Println(response)
	default:
		fmt.Println(fmt.Errorf("command \"%s\" is not valid", store.command))
	}
}

func (st *store) listFiles() {
	bodyBytes, err := st.createAndExecuteHTTPRequest(http.MethodGet, "listfiles", nil)
	if err != nil {
		fmt.Printf("error occured while fetching the files : %v", err)
	}

	if len(bodyBytes) == 0 {
		fmt.Println("no files exist on server")
	}

	var files []string
	if err := json.Unmarshal(bodyBytes, &files); err != nil {
		fmt.Println("error while reading the response from server")
	}

	for _, file := range files {
		fmt.Println(" ", file)
	}
}

func (st *store) addFiles() {
	if len(st.options) == 0 {
		fmt.Println("no files are specified")
	}

	files := []file{}

	for _, v := range st.options {
		name, content, err := st.getFileContent(v)
		if err != nil {
			fmt.Printf("error occured while reading the files %v", err)
			os.Exit(1)
		}
		files = append(files, file{
			Name:    name,
			Content: content,
		})
	}

	_, err := st.createAndExecuteHTTPRequest(http.MethodPost, "addfiles", files)
	if err != nil {
		fmt.Printf("error occured while adding the files :%v", err)
		os.Exit(1)
	}

	fmt.Println("files added successfully")
}

func (st *store) updateFiles() {
	if len(st.options) == 0 {
		fmt.Println("no files are specified")
	}

	files := []file{}

	for _, v := range st.options {
		name, content, err := st.getFileContent(v)
		if err != nil {
			fmt.Printf("error occured while reading the files %v", err)
			os.Exit(1)
		}
		files = append(files, file{
			Name:    name,
			Content: content,
		})
	}

	_, err := st.createAndExecuteHTTPRequest(http.MethodPut, "updatefiles", files)
	if err != nil {
		fmt.Printf("error occured while updating the files :%v", err)
		os.Exit(1)
	}

	fmt.Println("files updated successfully")
}

func (st *store) removeFile() {
	if len(st.options) == 0 {
		fmt.Println("no files are specified")
	}

	if len(st.options) > 1 {
		fmt.Println("more than one files are specified")
	}

	fileStructure := strings.Split(st.options[0], "/")
	_, err := st.createAndExecuteHTTPRequest(http.MethodDelete, "removefile", &file{
		Name: fileStructure[len(fileStructure)-1],
	})
	if err != nil {
		fmt.Printf("error occured while deleting the file :%v", err)
		os.Exit(1)
	}

	fmt.Println("file deleted successfully")
}

func (st *store) getFileContent(fPath string) (string, []byte, error) {
	fileStructure := strings.Split(fPath, "/")
	if fileStructure[0] == fPath {
		fileStructure = strings.Split(fPath, `\`)
	}
	fileContent, err := ioutil.ReadFile(fPath)
	if err != nil {
		return "", nil, err
	}

	return fileStructure[len(fileStructure)-1], fileContent, nil
}

func wordCounts(*http.Client) (interface{}, error) {
	return "word counts called", nil
}

func frequencyWords(*http.Client) (interface{}, error) {
	return "frequency words called", nil
}

func (st *store) createAndExecuteHTTPRequest(method, url string, reqBody interface{}) ([]byte, error) {
	reqUrl := baseURL + url
	requestBody := []byte{}
	var err error

	if reqBody != nil {
		requestBody, err = json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("error while marshaling the request body with %w", err)
		}
	}

	req, err := http.NewRequest(method, reqUrl, bytes.NewReader(requestBody))
	if err != nil {
		return nil, err
	}

	res, err := st.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		return bodyBytes, nil
	}

	return nil, fmt.Errorf("called failed with statuscode: %v", res.StatusCode)
}
