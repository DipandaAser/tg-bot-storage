package rest_client

import (
	"bytes"
	config_test "github.com/DipandaAser/tg-bot-storage/config-test"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func init() {
	chatId = config_test.GetConfig(configFilePath).ChatID
	draftChatId = config_test.GetConfig(configFilePath).DraftChatID
}

const (
	skipMsg        = "Skipping testing in non CI environment"
	configFilePath = "../../config-test/config.yaml"
)

var (
	apiHost     = os.Getenv("api_host")
	apiKey      = os.Getenv("api_key")
	chatId      int64
	draftChatId int64
)

func ckeckSkip(t *testing.T) {
	if os.Getenv("CI") != "true" {
		t.Skip(skipMsg)
	}
}

func TestRestClient_DownloadFileBuffer(t *testing.T) {
	// check if we are in a CI environnement with the api runinng
	ckeckSkip(t)
	client := NewRestClient(apiHost, apiKey)
	fileContent := []byte("test")
	fileName := "test.txt"
	msgIdentifier, err := client.UploadFileBuffer(chatId, fileName, fileContent)
	if err != nil {
		t.Error(err)
		return
	}

	if !assert.NotEmpty(t, msgIdentifier) {
		return
	}

	downloadResult, err := client.DownloadFileBuffer(msgIdentifier, draftChatId)
	if err != nil {
		t.Error(err)
		return
	}

	if assert.NotEmpty(t, downloadResult) {
		assert.Equal(t, downloadResult.FileInfo.Name, fileName)
		assert.Equal(t, downloadResult.Data, fileContent)
	}
}

func TestRestClient_DownloadFileReader(t *testing.T) {
	// check if we are in a CI environnement with the api runinng
	ckeckSkip(t)
	client := NewRestClient(apiHost, apiKey)
	fileContent := "test"
	fileName := "test.txt"
	msgIdentifier, err := client.UploadFileReader(chatId, fileName, strings.NewReader(fileContent))
	if err != nil {
		t.Error(err)
		return
	}

	if !assert.NotEmpty(t, msgIdentifier) {
		return
	}

	downloadResult, err := client.DownloadFileReader(msgIdentifier, draftChatId)
	if err != nil {
		t.Error(err)
		return
	}

	if !assert.NotEmpty(t, downloadResult) {
		return
	}

	assert.Equal(t, downloadResult.FileInfo.Name, fileName)
	buffer := bytes.Buffer{}
	_, err = buffer.ReadFrom(downloadResult.Data)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, buffer.String(), fileContent)
}

func TestRestClient_UploadFileBuffer(t *testing.T) {
	// check if we are in a CI environnement with the api runinng
	ckeckSkip(t)
	client := NewRestClient(apiHost, apiKey)

	msgIdentifier, err := client.UploadFileBuffer(chatId, "test.txt", []byte("test"))
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotEmpty(t, msgIdentifier)
}

func TestRestClient_UploadFileReader(t *testing.T) {
	// check if we are in a CI environnement with the api runinng
	ckeckSkip(t)
	client := NewRestClient(apiHost, apiKey)

	msgIdentifier, err := client.UploadFileReader(chatId, "test.txt", strings.NewReader("test"))
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotEmpty(t, msgIdentifier)
}
