package bot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DipandaAser/tg-bot-storage/pkg/storage"
	tb "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"net/http"
)

type Client struct {
	botToken string
	b        *tb.BotAPI
}

type botResult struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
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

func (c *Client) DownloadFileReader(identifier storage.MessageIdentifier, copyChat int64) (io.ReadCloser, error) {
	copiedMessage, err := c.b.Send(tb.ForwardConfig{
		BaseChat:   tb.BaseChat{ChatID: copyChat},
		FromChatID: identifier.ChatId,
		MessageID:  identifier.MessageId,
	})
	if err != nil {
		return nil, err
	}

	if copiedMessage.Document == nil {
		return nil, errors.New("message doesn't contain document")
	}

	url, err := c.b.GetFileDirectURL(copiedMessage.Document.FileID)
	if err != nil {
		return nil, err
	}

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("request fail with status code: %v, %s", response.StatusCode, response.Status)
		// read json error into the botResult struct
		var res botResult
		if json.NewDecoder(response.Body).Decode(&res) != nil {
			return nil, err
		}

		err = fmt.Errorf("code: %v, message, %s", res.ErrorCode, res.Description)
	}

	return response.Body, nil
}

func (c *Client) DownloadFileBuffer(identifier storage.MessageIdentifier, copyChat int64) ([]byte, error) {
	reader, err := c.DownloadFileReader(identifier, copyChat)
	if err != nil {
		return nil, err
	}

	defer reader.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(reader); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
