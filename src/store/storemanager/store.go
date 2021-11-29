package storemanager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// StoreManager is an interface which is exposing functionalities
type StoreManager interface {
	Command() string
	ListFiles()
	AddFiles()
	UpdateFiles()
	RemoveFile()
	WordCounts()
	WordFrequency()
}

// wordCountResponse is response when counting words
type wordCountResponse struct {
	Count int `json:"count"`
}

// wordFrequencyResponse is response when frequent words counted
type wordFrequencyResponse struct {
	Words []string `json:"words"`
}

// wordFrequencyRequest is request when frequent words counted
type wordFrequencyRequest struct {
	Limit int    `json:"limit"`
	Order string `json:"order"`
}

// file is representing file details name and content
type file struct {
	Name    string `json:"name"`
	Content []byte `json:"content"`
}

// store is implementing all the function that executes on different commands
type store struct {
	client  *http.Client
	command string
	options []string
	baseURL string
}

// NewStore is constructor to create store instance
func NewStore(baseURL string) StoreManager {
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
		baseURL: baseURL,
	}
}

func (st *store) Command() string {
	return st.command
}

func (st *store) ListFiles() {
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

func (st *store) AddFiles() {
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

func (st *store) UpdateFiles() {
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

func (st *store) RemoveFile() {
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

func (st *store) WordCounts() {
	bodyBytes, err := st.createAndExecuteHTTPRequest(http.MethodGet, "wordscount", nil)
	if err != nil {
		fmt.Printf("error occured while getting the words count : %v", err)
	}

	if len(bodyBytes) == 0 {
		fmt.Println("no files exist on server")
	}

	wordCountResponse := &wordCountResponse{}

	if err := json.Unmarshal(bodyBytes, &wordCountResponse); err != nil {
		fmt.Println("error while reading the response from server")
	}

	fmt.Printf("Total words are %v", wordCountResponse.Count)
}

func (st *store) WordFrequency() {
	wordFrequency := &wordFrequencyRequest{
		Limit: 10,
		Order: "asc",
	}

	for i, v := range st.options {
		if strings.HasPrefix(v, "--limit") || strings.HasPrefix(v, "-n") && i < len(st.options)-1 {
			limit, err := strconv.Atoi(st.options[i+1])
			if err != nil {
				fmt.Printf("invalid limit provided %v", st.options[i+1])
				os.Exit(1)
			}
			wordFrequency.Limit = limit
		}

		if strings.HasPrefix(v, "--order=") && i < len(st.options) {
			order := strings.TrimPrefix(st.options[i], "--order=")
			if order != "asc" && order != "dsc" {
				fmt.Printf("invalid order provided %v", st.options[i])
				os.Exit(1)
			}
			wordFrequency.Order = order
		}
	}

	bodyBytes, err := st.createAndExecuteHTTPRequest(http.MethodGet, "wordsfrequency", wordFrequency)
	if err != nil {
		fmt.Printf("error occured while getting the words count : %v", err)
	}

	if len(bodyBytes) == 0 {
		fmt.Println("no files exist on server")
	}

	wordFrequencyResponse := &wordFrequencyResponse{}

	if err := json.Unmarshal(bodyBytes, &wordFrequencyResponse); err != nil {
		fmt.Println("error while reading the response from server")
	}

	fmt.Printf("words are %v", wordFrequencyResponse.Words)
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

func (st *store) createAndExecuteHTTPRequest(method, url string, reqBody interface{}) ([]byte, error) {
	reqUrl := st.baseURL + url
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
