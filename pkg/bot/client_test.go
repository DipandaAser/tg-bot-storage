package bot

import (
	"bytes"
	"github.com/joho/godotenv"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"
)

const (
	ENVBOTTOKEN = "BOT_TOKEN"
	ENVCHATID   = "CHAT_ID"
)

func init() {
	_ = godotenv.Load("../../.env.test")
}

func Test_UploadFileReader(t *testing.T) {
	client, err := NewClient(os.Getenv(ENVBOTTOKEN))
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Run("Upload one file", func(t *testing.T) {
		chatId, _ := strconv.ParseInt(os.Getenv(ENVCHATID), 10, 64)
		data := bytes.NewReader([]byte("data"))
		_, err = client.UploadFileReader(chatId, t.Name(), data)
		if err != nil {
			t.Error(err)
			return
		}
	})
}

func Test_UploadFileBuffer(t *testing.T) {
	client, err := NewClient(os.Getenv(ENVBOTTOKEN))
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Run("Send one file", func(t *testing.T) {
		chatId, _ := strconv.ParseInt(os.Getenv(ENVCHATID), 10, 64)
		data := []byte("data")
		_, err := client.UploadFileBuffer(chatId, t.Name(), data)
		if err != nil {
			t.Error(err)
			return
		}
	})
}

func Test_DownloadFileReader(t *testing.T) {
	client, err := NewClient(os.Getenv(ENVBOTTOKEN))
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Run("Download one file", func(t *testing.T) {
		chatId, _ := strconv.ParseInt(os.Getenv(ENVCHATID), 10, 64)
		fileContent := "data"
		msgIdentifier, err := client.UploadFileReader(chatId, t.Name(), strings.NewReader(fileContent))
		if err != nil {
			t.Error(err)
			return
		}

		data, err := client.DownloadFileReader(msgIdentifier, chatId)
		if err != nil {
			t.Error(err)
			return
		}

		contentDownloaded, err := ioutil.ReadAll(data)
		if err != nil {
			t.Error(err)
			return
		}

		if string(contentDownloaded) != fileContent {
			t.Errorf("Expected %s, got %s", fileContent, string(contentDownloaded))
			return
		}
	})
}
