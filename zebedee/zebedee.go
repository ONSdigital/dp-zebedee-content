package zebedee

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ONSdigital/dp-zebedee-content/content"
	"github.com/ONSdigital/dp-zebedee-content/files"
	log "github.com/daiLlew/funkylog"
	"github.com/pkg/errors"
)

var (
	RootEmptyErr            = errors.New("root dir path required but was empty")
	defaultServiceAuthToken = "fc4089e2e12937861377629b0cd96cf79298a4c5d329a2ebb96664c88df77b67"

	cmsDirs = []string{
		"master",
		"collections",
		"publishing-log",
		"users",
		"sessions",
		"services",
		"permissions",
		"teams",
		"transactions",
		"launchpad",
		"application-keys",
		"keyring",
	}
)

// CreateDirStructure TODO
func CreateDirStructure(path string) error {
	if len(path) == 0 {
		return RootEmptyErr
	}

	exists, err := files.Exists(path)
	if err != nil {
		return err
	}

	if exists {
		log.Info("removing exists zebedee content directory")
		if err := os.RemoveAll(path); err != nil {
			return errors.WithMessage(err, "error removing existing content_dir")
		}
	}

	log.Info("creating fresh Zebedee content directory structure")
	if err := os.MkdirAll(path, 0700); err != nil {
		return err
	}

	for _, dir := range cmsDirs {
		if err := os.Mkdir(filepath.Join(path, dir), 0700); err != nil {
			return errors.WithMessage(err, fmt.Sprintf("error while attempting to create zebedee directory: %s", dir))
		}
	}

	return nil
}

func DownloadAndUnzipExampleContent(root string, downloader content.Downloader) error {
	zipFile := filepath.Join(root, "example-cms-content.zip")
	if _, err := content.DownloadExampleZip(zipFile, "example-cms-content.zip", "developer-cms-content", downloader); err != nil {
		return err
	}

	unzipDest := filepath.Join(root, "zebedee", "master")
	log.Info("unzipping %s to %s", zipFile, unzipDest)

	cmd := exec.Command("unzip", "-q", "-o", zipFile, "-d", unzipDest)
	cmd.Stdout = os.Stdout
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return errors.WithMessage(err, "error unzipping content zip")
	}

	return nil
}

func CreateServiceAccount(servicesDir string) (string, error) {
	log.Info("Generating new CMD service account")

	jsonB, err := json.Marshal(map[string]interface{}{"id": "Weyland-Yutani Corporation"})
	if err != nil {
		return "", errors.Wrap(err, "error marshaling service account JSON")
	}

	serviceAuthToken := os.Getenv("SERVICE_AUTH_TOKEN")
	if len(serviceAuthToken) == 0 {
		serviceAuthToken = defaultServiceAuthToken
	}

	filename := filepath.Join(servicesDir, serviceAuthToken+".json")
	err = ioutil.WriteFile(filename, jsonB, 0644)
	if err != nil {
		return "", errors.Wrap(err, "error writing service account JSON to file")
	}

	return serviceAuthToken, nil
}
