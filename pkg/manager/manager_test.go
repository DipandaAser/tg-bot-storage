package manager

import (
	"bytes"
	"fmt"
	config_test "github.com/DipandaAser/tg-bot-storage/config-test"
	"github.com/joho/godotenv"
	"log"
	"sync"
	"testing"
)

const (
	configFilePath = "../../config-test/config.yaml"
)

func init() {
	_ = godotenv.Load("../../.env.test")
}

func Test_UploadFileReaderWithOneBot(t *testing.T) {

	client, err := NewManager(config_test.GetConfig(configFilePath).Tokens[0])
	if err != nil {
		t.Fatal(err)
		return
	}

	go client.StartUploaderManager()

	t.Run("Stress Upload Test", func(t *testing.T) {
		wg := sync.WaitGroup{}
		lock := sync.Mutex{}
		chatId := config_test.GetConfig(configFilePath).ChatID
		count := 0
		total := 5
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
	client, err := NewManager(config_test.GetConfig(configFilePath).Tokens...)
	if err != nil {
		t.Fatal(err)
		return
	}

	go client.StartUploaderManager()

	t.Run("Stress Upload Test", func(t *testing.T) {
		wg := sync.WaitGroup{}
		lock := sync.Mutex{}
		chatId := config_test.GetConfig(configFilePath).ChatID
		count := 0
		total := 5
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
	client, err := NewManager(config_test.GetConfig(configFilePath).Tokens...)
	if err != nil {
		t.Fatal(err)
		return
	}

	go client.StartUploaderManager()

	chatId := config_test.GetConfig(configFilePath).ChatID
	data := bytes.NewReader([]byte("data"))
	fileIdentifier, err := client.UploadFileReader(chatId, t.Name(), data)
	if err != nil {
		t.Fatal(err)
	}

	draftChatId := config_test.GetConfig(configFilePath).DraftChatID

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
		total := 5
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
