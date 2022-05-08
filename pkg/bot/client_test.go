package bot

import (
	"bytes"
	"github.com/joho/godotenv"
	"os"
	"strconv"
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
		return
	}

	t.Run("Upload one file", func(t *testing.T) {
		chatId, _ := strconv.ParseInt(os.Getenv(ENVCHATID), 10, 64)
		data := bytes.NewReader([]byte("data"))
		_, err = client.UploadFileReader(chatId, "test.txt", data)
		if err != nil {
			t.Error(err)
			return
		}
	})
}

func Test_UploadFileBuffer(t *testing.T) {
	client, err := NewClient(os.Getenv(ENVBOTTOKEN))
	if err != nil {
		return
	}

	t.Run("Send one file", func(t *testing.T) {
		chatId, _ := strconv.ParseInt(os.Getenv(ENVCHATID), 10, 64)
		data := []byte("data")
		_, err := client.UploadFileBuffer(chatId, "test.txt", data)
		if err != nil {
			t.Error(err)
			return
		}
	})
}
