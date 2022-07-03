package storage

import (
	v1 "github.com/DipandaAser/tg-bot-storage/pkg/models/v1"
	"io"
)

// Storage is the interface that wraps files upload and download
type Storage interface {
	UploadFileBuffer(int64, string, []byte) (v1.MessageIdentifier, error)
	UploadFileReader(int64, string, io.Reader) (v1.MessageIdentifier, error)
	DownloadFileBuffer(identifier v1.MessageIdentifier, copyChat int64) (*v1.DownloadBufferResult, error)
	DownloadFileReader(identifier v1.MessageIdentifier, copyChat int64) (*v1.DownloadReaderResult, error)
}
