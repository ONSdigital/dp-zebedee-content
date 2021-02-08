package cms

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ONSdigital/dp-zebedee-content/files"
	"github.com/ONSdigital/dp-zebedee-content/out"
	"github.com/pkg/errors"
)

var (
	defaultDatasetAPIAuthToken = "FD0108EA-825D-411C-9B1D-41EF7727F465"
	defaultDatasetAPIURL       = "http://localhost:22000"
	defaultServiceAuthToken    = "fc4089e2e12937861377629b0cd96cf79298a4c5d329a2ebb96664c88df77b67"
)

// Build creates the Zebedee CMS directory structure
func (b *Builder) GenerateCMSContent() error {
	out.Info("Generating CMS file structure and default content")

	b.serviceAccountID = defaultServiceAuthToken
	b.datasetAPIAuthToken = defaultDatasetAPIAuthToken
	b.datasetAPIURL = defaultDatasetAPIURL

	if err := b.cleanAndPrepare(b.rootDir); err != nil {
		return err
	}

	if err := b.createDirs(); err != nil {
		return err
	}

	err := b.copyContentZipToMaster()
	if err != nil {
		return err
	}

	err = b.unzipContentInMaster()
	if err != nil {
		return err
	}

	err = b.removeContentZipFromMaster()
	if err != nil {
		return err
	}

	err = b.createServiceAccount()
	if err != nil {
		return err
	}
	return nil
}

func (b *Builder) cleanAndPrepare(contentDir string) error {
	exists, err := files.Exists(contentDir)
	if err != nil {
		return err
	}

	if exists {
		out.Info("Removing exists Zebedee content directory")
		if err := os.RemoveAll(contentDir); err != nil {
			return errors.WithMessage(err, "error removing existing content_dir")
		}
	}

	out.Info("Generating Zebedee content directory")
	if err := os.MkdirAll(contentDir, 0700); err != nil {
		return err
	}
	return nil
}

func (b *Builder) createDirs() error {
	for _, dir := range b.dirs() {
		cmd := exec.Command("mkdir", dir)
		cmd.Stderr = b.Out
		cmd.Stdout = b.OutErr

		if err := cmd.Run(); err != nil {
			return errors.Wrap(err, fmt.Sprintf("error while attempting to create zebedee directory: %s", dir))
		}
	}

	out.Info("Zebedee directory structure created successfully")
	return nil
}

func (b *Builder) copyContentZipToMaster() error {
	cmd := newCommand("cp", "", defaultContentZip, b.masterDir)

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "error executing command: copyContentZipToMaster")
	}
	return nil
}

func (b *Builder) unzipContentInMaster() error {
	out.Info("Populating Zebedee master directory with default content")
	cmd := newCommand("unzip", b.masterDir, "-q", defaultContentZip)

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "error executing command: unzipContentInMaster")
	}
	return nil
}

func (b *Builder) removeContentZipFromMaster() error {
	cmd := newCommand("rm", b.masterDir, defaultContentZip)

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "error executing command: removeContentZipFromMaster")
	}
	return nil
}

func (b *Builder) createServiceAccount() error {
	out.Info("Generating new CMD service account")

	jsonB, err := json.Marshal(map[string]interface{}{"id": "Weyland-Yutani Corporation"})
	if err != nil {
		return errors.Wrap(err, "error marshaling service account JSON")
	}

	filename := filepath.Join(b.servicesDir, b.serviceAccountID+".json")
	err = ioutil.WriteFile(filename, jsonB, 0644)
	if err != nil {
		return errors.Wrap(err, "error writing service account JSON to file")
	}

	return nil
}

func (b *Builder) dirs() []string {
	return []string{
		b.zebedeeDir,
		b.masterDir,
		b.collectionsDir,
		b.publishLogDir,
		b.usersDir,
		b.sessionsDir,
		b.permissionsDir,
		b.teamsDir,
		b.transactionsDir,
		b.launchPadDir,
		b.appKeysDir,
		b.servicesDir,
		b.keyRingDir,
	}
}

func newCommand(name string, dir string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}
