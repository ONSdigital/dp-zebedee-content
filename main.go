package main

import (
	"context"
	"errors"
	"flag"
	"os"

	"github.com/ONSdigital/dp-zebedee-content/cms"
	"github.com/ONSdigital/dp-zebedee-content/scripts"
	"github.com/ONSdigital/log.go/log"
)

func main() {
	if err := run(); err != nil {
		log.Event(context.Background(), "unexpected error", log.Error(err))
		os.Exit(1)
	}
}

func run() error {
	log.Namespace = "zebedee-content"
	contentDir := flag.String("content_dir", "", "The directory in which to build Zebedee directory structure and unpack the default content")
	projectDir := flag.String("project_dir", "", "The root directory of your Zebedee project")
	enableCMD := flag.Bool("enable_cmd", false, "If enabled the generated run script will have the CMD features enabled")
	flag.Parse()

	if *contentDir == "" {
		return errors.New("please specify a content_dir - see help (-h) for more details")
	}

	if *projectDir == "" {
		return errors.New("please specify the project_dir - see help (-h) for more details")
	}

	builder, err := cms.New(*contentDir, *enableCMD)
	if err != nil {
		return err
	}

	err = builder.GenerateCMSContent()
	if err != nil {
		return err
	}

	t := builder.GetRunTemplate()

	var file string
	file, err = scripts.GenerateCMSRunScript(t)
	if err != nil {
		return err
	}

	scriptLocation, err := scripts.CopyToProjectDir(*projectDir, file)
	if err != nil {
		return err
	}

	log.Event(context.Background(), "successfully generated zebedee file structure and default content you can use the generated run-cms.sh file to run the application", log.Data{
		"run_script_location":      scriptLocation,
		cms.EnableCMDEnv:           t.EnableDatasetImport,
		cms.DatasetAPIAuthTokenEnv: t.DatasetAPIAuthToken,
		cms.DatasetAPIURLEnv:       t.DatasetAPIURL,
		cms.ServiceAuthTokenEnv:    t.ServiceAuthToken,
	})

	return nil
}

func errorAndExit(err error) {

}
