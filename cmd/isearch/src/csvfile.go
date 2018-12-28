package src

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Parser parser the file and output de result file
func Parser(input string) (content [][]string, err error) {
	file, err := ReadFile(input)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, body, err := ReadContent(file)
	if err != nil {
		fmt.Printf("fail when reading input file: %v\n", err)
		return nil, err
	}

	return body, nil
}

// ReadFile read file by input
func ReadFile(input string) (file *os.File, err error) {
	goPaths := filepath.SplitList("GOPATH")
	if len(goPaths) == 0 {
		return nil, errors.New("GOPATH environment variable is not set or empty")
	}

	goRoot := runtime.GOROOT()
	if goRoot == "" {
		return nil, errors.New("GOROOT environment variable is not set or empty")
	}

	absPath, err := filepath.Abs(input)
	if err != nil {
		return nil, err
	}

	file, err = os.Open(absPath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// ReadContent read the content of the csv file , handler func(string)
func ReadContent(file *os.File) (header []string, body [][]string, err error) {
	r := csv.NewReader(file)
	// 逐行读取
	records, err := r.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	header = records[0]
	body = records[1:]

	return header, body, nil
}
