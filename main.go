package main

import (
	"context"
	"os"

	"github.com/ONSdigital/dp-zebedee-content/cms"
	"github.com/ONSdigital/dp-zebedee-content/out"
	"github.com/ONSdigital/dp-zebedee-content/scripts"
	"github.com/ONSdigital/log.go/log"
	"github.com/spf13/cobra"
)

func main() {
	if err := run(); err != nil {
		log.Event(context.Background(), "unexpected error", log.Error(err))
		os.Exit(1)
	}
}

func run() error {
	log.Namespace = "zebedee-content"

	root := &cobra.Command{
		Use:   "content",
		Short: "Generates default website content, users, permissions and service accounts required to run an instance of Zebedee CMS",
		RunE: func(cmd *cobra.Command, args []string) error {
			contentDir, err := cmd.Flags().GetString("content")
			if err != nil {
				return err
			}

			zebedeeDir, err := cmd.Flags().GetString("zebedee")
			if err != nil {
				return err
			}

			if contentDir == "" || zebedeeDir == "" {
				cmd.Help()
				return nil
			}

			builder, err := cms.New(contentDir, true)
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

			scriptLocation, err := scripts.CopyToProjectDir(zebedeeDir, file)
			if err != nil {
				return err
			}

			out.Info("Successfully generated zebedee directory structure and default content")
			out.InfoFHighlight("A script to run Zebedee CMS has been generated and placed in your Zebedee project dir: %s", scriptLocation)

			return nil
		},
	}

	root.Flags().StringP("content", "c", "","The directory to generate website content under (required)")
	root.Flags().StringP("zebedee", "z", "", "The directory of your Zebedee project (required)")

	err := root.Execute()
	if err != nil {
		return err
	}
	return nil
}