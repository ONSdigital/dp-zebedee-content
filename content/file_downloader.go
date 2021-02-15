package content

import (
	"io"
	"os"

	"github.com/ONSdigital/dp-zebedee-content/files"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/daiLlew/funkylog"
	"github.com/pkg/errors"
)

var (
	errInvalidDestFile   = errors.New("download destination required but was empty")
	errInvalidFilename   = errors.New("filename required but was empty")
	errInvalidBucketName = errors.New("s3 bucket name required but was empty")
	errDownloaderNil     = errors.New("Downloader required but was nil")
)

// Downloader defines an object for downloading something from an s3 bucket.
type Downloader interface {
	Download(w io.WriterAt, input *s3.GetObjectInput, options ...func(downloader *s3manager.Downloader)) (n int64, err error)
}

// DownloadExampleZip download the example content zip from the S3 bucket.
func DownloadExampleZip(dest, filename, bucketName string, downloader Downloader) (int64, error) {
	if err := validate(dest, filename, bucketName, downloader); err != nil {
		return 0, err
	}

	exists, errExist := files.Exists(dest)
	if errExist != nil {
		return 0, errExist
	}

	if exists {
		log.Info("content zip %s already exists skipping download", dest)
		return 0, nil
	}

	file, err := os.Create(dest)
	if err != nil {
		errors.Errorf("Unable to open file %q, %v", filename, err)
	}

	defer file.Close()

	log.Info("downloading %s from S3 bucket", filename)

	var numBytes int64
	numBytes, err = downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
	})

	// If the download error is not nil delete the dest file before returning
	if err != nil {
		os.Remove(dest)
		return 0, err
	}

	log.Info("download completed successfully")
	return numBytes, nil
}

func validate(dest, filename, bucketName string, downloader Downloader) error {
	if len(dest) == 0 {
		return errInvalidDestFile
	}

	if len(filename) == 0 {
		return errInvalidFilename
	}

	if len(bucketName) == 0 {
		return errInvalidBucketName
	}

	if downloader == nil {
		return errDownloaderNil
	}
	return nil
}
