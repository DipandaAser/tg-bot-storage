package bot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	v1 "github.com/DipandaAser/tg-bot-storage/pkg/models/v1"
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

func (c *Client) UploadFileReader(chatId int64, fileName string, fileReader io.Reader) (v1.MessageIdentifier, error) {

	tbFile := tb.FileReader{
		Name:   fileName,
		Reader: fileReader,
	}

	sentMsg, err := c.b.Send(tb.NewDocument(chatId, tbFile))
	if err != nil {
		return v1.MessageIdentifier{}, err
	}

	return v1.MessageIdentifier{
		ChatId:       chatId,
		MessageId:    sentMsg.MessageID,
		FileUniqueId: sentMsg.Document.FileUniqueID,
	}, nil
}

func (c *Client) UploadFileBuffer(chatId int64, fileName string, fileData []byte) (v1.MessageIdentifier, error) {
	reader := bytes.NewReader(fileData)
	return c.UploadFileReader(chatId, fileName, reader)
}

func (c *Client) DownloadFileReader(identifier v1.MessageIdentifier, copyChat int64) (*v1.DownloadReaderResult, error) {
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
		return nil, err
	}

	return &v1.DownloadReaderResult{
		Data: response.Body,
		FileInfo: v1.FileInfo{
			Size:        int64(copiedMessage.Document.FileSize),
			Name:        copiedMessage.Document.FileName,
			ContentType: copiedMessage.Document.MimeType,
		},
	}, nil
}

func (c *Client) DownloadFileBuffer(identifier v1.MessageIdentifier, copyChat int64) (*v1.DownloadBufferResult, error) {
	result, err := c.DownloadFileReader(identifier, copyChat)
	if err != nil {
		return nil, err
	}

	defer result.Data.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(result.Data); err != nil {
		return nil, err
	}

	return &v1.DownloadBufferResult{
		Data:     buf.Bytes(),
		FileInfo: result.FileInfo,
	}, nil
}

func (c *Client) GetToken() string {
	return c.botToken
}
