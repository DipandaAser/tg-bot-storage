package manager

import (
	"bytes"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
)

const (
	ENVBOTTOKEN    = "BOT_TOKEN"
	ENVBOTTOKENS   = "BOT_TOKENS"
	ENVCHATID      = "CHAT_ID"
	ENVDRAFTCHATID = "DRAFT_CHAT_ID"
)

func init() {
	_ = godotenv.Load("../../.env.test")
}

func Test_UploadFileReaderWithOneBot(t *testing.T) {

	client, err := NewManager(os.Getenv(ENVBOTTOKEN))
	if err != nil {
		t.Fatal(err)
		return
	}

	go client.StartUploaderManager()

	t.Run("Stress Upload Test", func(t *testing.T) {
		wg := sync.WaitGroup{}
		lock := sync.Mutex{}
		chatId, _ := strconv.ParseInt(os.Getenv(ENVCHATID), 10, 64)
		count := 0
		total := 25
		for i := 0; i < total; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				data := bytes.NewReader([]byte("data"))
				_, err = client.UploadFileReader(chatId, fmt.Sprintf("%s %d", t.Name(), i), data)
				if err != nil {
					t.Error(err)
					return
				}
				lock.Lock()
				count++
				log.Printf("Uploaded %d/%d files", count, total)
				lock.Unlock()
			}(i)
		}
		wg.Wait()
	})
}

func Test_UploadFileReaderWithMultipleBot(t *testing.T) {

	botsTokens := strings.Split(os.Getenv(ENVBOTTOKENS), ",")
	client, err := NewManager(botsTokens...)
	if err != nil {
		t.Fatal(err)
		return
	}

	go client.StartUploaderManager()

	t.Run("Stress Upload Test", func(t *testing.T) {
		wg := sync.WaitGroup{}
		lock := sync.Mutex{}
		chatId, _ := strconv.ParseInt(os.Getenv(ENVCHATID), 10, 64)
		count := 0
		total := 60
		for i := 0; i < total; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				data := bytes.NewReader([]byte("data"))
				_, err = client.UploadFileReader(chatId, fmt.Sprintf("%s %d", t.Name(), i), data)
				if err != nil {
					t.Error(err)
					return
				}
				lock.Lock()
				count++
				log.Printf("Uploaded %d/%d files", count, total)
				lock.Unlock()
			}(i)
		}
		wg.Wait()
	})
}

func Test_DownloadFileReader(t *testing.T) {
	client, err := NewManager(strings.Split(os.Getenv(ENVBOTTOKENS), ",")...)
	if err != nil {
		t.Fatal(err)
		return
	}

	go client.StartUploaderManager()

	chatId, _ := strconv.ParseInt(os.Getenv(ENVCHATID), 10, 64)
	data := bytes.NewReader([]byte("data"))
	fileIdentifier, err := client.UploadFileReader(chatId, t.Name(), data)
	if err != nil {
		t.Fatal(err)
	}

	draftChatId, _ := strconv.ParseInt(os.Getenv(ENVDRAFTCHATID), 10, 64)

	t.Run("Download one file", func(t *testing.T) {
		_, err := client.DownloadFileReader(fileIdentifier, draftChatId)
		if err != nil {
			t.Error(err)
			return
		}
	})

	t.Run("Stress Multiple Download", func(t *testing.T) {
		wg := sync.WaitGroup{}
		lock := sync.Mutex{}
		count := 0
		total := 60
		for i := 0; i < total; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				_, err := client.DownloadFileReader(fileIdentifier, draftChatId)
				if err != nil {
					t.Error(err)
					return
				}
				lock.Lock()
				count++
				log.Printf("Downloaded %d/%d files", count, total)
				lock.Unlock()
			}()
		}
		wg.Wait()
	})
}
