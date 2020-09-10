package commands

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ONSdigital/dp-zebedee-content/cms"
	"github.com/ONSdigital/dp-zebedee-content/out"
	"github.com/ONSdigital/dp-zebedee-content/scripts"
	"github.com/spf13/cobra"
)

const (
	tilde = "~"
	home  = "HOME"

	contentFlag = "content"
	zebedeeFlag = "zebedee"
)

func GetRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "dp-zebedee-content",
		Short: "Cli tool for ONS website developers. Generates default content and directory structure required to run an instance of Zebedee CMS.",
	}

	root.AddCommand(getGenerateCommand())
	return root
}

func getGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate the Zebedee directory structure and populate with default web content, users accounts, user permissions, teams and service accounts.",
		RunE: func(cmd *cobra.Command, args []string) error {
			contentDir, err := getPathFlag(contentFlag, cmd)
			if err != nil {
				return err
			}

			zebedeeDir, err := getPathFlag(zebedeeFlag, cmd)
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

			out.InfoFHighlight("Successfully generated Zebedee CMS file structure and default content: %s", contentDir)
			out.InfoFHighlight("Successfully generated Zebedee CMS run script: %s", scriptLocation)
			return nil
		},
	}

	cmd.Flags().StringP(contentFlag, "c", "", "The directory to generate the default Zebedee CMS content under (Required)")
	cmd.Flags().StringP(zebedeeFlag, "z", "", "The directory of your Zebedee project (Required)")

	return cmd
}

func getPathFlag(name string, cmd *cobra.Command) (string, error) {
	p, err := cmd.Flags().GetString(name)
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(p, tilde) {
		p = strings.TrimPrefix(p, tilde)
		p = filepath.Join(os.Getenv(home), p)
	}

	return p, nil
}
