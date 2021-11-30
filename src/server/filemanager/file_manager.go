package filemanager

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// FileManager is an interface which is exposing functionalities
type FileManager interface {
	ListFiles(*http.Request) (interface{}, error)
	AddFiles(*http.Request) (interface{}, error)
	UpdateFiles(*http.Request) (interface{}, error)
	RemoveFile(*http.Request) (interface{}, error)
	WordCounts(*http.Request) (interface{}, error)
	WordFrequency(*http.Request) (interface{}, error)
}

// file is representing file details
type file struct {
	Name    string `json:"name"`
	Content []byte `json:"content"`
}

// wordFrequencyRequest is representing the frequent words request
type wordFrequencyRequest struct {
	Limit int    `json:"limit"`
	Order string `json:"order"`
}

// wordFrequencyResponse is representing frequent words response
type wordFrequencyResponse struct {
	Words []string `json:"words"`
}

// wordCountResponse is representing word count response
type wordCountResponse struct {
	Count int `json:"count"`
}

// file manager is implementing FileManager
type fileManager struct{}

func NewFileManager() FileManager {
	return &fileManager{}
}

// ListFiles is returning list of all the files stored
func (fm *fileManager) ListFiles(_ *http.Request) (interface{}, error) {
	return readDir("../files")
}

// readDir is returning all the files inside in a file path
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

// AddFiles is creating files that are coming in the request
// If files already exist then returning error
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
	return nil, nil
}

// UpdateFiles is updating/creating files coming in the request
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
	return nil, nil
}

// RemoveFile is deleting the file coming in the request
// If file does not exist then returning error
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

// WordCounts is counting all the words from all files stored and
// returning the count (unique word count)
func (fm *fileManager) WordCounts(*http.Request) (interface{}, error) {
	wordCounts, err := getWordCounts()
	if err != nil {
		return nil, err
	}
	return &wordCountResponse{
		Count: len(wordCounts),
	}, nil
}

// WordFrequency is counting all the words and
// returning number of words mentioned in the request as limit
// in ascending or descending format as mentioned by order in the request
func (fm *fileManager) WordFrequency(r *http.Request) (interface{}, error) {
	wordFrequency := &wordFrequencyRequest{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(wordFrequency); err != nil {
		return nil, fmt.Errorf("word frequency request body decoding failed with %v", err)
	}

	wordCounts, err := getWordCounts()
	if err != nil {
		return nil, err
	}

	words := make([]string, 0, len(wordCounts))
	for word := range wordCounts {
		words = append(words, word)
	}

	sort.SliceStable(words, func(i, j int) bool {
		return wordCounts[words[i]] < wordCounts[words[j]]
	})

	wordFrequencyResponse := &wordFrequencyResponse{}

	if wordFrequency.Limit > len(words) {
		wordFrequency.Limit = len(words)
	}

	wordFrequencyResponse.Words = words[len(words)-wordFrequency.Limit:]
	switch wordFrequency.Order {
	case "asc":
		// no-op
	case "dsc":
		for i := 0; i < (len(wordFrequencyResponse.Words) / 2); i++ {
			wordFrequencyResponse.Words[i], wordFrequencyResponse.Words[len(wordFrequencyResponse.Words)-i-1] = wordFrequencyResponse.Words[len(wordFrequencyResponse.Words)-i-1], wordFrequencyResponse.Words[i]
		}
	default:
		return nil, fmt.Errorf("invalid order type %v %v", wordFrequency.Order, wordFrequency.Limit)
	}

	return wordFrequencyResponse, nil
}

// getWordCounts is counting words from all the files
// reading them in go routines (wordCounts)
func getWordCounts() (map[string]int, error) {
	files, err := readDir("../files")
	if err != nil {
		return nil, fmt.Errorf("error while reading files for word count %v", err)
	}

	c := make(chan map[string]int, len(files))
	errChan := make(chan error)
	for _, file := range files {
		go wordCounts(c, errChan, file)
	}

	wordCounts := make(map[string]int)
	for i := 0; i < len(files); i++ {
		select {
		case counts := <-c:
			for word, count := range counts {
				wordCounts[word] += count
			}
		case err := <-errChan:
			return nil, err
		}
	}
	return wordCounts, nil
}

// wordCounts is reading words from a file (fileName)
// if err occured writing error to errChan
// otherwise writing words count map to c chan
func wordCounts(c chan<- map[string]int, errChan chan<- error, fileName string) {
	wordCounts := make(map[string]int)
	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		errChan <- fmt.Errorf("error while opening file %v with error %v", fileName, err)
	}

	rdr := bufio.NewReader(file)
	for {
		line, err := rdr.ReadString('\n')
		if line == "" {
			break
		} else {
			words := strings.Fields(line)
			for _, word := range words {
				word = strings.ToLower(word)
				wordCounts[word]++
			}
		}
		if err != nil {
			if err != io.EOF {
				errChan <- fmt.Errorf("error while reading from file %v with error %v", fileName, err)
				break
			}
			break
		}
	}
	c <- wordCounts
}

// getFilePath is returning the Absolute file path of a file
// file path is used to store and retrieve the file
func getFilePath(fileName string) (string, error) {
	filePath, err := filepath.Abs("../files/" + fileName)
	if err != nil {
		return "", fmt.Errorf("error while creating file path %v", err)
	}
	return filePath, nil
}

// checkFileStatus is checking whether a file exist of not
// if file exists then it is an error
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

// createFile is creating a file
// if file already exists then it is an error
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

// updateFile is updating/creating a file
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

// removeFile is removing a file
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
