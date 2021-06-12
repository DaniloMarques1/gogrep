package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type Match struct {
	fileName    string
	lineNumber  int
	lineContent string
}

func main() {
	args, err := readArguments()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	searchString := args[1]
	fileName := args[2]
	if len(args) == 3 {
		execute(fileName, searchString)
	} else if len(args) == 4 {
		if args[3] != "-R" {
			fmt.Println("gogrep: Wrong flag")
			os.Exit(1)
		}

		err = handleRecursiveSearch(searchString, fileName)
	}
}

func execute(fileName, searchString string) {
	lines, err := readFileLines(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	matches := searchForMatches(searchString, lines)
	printMatches(matches)
}

// will read and return all arguments given
func readArguments() ([]string, error) {
	args := os.Args
	if len(args) < 3 {
		return nil, errors.New("Not enough arguments")
	}

	return args, nil
}

// will read the file and return a slice of lines
func readFileLines(fileName string) ([]string, error) {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	fileContent := string(bytes)
	lines := strings.Split(fileContent, "\n")

	return lines, nil
}

// receive a slice of lines and search for the searchString
func searchForMatches(searchString string, lines []string) []Match {
	matches := make([]Match, 0)
	for idx, line := range lines {
		if strIdx := strings.Index(line, searchString); strIdx > -1 {
			match := Match{
				lineNumber:  idx + 1,
				lineContent: line,
			}
			matches = append(matches, match)
		}
	}

	return matches
}

// will print to the user each line and number the search
// string appeared
func printMatches(matches []Match) {
	// TODO
	for _, match := range matches {
		fmt.Printf("%v: %v\n", match.lineNumber, match.lineContent)
	}
}

func handleRecursiveSearch(searchString, dirName string) error {
	file, err := os.Open(dirName)
	if err != nil {
		return err
	}
	defer file.Close()

	files, err := readFilesFromDir(file)
	for _, fileName := range files {
		execute(fileName, searchString)
	}

	return nil
}

// will receive a opened directory and return a list with
// the name of all the files in dir
func readFilesFromDir(dir *os.File) ([]string, error) {
	files, err := dir.Readdirnames(0)
	if err != nil {
		return nil, err
	}

	return files, nil
}
