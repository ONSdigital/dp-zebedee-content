package main

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/ONSdigital/dp-zebedee-content/scripts/service-accounts/service"
	"github.com/ONSdigital/log.go/log"
)

var services = []string{
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

var r *rand.Rand

func init() {
	s1 := rand.NewSource(time.Now().UnixNano())
	r = rand.New(s1)
}

func main() {
	replaceDir := flag.Bool("replace", false, "the service account directory path to be replaced")
	pathToServices := flag.String("set-path", "", "path to services")
	serviceDir := flag.String("dir", "/zebedee/services", "the service account directory path")
	updateMakefiles := flag.Bool("update-mk", false, "an indicator to update repo makefiles to add service auth token, must be set to true or false, default value is false")

	flag.Parse()

	if err := service.CreateDirectory(*serviceDir, *replaceDir); err != nil {
		log.Event(nil, "error checking service dir", log.Error(err))
		os.Exit(1)
	}

	if *updateMakefiles {
		// Check path to services is being set
		if *pathToServices == "" {
			err := errors.New("missing path to services")
			log.Event(nil, "unable to update makefiles, missing path to services", log.Error(err))
			os.Exit(1)
		}
	}

	existingAccounts, err := service.LoadExisting(*serviceDir)
	if err != nil {
		log.Event(nil, "error loading existing service accounts", log.Error(err))
		os.Exit(1)
	}

	if len(existingAccounts) > 0 {
		log.Event(nil, "found existing service accounts skipped", log.Data{"existing_accounts": existingAccounts})
	}

	for _, s := range services {
		logD := log.Data{"service": s}

		if _, ok := existingAccounts[s]; ok {
			continue
		}

		acc, err := service.CreateJSONAccount(r, *serviceDir, s)
		if err != nil {
			log.Event(nil, "error creating service account", log.Error(err), logD)
			os.Exit(1)
		}

		if *updateMakefiles {
			filePath := *pathToServices + "/" + s + "/Makefile"

			if err := service.UpdateMakefile(filePath, acc.ID); err != nil {
				log.Event(nil, "unable to update makefile for service", log.Error(err), logD)
			}
		}

		fmt.Printf("%s\t %s\n", acc.ID, acc.Name)
	}

	log.Event(nil, "successfully created service accounts")
}
