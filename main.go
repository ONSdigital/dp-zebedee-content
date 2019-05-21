package main

import (
	"flag"
	"os"

	"github.com/ONSdigital/dp-zebedee-content/cms"
	"github.com/ONSdigital/dp-zebedee-content/scripts"
	"github.com/ONSdigital/log.go/log"
)

func main() {
	log.Namespace = "zebedee-content"
	contentDir := flag.String("content_dir", "", "The directory in which to build Zebedee directory structure and unpack the default content")
	projectDir := flag.String("project_dir", "", "The root directory of your Zebedee project")
	enableCMD := flag.Bool("enable_cmd", false, "If enabled the generated run script will have the CMD features enabled")
	flag.Parse()

	if *contentDir == "" {
		log.Event(nil, "please specify a content_dir - see help (-h) for more details")
		os.Exit(1)
	}

	if *projectDir == "" {
		log.Event(nil, "please specify the project_dir - see help (-h) for more details")
		os.Exit(1)
	}

	builder, err := cms.New(*contentDir, *enableCMD)
	if err != nil {
		errorAndExit(err)
	}

	err = builder.GenerateCMSContent()
	if err != nil {
		errorAndExit(err)
	}

	t := builder.GetRunTemplate()

	var file string
	file, err = scripts.GenerateCMSRunScript(t)
	if err != nil {
		errorAndExit(err)
	}

	scriptLocation, err := scripts.CopyToProjectDir(*projectDir, file)
	if err != nil {
		errorAndExit(err)
	}
	log.Event(nil, "successfully generated zebedee file structure and default content you can use the generated run-cms.sh file to run the application", log.Data{
		"run_script_location":      scriptLocation,
		cms.EnableCMDEnv:           t.EnableDatasetImport,
		cms.DatasetAPIAuthTokenEnv: t.DatasetAPIAuthToken,
		cms.DatasetAPIURLEnv:       t.DatasetAPIURL,
		cms.ServiceAuthTokenEnv:    t.ServiceAuthToken,
	})
}

func errorAndExit(err error) {
	log.Event(nil, "unexpected error", log.Error(err))
	os.Exit(1)
}
