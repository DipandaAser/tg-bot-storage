package storage

import "io"

// Storage is the interface that wraps files upload and download
type Storage interface {
	UploadFileBuffer(int64, []byte) (MessageIdentifier, error)
	UploadFileReader(int64, io.Reader) (MessageIdentifier, error)
	DownloadFileBuffer(identifier MessageIdentifier) ([]byte, error)
	DownloadFileReader(identifier MessageIdentifier) (io.ReadCloser, error)
}
