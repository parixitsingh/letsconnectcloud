package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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
	store := newStore()

	switch store.command {
	case LS:
		store.listFiles()
	case ADD:
		response, err := store.addFiles()
		if err != nil {
			fmt.Printf("error occured while adding the files %v", err)
			os.Exit(1)
		}
		fmt.Println(response)
	// case "update":
	// 	response, err := updateFiles(httpClient, args[1:])
	// 	if err != nil {
	// 		fmt.Printf("error occured while updating the files %v", err)
	// 		os.Exit(1)
	// 	}
	// 	fmt.Println(response)
	// case "rm":
	// 	response, err := removeFiles(httpClient, args[1:])
	// 	if err != nil {
	// 		fmt.Printf("error occured while removing the files %v", err)
	// 		os.Exit(1)
	// 	}
	// 	fmt.Println(response)
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
	req, err := http.NewRequest(http.MethodGet, baseURL+"listfiles", nil)
	if err != nil {
		fmt.Printf("not able to fetch files due following error : %v", err)
	}

	bodyBytes, err := st.executeHTTPRequest(req)
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

func (st *store) addFiles() (interface{}, error) {
	if len(st.options) == 0 {
		return "no files are specified", nil
	}
	return "add files called", nil
}

// func updateFiles(_ *http.Client, files []string) (interface{}, error) {
// 	if len(files) == 0 {
// 		return "no files are specified", nil
// 	}
// 	return "update files called", nil
// }

// func removeFiles(*http.Client,files []string) (interface{}, error) {
// 	if len(files) == 0 {
// 		return "no files are specified", nil
// 	}
// 	return "remove files called", nil
// }

// func wordCounts(*http.Client) (interface{}, error) {
// 	return "word counts called", nil
// }

// func frequencyWords(*http.Client) (interface{}, error) {
// 	return "frequency words called", nil
// }

/*
Q. In update and remove should it be single or multiple files?
Q.
*/

func (st *store) executeHTTPRequest(req *http.Request) ([]byte, error) {
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
