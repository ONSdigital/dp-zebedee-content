package commands

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ONSdigital/dp-zebedee-content/cms"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spf13/cobra"
)

const (
	tilde       = "~"
	home        = "HOME"
	contentFlag = "content"
)

func GetRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "dp-zebedee-content",
		Short: "Cli tool for ONS website developers. Generates default content and directory structure required to run an instance of Zebedee CMS.",
	}

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
	})

	downloader := s3manager.NewDownloader(sess)
	cmd := getGenerateCommand(downloader)

	root.AddCommand(cmd)
	return root
}

func getGenerateCommand(downloader *s3manager.Downloader) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate the Zebedee directory structure and populate with default web content, users accounts, user permissions, teams and service accounts.",
		RunE: func(cmd *cobra.Command, args []string) error {
			contentRootDir, err := getPathFlag(contentFlag, cmd)
			if err != nil {
				return err
			}

			if len(contentRootDir) == 0 {
				cmd.Help()
				return nil
			}

			return cms.Setup(contentRootDir, downloader)
		},
	}

	cmd.Flags().StringP(contentFlag, "c", "", "The directory to generate the default Zebedee CMS content under (Required)")

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
