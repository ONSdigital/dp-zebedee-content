package cms

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/daiLlew/funkylog"
	"github.com/fatih/color"
	"github.com/pkg/errors"
)

var (
	errInvalidTargetFile   = errors.New("download target file required but was empty")
	errContentRootDirEmpty = errors.New("content root dir path required but was empty")
	errZipSrcEmpty         = errors.New("zip source file required but was empty")
	errZipDestEmpty        = errors.New("zip destination file required but was empty")
	errServicesDirEmpty    = errors.New("service dir required but was empty")

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

	defaultServiceAuthToken = "fc4089e2e12937861377629b0cd96cf79298a4c5d329a2ebb96664c88df77b67"
	zipName                 = "cms-content.zip"
	bucketName              = "ons-dp-sandbox-developer-cms-content"
	serviceTokenEnvVar      = "SERVICE_AUTH_TOKEN"
	defaultServiceAccount   = serviceAccount{
		ID: "Weyland-Yutani Corporation",
	}
)

type serviceAccount struct {
	ID string `json:"id"`
}

// Setup the CMS content.
func Setup(cmsRootDir string) error {
	if len(cmsRootDir) == 0 {
		return errContentRootDirEmpty
	}

	zebedeeDir, err := CreateDirStructure(cmsRootDir)
	if err != nil {
		return errors.WithMessage(err, "error creating CMS directory structure")
	}

	contentZip := filepath.Join(cmsRootDir, zipName)
	if err := DownloadContentZip(contentZip); err != nil {
		return errors.WithMessage(err, "error attempting to download content zip")
	}

	masterDir := filepath.Join(zebedeeDir, "master")
	if err := UnzipContent(contentZip, masterDir); err != nil {
		return errors.WithMessage(err, "error while attempting to unzip content into master dir")
	}

	servicesDir := filepath.Join(zebedeeDir, "services")
	serviceAuthToken, serviceAccErr := CreateServiceAccount(servicesDir)
	if serviceAccErr != nil {
		return errors.WithMessage(err, "error creating service auth token for CMS")
	}

	log.Info(":exclamation: Add the following to your env vars if they do not already exist :exclamation:")

	color.Yellow("\n\texport SERVICE_AUTH_TOKEN=%s", serviceAuthToken)
	color.Yellow("\texport zebedee_root=%s\n\n", cmsRootDir)

	log.Info("restart zebedee if already running and ensure the correct %q configuration is being applied (full app configuration is logged on start up)", "zebedee_root")
	log.Info("set up CMS content completed successfully :tada::rocket:")
	return nil
}

// CreateDirStructure creates the directory structure required to run Zebedee CMS.
func CreateDirStructure(cmsRootDir string) (string, error) {
	if len(cmsRootDir) == 0 {
		return "", errContentRootDirEmpty
	}

	zebedeeDir := filepath.Join(cmsRootDir, "zebedee")

	exists, err := fileExists(zebedeeDir)
	if err != nil {
		return "", err
	}

	if exists {
		log.Info("an existing zebedee content directory already exists under %q deleting and regenerating", zebedeeDir)
		if err := os.RemoveAll(zebedeeDir); err != nil {
			return "", errors.WithMessage(err, "error removing existing content_dir")
		}
	}

	log.Info("creating Zebedee content directories under: %s", zebedeeDir)
	if err := os.MkdirAll(zebedeeDir, 0700); err != nil {
		return "", err
	}

	for _, dir := range cmsDirs {
		if err := os.Mkdir(filepath.Join(zebedeeDir, dir), 0700); err != nil {
			return "", errors.WithMessage(err, fmt.Sprintf("error while attempting to create zebedee directory: %s", dir))
		}
	}

	return zebedeeDir, nil
}

// DownloadContentZip download the example CMS content zip from the S3 bucket.
func DownloadContentZip(target string) error {
	if len(target) == 0 {
		return errInvalidTargetFile
	}

	exists, errExist := fileExists(target)
	if errExist != nil {
		return errExist
	}

	if exists {
		log.Info("content zip %s already exists skipping download from s3 bucket", target)
		return nil
	}

	log.Info("downloading %s from S3 bucket", zipName)

	var numBytes int64

	err := checkLsExists()
	if err != nil {
		return err
	}

	path := "s3://" + bucketName + "/" + zipName
	cmd := exec.Command("aws", "s3", "cp", path, target)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	err = cmd.Run()
	if err != nil {
		return err
	}
	//	outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
	//	fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)

	log.Info("successfully downloaded CMS content zip %d bytes", numBytes)
	return nil
}

func checkLsExists() error {
	_, err := exec.LookPath("aws")
	if err != nil {
		return err
	}
	return nil
}

// UnzipContent unzip a zip file (src) to the specified location (dest)
func UnzipContent(src, dest string) error {
	if len(src) == 0 {
		return errZipSrcEmpty
	}

	if len(dest) == 0 {
		return errZipDestEmpty
	}

	zipR, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer zipR.Close()

	extractAndWriteFile := func(entry *zip.File) error {
		fpath := filepath.Join(dest, entry.Name)

		// If dir create full dir path and return
		if entry.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, os.ModeDir|os.ModePerm); err != nil {
				return err
			}
			return nil
		}

		// Otherwise its a file...
		err = os.MkdirAll(filepath.Dir(fpath), os.ModeDir|os.ModePerm)
		if err != nil {
			return err
		}

		f, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, entry.Mode())
		if err != nil {
			return err
		}

		defer f.Close()

		rc, errOpenEntry := entry.Open()
		if errOpenEntry != nil {
			return err
		}

		defer rc.Close()

		_, err = io.Copy(f, rc)
		if err != nil {
			return err
		}

		return nil
	}

	for _, entry := range zipR.File {
		if err := extractAndWriteFile(entry); err != nil {
			return err
		}
	}

	log.Info("successfully unzipped content under %s", dest)
	return nil
}

// CreateServiceAccount create a service account Zebedee CMS. If env var SERVICE_AUTH_TOKEN already exists a service
// account with this ID will be created, if env SERVICE_AUTH_TOKEN does not exist a new ID & service account will be generated.
func CreateServiceAccount(servicesDir string) (string, error) {
	if len(servicesDir) == 0 {
		return "", errServicesDirEmpty
	}

	exists, existsErr := fileExists(servicesDir)
	if existsErr != nil {
		return "", existsErr
	}

	if !exists {
		return "", errors.Errorf("servicesDir does not exist %s", servicesDir)
	}

	jsonB, err := json.Marshal(defaultServiceAccount)
	if err != nil {
		return "", errors.Wrap(err, "error marshaling service account JSON")
	}

	serviceAuthToken := os.Getenv(serviceTokenEnvVar)
	if len(serviceAuthToken) == 0 {
		log.Info("env var %q not found creating a new service account", serviceTokenEnvVar)
		serviceAuthToken = defaultServiceAuthToken
	} else {
		log.Info("creating service account from existing %q env var", serviceTokenEnvVar)
	}

	filename := filepath.Join(servicesDir, serviceAuthToken+".json")
	err = ioutil.WriteFile(filename, jsonB, 0644)
	if err != nil {
		return "", errors.Wrap(err, "error writing service account JSON to file")
	}

	log.Info("create service account completed successfully")
	return serviceAuthToken, nil
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
