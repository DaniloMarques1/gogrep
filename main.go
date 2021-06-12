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
		err = handleRecursiveSearch(searchString, fileName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func execute(fileName, searchString string) {
	lines, err := readFileLines(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	matches := searchForMatches(fileName, searchString, lines)
	printMatches(matches)
}

// will read and return all arguments given
func readArguments() ([]string, error) {
	args := os.Args
	if len(args) < 3 {
		return nil, errors.New("gogrep: Not enough arguments")
	}

	if len(args) == 4 && args[3] != "-R" {
		return nil, errors.New("gogrep: Wrong flag given")
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
func searchForMatches(fileName, searchString string, lines []string) []Match {
	matches := make([]Match, 0)
	for idx, line := range lines {
		if strIdx := strings.Index(line, searchString); strIdx > -1 {
			match := Match{
				fileName:    fileName,
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
	for _, match := range matches {
		fmt.Printf("%v:%v: %v\n", match.fileName, match.lineNumber, match.lineContent)
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
		fileInfo, err := os.Stat(dirName + "/" + fileName)
		if err != nil {
			return err
		}

		fullPath := dirName + "/" + fileName
		if fileInfo.IsDir() {
			err = handleRecursiveSearch(searchString, fullPath)
			if err != nil {
				return err
			}
		} else {
			execute(fullPath, searchString)
		}
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
