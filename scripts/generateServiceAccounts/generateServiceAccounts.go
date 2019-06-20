package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
	"flag"
)

var (
	serviceIDChars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	r              *rand.Rand
	serviceList       = []string{"dp-import-reporter",
		"babbage",
		"florence",
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
		"dp-hierarchy-api",
		"dp-hierarchy-builder",
		"dp-dataset-exporter",
		"dp-dataset-exporter-xlsx",
		"dp-search-builder",
		"dp-search-api",
		"dp-api-router",
		"dp-download-service",
		"dp-frontend-geography-controller",
		"dp-identity-api"}}
)

func main() {

	specificService := flag.String("svc", "", "a specific service to generate an account for")
	withID := flag.String("id", "", "generate a specific service account with a user supplied ID")
	flag.Parse()

	if *specificService == "" && *withID != "" {
		fmt.Println("Aborting, you must specify a lone service to use a pre-defined ID")
		os.Exit(1)
	}

	var services []string
	if *specificService != "" {
		services = []string{*specificService}
	} else {
		services = serviceList
	}

	for _, service := range services {

		account := map[string]interface{}{
			"id": service,
		}

		var id string
		if *withID == "" {
			id = newRandomID(64)
		} else {
			id = *withID
		}

		jsonData, err := json.Marshal(account)
		if err != nil {
			fmt.Sprintf("Failed when attempting to marshal json for service: ", service)
			panic(err)
		}

		jsonFile, err := os.Create("./" + id + ".json")
		if err != nil {
			fmt.Sprintf("Failed when generating json for service: ", service)
			panic(err)
		}
		defer jsonFile.Close()
		jsonFile.Write(jsonData)

		fmt.Println(id, service)

	}

}

func init() {
	s1 := rand.NewSource(time.Now().UnixNano())
	r = rand.New(s1)
}

func newRandomID(size int) string {
	id := ""
	for len(id) < size {
		id += string(serviceIDChars[r.Intn(len(serviceIDChars))])
	}

	return id
}
