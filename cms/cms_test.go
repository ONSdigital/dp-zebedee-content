package cms

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/ONSdigital/dp-zebedee-content/cms/mocks"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateDirStructure(t *testing.T) {
	Convey("should return error if root is empty", t, func() {
		zebedeeDir, err := CreateDirStructure("")
		So(err, ShouldEqual, errContentRootDirEmpty)
		So(zebedeeDir, ShouldBeEmpty)
	})

	Convey("should create expected directories given a valid root path", t, func() {
		rootDir := "test-content"

		zebedeeDir, err := CreateDirStructure(rootDir)
		defer tearDownCMSDirs(rootDir)

		So(err, ShouldBeNil)
		So(zebedeeDir, ShouldEqual, filepath.Join(rootDir, "zebedee"))

		for _, dir := range cmsDirs {
			path := filepath.Join(zebedeeDir, dir)
			exists, err := fileExists(path)
			So(err, ShouldBeNil)
			So(exists, ShouldBeTrue)
		}
	})
}

func TestDownloadContentZip(t *testing.T) {
	contentZip := "content.zip"

	Convey("should return an error if target is empty", t, func() {
		err := DownloadContentZip("", nil)
		So(err, ShouldEqual, errInvalidTargetFile)
	})

	Convey("should return an error if downloader is nil", t, func() {
		err := DownloadContentZip(contentZip, nil)
		So(err, ShouldEqual, errDownloaderNil)
	})

	Convey("should return nil if the content.zip already exists", t, func() {
		downloader := &mocks.DownloaderMock{}
		existingContent, err := os.Create(contentZip)
		So(err, ShouldBeNil)

		// test tear down
		defer func() {
			existingContent.Close()
			os.Remove(contentZip)
		}()

		err = DownloadContentZip(contentZip, downloader)

		So(err, ShouldBeNil)
		So(downloader.GetCalls(), ShouldHaveLength, 0)
	})

	Convey("should return error downloader fails to download the content.zip", t, func() {
		expectedErr := errors.New("download failed")
		downloader := mocks.ErroringDownloader(expectedErr)

		err := DownloadContentZip(contentZip, downloader)

		So(err, ShouldEqual, expectedErr)
		So(downloader.GetCalls(), ShouldHaveLength, 1)

		exists, existsErr := fileExists(contentZip)
		So(existsErr, ShouldBeNil)
		So(exists, ShouldBeFalse)
	})

	Convey("should return nil if the downloader successfully downloads the content.zip", t, func() {
		downloader := mocks.SuccessfulDownloader(func() (int64, error) {
			err := ioutil.WriteFile(contentZip, []byte("hello world"), 0700)
			return 0, err
		})

		defer os.Remove(contentZip)

		err := DownloadContentZip(contentZip, downloader)

		So(err, ShouldBeNil)
		So(downloader.GetCalls(), ShouldHaveLength, 1)

		exists, existsErr := fileExists(contentZip)
		So(existsErr, ShouldBeNil)
		So(exists, ShouldBeTrue)

		b, readErr := ioutil.ReadFile(contentZip)
		So(readErr, ShouldBeNil)
		So(string(b), ShouldResemble, "hello world")
	})
}

func TestUnzipContent(t *testing.T) {
	contentZip := "content.zip"

	Convey("should return error if src is empty", t, func() {
		So(UnzipContent("", ""), ShouldEqual, errZipSrcEmpty)
	})

	Convey("should return error if dest is empty", t, func() {
		So(UnzipContent(contentZip, ""), ShouldEqual, errZipDestEmpty)
	})

	Convey("should unzip content", t, func() {
		outputDir := "test/output"

		err := UnzipContent("test/test-content.zip", outputDir)
		So(err, ShouldBeNil)

		defer os.RemoveAll(outputDir)

		assertContentDir(filepath.Join(outputDir, "economy"))
		assertContentFile(filepath.Join(outputDir, "economy", "data.json"), "economy data")

		assertContentDir(filepath.Join(outputDir, "economy", "gdp"))
		assertContentFile(filepath.Join(outputDir, "economy", "gdp", "data.json"), "gdp data")
	})
}

func TestCreateServiceAccount(t *testing.T) {
	Convey("should return error if servicesDir is empty", t, func() {
		serviceToken, err := CreateServiceAccount("")

		So(err, ShouldEqual, errServicesDirEmpty)
		So(serviceToken, ShouldBeEmpty)
	})

	Convey("should return error if servicesDir does not exist", t, func() {
		serviceToken, err := CreateServiceAccount("nodirhere")

		So(err.Error(), ShouldResemble, "servicesDir does not exist nodirhere")
		So(serviceToken, ShouldBeEmpty)
	})

	Convey("should create service account", t, func() {
		serviceToken, err := CreateServiceAccount("test")

		So(err, ShouldBeNil)

		path := fmt.Sprintf("test/%s.json", serviceToken)
		defer os.Remove(path)

		exists, errExists := fileExists(path)
		So(errExists, ShouldBeNil)
		So(exists, ShouldBeTrue)

		b, errRead := ioutil.ReadFile(path)
		So(errRead, ShouldBeNil)

		var actual serviceAccount
		err = json.Unmarshal(b, &actual)
		So(err, ShouldBeNil)
		So(actual, ShouldResemble, defaultServiceAccount)
	})
}

func assertContentDir(path string) {
	exists, existsErr := fileExists(path)
	So(existsErr, ShouldBeNil)
	So(exists, ShouldBeTrue)

	infos, err := os.Stat(path)
	So(err, ShouldBeNil)
	So(infos.IsDir(), ShouldBeTrue)
}

func assertContentFile(path string, expectedContent string) {
	exists, existsErr := fileExists(path)
	So(existsErr, ShouldBeNil)
	So(exists, ShouldBeTrue)

	infos, err := os.Stat(path)
	So(err, ShouldBeNil)
	So(infos.IsDir(), ShouldBeFalse)

	var b []byte
	b, err = ioutil.ReadFile(path)
	So(err, ShouldBeNil)
	So(string(b), ShouldEqual, expectedContent)
}

func tearDownCMSDirs(rootDir string) {
	if exists, err := fileExists(rootDir); exists {
		So(err, ShouldBeNil)

		err := os.RemoveAll(rootDir)
		So(err, ShouldBeNil)
	}
}
