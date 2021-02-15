package zebedee

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ONSdigital/dp-zebedee-content/files"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateDirStructure(t *testing.T) {
	Convey("should return error if root is empty", t, func() {
		err := CreateDirStructure("")
		So(err, ShouldEqual, RootEmptyErr)
	})

	Convey("should create expected directories given a valid root path", t, func() {
		rootDir := "test-content"

		err := CreateDirStructure(rootDir)
		defer tearDown(rootDir)

		So(err, ShouldBeNil)

		for _, dir := range cmsDirs {
			exists, err := files.Exists(filepath.Join(rootDir, dir))
			So(err, ShouldBeNil)
			So(exists, ShouldBeTrue)
		}
	})
}

func tearDown(rootDir string) {
	if exists, _ := files.Exists(rootDir); exists {
		err := os.RemoveAll(rootDir)
		So(err, ShouldBeNil)
	}
}
