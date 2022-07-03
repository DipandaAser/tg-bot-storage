package manager

import (
	"bytes"
	"container/list"
	"errors"
	"github.com/DipandaAser/tg-bot-storage/pkg/bot"
	v1 "github.com/DipandaAser/tg-bot-storage/pkg/models/v1"
	"io"
	"sync"
	"time"
)

const (
	waitingTimeInSec = 4
)

type Manager struct {
	// bots is a map of bot tokens to bot clients
	bots  map[string]*bot.Client
	mutex sync.Mutex
	// freeBots is a list of free bots. this list contains bot token. that bot token can be use as key for bots map
	freeBots *list.List
	// freeChan is a channel used to make a bot available after and upload. channel of string where string is the bot token
	freeChan chan string
	// requestBotChan is a channel of channel used to request a bot. the channel of string is used to send the bot token
	requestBotChan chan chan string
	// requestList is a list of waiting requests. new requests are added when no bot is available
	requestList *list.List
}

// NewManager creates a new Manager. botsTokens is a slice of distinct bot tokens
func NewManager(botsTokens ...string) (*Manager, error) {

	if len(botsTokens) == 0 {
		return nil, errors.New("no bots tokens provided")
	}

	// check if there are duplicates
	for i := 0; i < len(botsTokens); i++ {
		for j := i + 1; j < len(botsTokens); j++ {
			if botsTokens[i] == botsTokens[j] {
				return nil, errors.New("duplicate bot token")
			}
		}
	}

	// check if all bots are valid, and add them to the manager
	bots := make(map[string]*bot.Client)
	freeBots := list.New()
	for _, token := range botsTokens {
		client, err := bot.NewClient(token)
		if err != nil {
			return nil, err
		}
		bots[token] = client
		freeBots.PushBack(client)
	}

	return &Manager{
		bots:           bots,
		freeChan:       make(chan string),
		freeBots:       freeBots,
		requestList:    list.New(),
		requestBotChan: make(chan chan string),
	}, nil
}

// StartUploaderManager starts the uploader manager
//  This uploader manager will handle requests for getting available bots, and making them available again
func (m *Manager) StartUploaderManager() {
	for {
		select {
		case botToken := <-m.freeChan:
			m.mutex.Lock()

			// we free the bot only if there's no waiting request, otherwise we'll use it for the first waiting request
			if m.requestList.Len() == 0 {
				m.freeBots.PushBack(m.bots[botToken])
				m.mutex.Unlock()
				continue
			}

			// if there's some waiting request we directly use the new free bot
			request := m.requestList.Remove(m.requestList.Front()).(chan string)
			m.mutex.Unlock()
			request <- botToken
		case request := <-m.requestBotChan:
			m.mutex.Lock()

			// if there's no free bot we add the request to the queue of waiting request
			if m.freeBots.Len() == 0 {
				m.requestList.PushBack(request)
				m.mutex.Unlock()
				continue
			}
			freeBotToken := m.freeBots.Remove(m.freeBots.Front()).(*bot.Client).GetToken()
			m.mutex.Unlock()
			request <- freeBotToken
		}
	}
}

func (m *Manager) UploadFileReader(chatId int64, fileName string, fileReader io.Reader) (v1.MessageIdentifier, error) {

	botClient := m.getOneBot()

	// at the end of the function we need to make the bot available again so that it can be used by other requests
	defer time.AfterFunc(waitingTimeInSec*time.Second, func() {
		m.makeBotAvailable(botClient.GetToken())
	})

	msgIdentifier, err := botClient.UploadFileReader(chatId, fileName, fileReader)
	if err != nil {
		return v1.MessageIdentifier{}, err
	}

	return msgIdentifier, nil
}

func (m *Manager) UploadFileBuffer(chatId int64, fileName string, fileData []byte) (v1.MessageIdentifier, error) {
	reader := bytes.NewReader(fileData)
	return m.UploadFileReader(chatId, fileName, reader)
}

func (m *Manager) DownloadFileReader(identifier v1.MessageIdentifier, copyChat int64) (*v1.DownloadReaderResult, error) {
	var botClient *bot.Client
	for _, client := range m.bots {
		botClient = client
		break
	}

	if botClient == nil {
		return nil, errors.New("no bot available")
	}

	return botClient.DownloadFileReader(identifier, copyChat)
}

func (m *Manager) DownloadFileBuffer(identifier v1.MessageIdentifier, copyChat int64) (*v1.DownloadBufferResult, error) {
	result, err := m.DownloadFileReader(identifier, copyChat)
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

// getOneBot returns a bot from the free list and removes it from the free list
func (m *Manager) getOneBot() *bot.Client {
	request := make(chan string)
	m.requestBotChan <- request
	botToken := <-request

	return m.bots[botToken]
}

// makeBotAvailable adds a bot to the free list
func (m *Manager) makeBotAvailable(token string) {
	m.freeChan <- token
}

func (m *Manager) AddBot(botToken string) error {
	if _, ok := m.bots[botToken]; ok {
		return errors.New("bot already exists")
	}

	client, err := bot.NewClient(botToken)
	if err != nil {
		return err
	}

	m.mutex.Lock()
	m.bots[botToken] = client
	m.mutex.Unlock()

	m.makeBotAvailable(botToken)

	return nil
}
