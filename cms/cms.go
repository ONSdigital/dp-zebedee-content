package cms

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ONSdigital/dp-zebedee-content/files"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/daiLlew/funkylog"
	"github.com/pkg/errors"
)

var (
	errInvalidTargetFile   = errors.New("download target file required but was empty")
	errDownloaderNil       = errors.New("Downloader required but was nil")
	errContentRootDirEmpty = errors.New("content root dir path required but was empty")
	errZipSrcEmpty         = errors.New("zip source file required but was empty")
	errZipDestEmpty        = errors.New("zip destination file required but was empty")

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
	bucketName              = "developer-cms-content"
)

// Downloader defines an object for downloading something from an s3 bucket.
type Downloader interface {
	Download(w io.WriterAt, input *s3.GetObjectInput, options ...func(downloader *s3manager.Downloader)) (n int64, err error)
}

func Setup(cmsRootDir string, downloader Downloader) error {
	if len(cmsRootDir) == 0 {
		return errContentRootDirEmpty
	}

	zebedeeDir, err := CreateDirStructure(cmsRootDir)
	if err != nil {
		return errors.WithMessage(err, "error creating CMS directory structure")
	}

	contentZip := filepath.Join(cmsRootDir, zipName)
	if err := DownloadContentZip(contentZip, downloader); err != nil {
		return errors.WithMessage(err, "error attempting to download content zip")
	}

	masterDir := filepath.Join(zebedeeDir, "master")
	if err := UnzipContent(contentZip, masterDir); err != nil {
		return errors.WithMessage(err, "error while attempting to unzip content into master dir")
	}

	servicesDir := filepath.Join(zebedeeDir, "services")
	_, err = CreateServiceAccount(servicesDir)
	if err != nil {
		return errors.WithMessage(err, "error creating service auth token for CMS")
	}

	return nil
}

// CreateDirStructure TODO
func CreateDirStructure(cmsRootDir string) (string, error) {
	if len(cmsRootDir) == 0 {
		return "", errContentRootDirEmpty
	}

	zebedeeDir := filepath.Join(cmsRootDir, "zebedee")

	exists, err := files.Exists(zebedeeDir)
	if err != nil {
		return "", err
	}

	if exists {
		log.Info("removing exists zebedee content directory")
		if err := os.RemoveAll(zebedeeDir); err != nil {
			return "", errors.WithMessage(err, "error removing existing content_dir")
		}
	}

	log.Info("creating fresh Zebedee content directory structure")
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

// DownloadContentZip download the example content zip from the S3 bucket.
func DownloadContentZip(target string, downloader Downloader) error {
	if len(target) == 0 {
		return errInvalidTargetFile
	}

	if downloader == nil {
		return errDownloaderNil
	}

	exists, errExist := files.Exists(target)
	if errExist != nil {
		return errExist
	}

	if exists {
		log.Info("content zip %s already exists skipping download", target)
		return nil
	}

	downloadedFile, err := os.Create(target)
	if err != nil {
		errors.Errorf("Unable to open file %q, %v", target, err)
	}

	defer downloadedFile.Close()

	log.Info("downloading %s from S3 bucket", zipName)

	var numBytes int64
	numBytes, err = downloader.Download(downloadedFile, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(zipName),
	})

	// If the download error is not nil delete the dest file before returning
	if err != nil {
		os.Remove(target)
		return err
	}

	log.Info("successfully downloaded website content zip %d bytes", numBytes)
	return nil
}

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

	for _, f := range zipR.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}

		defer rc.Close()

		path := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(path, os.ModeDir|os.ModePerm); err != nil {
				return err
			}
			continue
		}

		err = os.MkdirAll(filepath.Dir(path), os.ModeDir|os.ModePerm)
		if err != nil {
			return err
		}

		f, err := os.Create(path)
		if err != nil {
			return err
		}

		defer f.Close()

		_, err = io.Copy(f, rc)
		if err != nil {
			return err
		}
	}

	log.Info("successfully unzipped content to %s", dest)
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
