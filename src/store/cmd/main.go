package main

import (
	"fmt"

	"store/storemanager"
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

func main() {
	// creating new store
	storeManager := storemanager.NewStore(baseURL)

	switch storeManager.Command() {
	case LS:
		storeManager.ListFiles()
	case ADD:
		storeManager.AddFiles()
	case "update":
		storeManager.UpdateFiles()
	case "rm":
		storeManager.RemoveFile()
	case "wc":
		storeManager.WordCounts()
	case "freq-words":
		storeManager.WordFrequency()
	default:
		fmt.Println(fmt.Errorf("command \"%s\" is not valid", storeManager.Command))
	}
}
