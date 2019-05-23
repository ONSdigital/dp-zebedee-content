package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

var (
	serviceIDChars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	r              *rand.Rand
	services       = []string{"babbage",
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
		"dp-identity-api"}
)

func main() {

	for _, service := range services {

		account := map[string]interface{}{
			"id": service,
		}

		id := newRandomID(64)

		jsonData, err := json.Marshal(account)
		if err != nil {
			panic(err)
		}

		jsonFile, err := os.Create("./" + id + ".json")
		if err != nil {
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
