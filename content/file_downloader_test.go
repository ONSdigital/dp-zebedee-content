package content

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/ONSdigital/dp-zebedee-content/content/mocks"
	"github.com/ONSdigital/dp-zebedee-content/files"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_DownloadExampleZip(t *testing.T) {
	Convey("should return an error if dest is empty", t, func() {
		nunBytes, err := DownloadExampleZip("", "", "", nil)

		So(err, ShouldEqual, errInvalidDestFile)
		So(nunBytes, ShouldEqual, 0)
	})

	Convey("should return an error if filename is empty", t, func() {
		nunBytes, err := DownloadExampleZip("dest", "", "", nil)

		So(err, ShouldEqual, errInvalidFilename)
		So(nunBytes, ShouldEqual, 0)
	})

	Convey("should return an error if bucket name is empty", t, func() {
		nunBytes, err := DownloadExampleZip("dest", "filename", "", nil)

		So(err, ShouldEqual, errInvalidBucketName)
		So(nunBytes, ShouldEqual, 0)
	})

	Convey("should return an error if downloader is nil", t, func() {
		nunBytes, err := DownloadExampleZip("content.zip", "filename", "bucket", nil)

		So(err, ShouldEqual, errDownloaderNil)
		So(nunBytes, ShouldEqual, 0)
	})

	Convey("should return an error if downloader.Download returns an error", t, func() {
		expectedErr := errors.New("unexpected downloader error")
		downloaderMock := mocks.ErroringDownloader(expectedErr)

		nunBytes, err := DownloadExampleZip("content.zip", "filename", "bucket", downloaderMock)

		So(err, ShouldEqual, expectedErr)
		So(nunBytes, ShouldEqual, 0)
		So(downloaderMock.GetCalls(), ShouldHaveLength, 1)
		So(downloaderMock.GetCalls()[0].Bucket, ShouldResemble, aws.String("bucket"))
		So(downloaderMock.GetCalls()[0].Key, ShouldResemble, aws.String("filename"))

		exists, existsErr := files.Exists("content.zip")
		So(exists, ShouldBeFalse)
		So(existsErr, ShouldBeNil)
	})

	Convey("Given valid dest, filename and bucket parameters ", t, func() {
		 var actualNumBytes int64

		downloaderMock := mocks.SuccessfulDownloader(func() (int64, error) {
			f, err := os.Create("content.txt")
			if err != nil {
				return 0, err
			}

			var n int
			n, err = f.Write([]byte("hello world!"))
			actualNumBytes = int64(n)

			return actualNumBytes, err
		})

		Convey("When DownloadExampleZip is called", func() {

			nunBytes, err := DownloadExampleZip("content.txt", "filename", "bucket", downloaderMock)
			defer os.Remove("content.txt")

			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
				So(nunBytes, ShouldEqual, actualNumBytes)

				So(downloaderMock.GetCalls(), ShouldHaveLength, 1)
				So(downloaderMock.GetCalls()[0].Bucket, ShouldResemble, aws.String("bucket"))
				So(downloaderMock.GetCalls()[0].Key, ShouldResemble, aws.String("filename"))
			})

			Convey("And the requested content is downloaded to the destination file", func() {
				exists, existsErr := files.Exists("content.txt")
				So(existsErr, ShouldBeNil)
				So(exists, ShouldBeTrue)

				b, readFileErr := ioutil.ReadFile("content.txt")
				So(readFileErr, ShouldBeNil)
				So(string(b), ShouldEqual, "hello world!")
			})
		})

	})
}
