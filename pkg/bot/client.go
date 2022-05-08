package bot

import (
	"bytes"
	"github.com/DipandaAser/tg-bot-storage/pkg/storage"
	tb "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
)

type Client struct {
	botToken string
	b        *tb.BotAPI
}

func NewClient(botToken string) (*Client, error) {
	myBot, err := tb.NewBotAPI(botToken)
	if err != nil {
		return nil, err
	}

	return &Client{
		botToken: botToken,
		b:        myBot,
	}, nil
}

func (c *Client) UploadFileReader(chatId int64, fileName string, fileReader io.Reader) (storage.MessageIdentifier, error) {

	tbFile := tb.FileReader{
		Name:   fileName,
		Reader: fileReader,
	}

	sentMsg, err := c.b.Send(tb.NewDocument(chatId, tbFile))
	if err != nil {
		return storage.MessageIdentifier{}, err
	}

	return storage.MessageIdentifier{
		ChatId:    chatId,
		MessageId: sentMsg.MessageID,
	}, nil
}

func (c *Client) UploadFileBuffer(chatId int64, fileName string, fileData []byte) (storage.MessageIdentifier, error) {
	reader := bytes.NewReader(fileData)
	return c.UploadFileReader(chatId, fileName, reader)
}
