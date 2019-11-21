package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/ONSdigital/log.go/log"
)

const stringMatch string = "export GOARCH?=$(shell go env GOARCH)"

var (
	existingServiceAuthToken, _ = regexp.Compile("export SERVICE_AUTH_TOKEN=[a-zA-Z0-9]+")
	serviceIDChars              = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	r                           *rand.Rand
	services                    = []string{
		"babbage",
		"florence",
		"zebedee",
		"dp-frontend-router",
		"dp-frontend-filter-dataset-controller",
		"dp-frontend-renderer",
		"dp-frontend-dataset-controller",
		"dp-import-tracker",
		"dp-dimension-extractor",
		"dp-dimension-importer",
		"dp-observation-extractor",
		"dp-observation-importer",
		"dp-import-api",
		"dp-code-list-api",
		"dp-dataset-api",
		"dp-filter-api",
		"dp-recipe-api",
		"dp-code-list-api",
		"dp-dataset-exporter",
		"dp-hierarchy-api",
		"dp-hierarchy-builder",
		"dp-dataset-exporter",
		"dp-dataset-exporter-xlsx",
		"dp-search-builder",
		"dp-search-api",
		"dp-api-router",
		"dp-download-service",
		"dp-frontend-geography-controller",
		"dp-identity-api",
	}
)

type ServiceAccount struct {
	ID   string `json:"-"`
	Name string `json:"id,omitempty"`
}

func init() {
	s1 := rand.NewSource(time.Now().UnixNano())
	r = rand.New(s1)
}

func main() {
	replaceDir := flag.Bool("replace", false, "the service account directory path to be replaced")
	pathToServices := flag.String("set-path", "", "path to services")
	serviceDir := flag.String("dir", "/zebedee/services", "the service account directory path")
	updateMakeFiles := flag.Bool("update-mk", false, "an indicator to update repo makefiles to add service auth token, must be set to true or false, default value is false")

	flag.Parse()

	if err := createServiceDir(*serviceDir, *replaceDir); err != nil {
		log.Event(nil, "error checking service dir", log.Error(err))
		os.Exit(1)
	}

	if *updateMakeFiles {
		// Check path to services is being set
		if *pathToServices == "" {
			err := errors.New("missing path to services")
			log.Event(nil, "unable to update makefiles, missing path to services", log.Error(err))
			os.Exit(1)
		}
	}

	existingAccounts, err := loadExisting(*serviceDir)
	if err != nil {
		log.Event(nil, "error loading existing service accounts", log.Error(err))
		os.Exit(1)
	}

	if len(existingAccounts) > 0 {
		log.Event(nil, "found existing service accounts skipped", log.Data{"existing_accounts": existingAccounts})
	}

	for _, service := range services {
		logD := log.Data{"service": service}

		if _, ok := existingAccounts[service]; ok {
			continue
		}

		acc, err := createServiceAccountJSON(*serviceDir, service)
		if err != nil {
			log.Event(nil, "error creating service account", log.Error(err), logD)
			os.Exit(1)
		}

		if *updateMakeFiles {
			filePath := *pathToServices + "/" + service + "/Makefile"

			if err := updateMakeFile(filePath, acc.ID); err != nil {
				log.Event(nil, "unable to update makefile for service", log.Error(err), logD)
			}
		}

		fmt.Printf("%s\t %s\n", acc.ID, acc.Name)
	}

	log.Event(nil, "successfully created service accounts")
}

func updateMakeFile(filePath, serviceAuthToken string) error {
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

func createServiceAccountJSON(path, service string) (*ServiceAccount, error) {
	serviceAcc := &ServiceAccount{
		ID:   newRandomID(64),
		Name: service,
	}

	logD := log.Data{
		"service": service,
	}

	jsonData, err := json.Marshal(serviceAcc)
	if err != nil {
		log.Event(nil, "error attempting to marshal service account to json", log.Error(err), logD)
		return nil, err
	}

	p := filepath.Join(path, serviceAcc.ID+".json")
	jsonFile, err := os.Create(p)
	if err != nil {
		log.Event(nil, "error creating service account json", log.Error(err), logD)
		return nil, err
	}
	defer func() {
		err := jsonFile.Close()
		if err != nil {
			log.Event(nil, "error closing service account json", log.Error(err), logD)
		}
	}()

	_, err = jsonFile.Write(jsonData)
	if err != nil {
		log.Event(nil, "error writing service account to file", log.Error(err), logD)
		return nil, err
	}
	return serviceAcc, nil
}

func createServiceDir(path string, replaceServiceDir bool) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Event(nil, "services dir does not exist, now creating")
		if err = os.MkdirAll(path, 0777); err != nil {
			return err
		}

		return nil
	}

	if replaceServiceDir {
		// Remove existing services directory
		log.Event(nil, "removing directory")
		os.RemoveAll(path + "/")
	}

	log.Event(nil, "services dir does not exist, now creating")
	if err := os.MkdirAll(path, 0777); err != nil {
		return err
	}

	return nil
}

func loadExisting(servicePath string) (map[string]*ServiceAccount, error) {
	results := make(map[string]*ServiceAccount, 0)

	err := filepath.Walk(servicePath, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".json" {
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			var acc ServiceAccount
			if err := json.Unmarshal(b, &acc); err != nil {
				return err
			}
			results[acc.Name] = &acc
		}
		return nil
	})

	return results, err
}

func newRandomID(size int) string {
	id := ""
	for len(id) < size {
		id += string(serviceIDChars[r.Intn(len(serviceIDChars))])
	}

	return id
}
