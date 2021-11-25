package main

import (
	"fmt"
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

func main() {
	// reading the args
	// skipping the 0 index as it gives program
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("store version is 1.0")
		os.Exit(0)
	}

	switch args[0] {
	case "ls":
		listFiles()
	case "add":
		response, err := addFiles(args[1:])
		if err != nil {
			fmt.Printf("error occured while adding the files %v", err)
			os.Exit(1)
		}
		fmt.Println(response)
	case "update":
		response, err := updateFiles(args[1:])
		if err != nil {
			fmt.Printf("error occured while updating the files %v", err)
			os.Exit(1)
		}
		fmt.Println(response)
	case "rm":
		response, err := removeFiles(args[1:])
		if err != nil {
			fmt.Printf("error occured while removing the files %v", err)
			os.Exit(1)
		}
		fmt.Println(response)
	case "wc":
		response, err := wordCounts()
		if err != nil {
			fmt.Printf("error occured while counting the words %v", err)
			os.Exit(1)
		}
		fmt.Println(response)
	case "freq-words":
		response, err := frequencyWords()
		if err != nil {
			fmt.Printf("error occured while counting the words frequencies %v", err)
			os.Exit(1)
		}
		fmt.Println(response)
	default:
		fmt.Println(fmt.Errorf("command \"%s\" is not valid", args[0]))
	}
}

func listFiles() error {
	fmt.Println("list files called")
	return nil
}

func addFiles(files []string) (interface{}, error) {
	if len(files) == 0 {
		return "no files are specified", nil
	}
	return "add files called", nil
}

func updateFiles(files []string) (interface{}, error) {
	if len(files) == 0 {
		return "no files are specified", nil
	}
	return "update files called", nil
}

func removeFiles(files []string) (interface{}, error) {
	if len(files) == 0 {
		return "no files are specified", nil
	}
	return "remove files called", nil
}

func wordCounts() (interface{}, error) {
	return "word counts called", nil
}

func frequencyWords() (interface{}, error) {
	return "frequency words called", nil
}

/*
Q. In update and remove should it be single or multiple files?
Q.
*/
