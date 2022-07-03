package v1

import "io"

// MessageIdentifier is a identifier for a telegram message
type MessageIdentifier struct {
	// ChatId represents the telegram chat id
	ChatId int64
	// MessageId represents the telegram message id
	MessageId int
	// FileUniqueId represents the telegram file unique id. This id can be used between bots
	FileUniqueId string
}

type FileInfo struct {
	Size        int64
	Name        string
	ContentType string
}

type DownloadReaderResult struct {
	Data     io.ReadCloser
	FileInfo FileInfo
}

type DownloadBufferResult struct {
	Data     []byte
	FileInfo FileInfo
}
