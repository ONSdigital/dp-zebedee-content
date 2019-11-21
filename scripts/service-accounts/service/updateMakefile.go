package service

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"regexp"
)

const stringMatch string = "export GOARCH?=$(shell go env GOARCH)"

var existingServiceAuthToken, _ = regexp.Compile("export SERVICE_AUTH_TOKEN=[a-zA-Z0-9]+")

// UpdateMakefile adds environment variable SERVICE_AUTH_TOKEN
// to a services Makefile
func UpdateMakefile(filePath, serviceAuthToken string) error {
	lines, err := file2lines(filePath)
	if err != nil {
		return err
	}

	fileContent := ""
	for _, line := range lines {
		if !existingServiceAuthToken.MatchString(line) {
			fileContent += line + "\n"
		}

		if line == stringMatch {
			fileContent += "export SERVICE_AUTH_TOKEN=" + serviceAuthToken + "\n"
		}
	}

	return ioutil.WriteFile(filePath, []byte(fileContent), 0644)
}

func file2lines(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return linesFromReader(f)
}

func linesFromReader(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
