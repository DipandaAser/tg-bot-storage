package storage

import "io"

// Storage is the interface that wraps files upload and download
type Storage interface {
	UploadFileBuffer(int64, string, []byte) (MessageIdentifier, error)
	UploadFileReader(int64, string, io.Reader) (MessageIdentifier, error)
	DownloadFileBuffer(identifier MessageIdentifier) ([]byte, error)
	DownloadFileReader(identifier MessageIdentifier) (io.ReadCloser, error)
}
