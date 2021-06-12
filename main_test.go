package main

import (
	"log"
	"os"
	"testing"
)

const (
	DIR_NAME       = "dir_test"
	FILE1          = DIR_NAME + "/file1"
	FILE_1_CONTENT = "this is the content of the file 1"
	FILE2          = DIR_NAME + "/file2"
	FILE_2_CONTENT = "this is the content of the file 2"
	FILE3          = DIR_NAME + "/file3"
	FILE_3_CONTENT = "this is the content of the file 3\nthis is a new line on the third file"
	FILES_QTD      = 3
)

func TestMain(m *testing.M) {
	createFiles()
	code := m.Run()
	cleanFiles()
	os.Exit(code)
}

func cleanFiles() {
	removeFile(FILE1)
	removeFile(FILE2)
	removeFile(FILE3)
	err := os.Remove(DIR_NAME)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
}

func removeFile(fileName string) {
	err := os.Remove(fileName)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
}

func createFiles() {
	err := os.Mkdir(DIR_NAME, os.ModePerm)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	createFile(FILE1, FILE_1_CONTENT)
	createFile(FILE2, FILE_2_CONTENT)
	createFile(FILE3, FILE_3_CONTENT)
}

func createFile(fileName, fileContent string) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	_, err = file.Write([]byte(fileContent))
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	file.Close()
}

func TestReadFilesFromDir(t *testing.T) {
	file, err := os.Open(DIR_NAME)
	if err != nil {
		t.Fatalf("Error opening dir %v\n", err)
	}
	files, err := readFilesFromDir(file)
	if err != nil {
		t.Fatalf("Error reading files from dir %v\n", err)
	}

	if len(files) != FILES_QTD {
		t.Fatalf("Wrong number of file names returned, expect %v got %v", FILES_QTD, len(files))
	}
}

func TestSearchForMacthes(t *testing.T) {
	lines := []string{"this is the first line", "now the second line", "and here is the thrid line"}
	searchString := "the"
	matches := searchForMatches("randomFile", searchString, lines)
	if len(matches) != 3 {
		t.Fatalf("Wrong matches returned. exepct 3 got %v\n", len(matches))
	}

	match := matches[0]
	if match.lineNumber != 1 {
		t.Fatalf("Wrong line number returned. exepct 1 got %v\n", match.lineNumber)
	}
	match = matches[1]
	if match.lineNumber != 2 {
		t.Fatalf("Wrong line number returned. exepct 2 got %v\n", match.lineNumber)
	}
	match = matches[2]
	if match.lineNumber != 3 {
		t.Fatalf("Wrong line number returned. exepct 3 got %v\n", match.lineNumber)
	}

}

func TestReadFileLines(t *testing.T) {
	lines, err := readFileLines(FILE1)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	if len(lines) != 1 {
		t.Fatalf("Wrong number of lines from %v. Expect 1 got %v", FILE1, len(lines))
	}

	lines, err = readFileLines(FILE3)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	if len(lines) != 2 {
		t.Fatalf("Wrong number of lines from %v. Expect 2 got %v", FILE3, len(lines))
	}

}

func TestReadArguments(t *testing.T) {
	os.Args = []string{"program", "searchString", "fileName"}
	args, err := readArguments()
	if err != nil {
		t.Fatalf("Error reading arguments%v\n", err)
	}
	if len(args) != 3 {
		t.Fatalf("Wrong arguments returned. Expect 3 got %v\n", len(args))
	}

	os.Args = []string{"program", "searchString", "fileName", "-R"}
	args, _ = readArguments()
	if len(args) != 4 {
		t.Fatalf("Wrong arguments returned. Expect 4 got %v\n", len(args))
	}

	os.Args = []string{"program", "searchString", "fileName", "-t"}
	_, err = readArguments()
	if err == nil {
		t.Fatalf("Should return error\n")
	}
}
