package mocks

import (
	"io"
	"sync"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// DownloaderMock mock implementation of a FileDownloader.
type DownloaderMock struct {
	mutex    sync.Mutex
	calls    []*s3.GetObjectInput
	stubFunc func() (int64, error)
}

// ErroringDownloader construct a new DownloaderMock that will return the specified error which invoked.
func ErroringDownloader(errToReturn error) *DownloaderMock {
	return &DownloaderMock{
		calls: make([]*s3.GetObjectInput, 0),
		stubFunc: func() (int64, error) {
			return 0, errToReturn
		},
	}
}

// SuccessfulDownloader construct a new DownloaderMock that will invoke the provided success func when called.
func SuccessfulDownloader(successFunc func() (int64, error)) *DownloaderMock {
	return &DownloaderMock{
		calls:    make([]*s3.GetObjectInput, 0),
		stubFunc: successFunc,
	}
}

// Download the mocked out implementation of the Download method.
func (d *DownloaderMock) Download(w io.WriterAt, input *s3.GetObjectInput, options ...func(downloader *s3manager.Downloader)) (n int64, err error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.calls = append(d.calls, input)

	return d.stubFunc()
}

// GetCalls return an array of *s3.GetObjectInput passed as a parameter to the invocation of the Download method.
func (d *DownloaderMock) GetCalls() []*s3.GetObjectInput {
	return d.calls
}
