package storage

// MessageIdentifier is a identifier for a telegram message
type MessageIdentifier struct {
	// ChatId represents the telegram chat id
	ChatId int64
	// MessageId represents the telegram message id
	MessageId int
}
