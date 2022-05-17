package storage

import "io"

// Storage is the interface that wraps files upload and download
type Storage interface {
	UploadFileBuffer(int64, string, []byte) (MessageIdentifier, error)
	UploadFileReader(int64, string, io.Reader) (MessageIdentifier, error)
	DownloadFileBuffer(identifier MessageIdentifier, copyChat int64) ([]byte, error)
	DownloadFileReader(identifier MessageIdentifier, copyChat int64) (io.ReadCloser, error)
}
