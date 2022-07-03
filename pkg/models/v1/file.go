package v1

type UploadFileQuery struct {
	FileName string `json:"file_name"`
	ChatId   int64  `json:"chat_id"`
}

type DownloadFileQuery struct {
	MessageIdentifier
	DraftChatId int64 `json:"draft_chat_id"`
}
