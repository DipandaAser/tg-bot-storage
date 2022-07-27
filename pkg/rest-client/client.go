package rest_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	v1 "github.com/DipandaAser/tg-bot-storage/pkg/models/v1"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type RestClient struct {
	apiKey string
	apiUrl string
	client http.Client
}

type apiError struct {
	Error string `json:"error"`
}

func NewRestClient(apiUrl, apiKey string) RestClient {
	return RestClient{
		apiKey: apiKey,
		apiUrl: strings.TrimSuffix(apiUrl, "/"),
		client: http.Client{},
	}
}

func (rc *RestClient) getApiUrl() string {
	return rc.apiUrl + "/api"
}

func (rc *RestClient) UploadFileReader(chatId int64, fileName string, fileReader io.Reader) (v1.MessageIdentifier, error) {
	params := url.Values{}
	params.Set("chat_id", fmt.Sprintf("%d", chatId))
	params.Set("file_name", fileName)
	params.Set("api-key", rc.apiKey)
	var finalUrl = fmt.Sprintf("%s/%s?%s", rc.getApiUrl(), "files", params.Encode())
	response, err := rc.client.Post(finalUrl, "application/octet-stream", fileReader)
	if err != nil {
		return v1.MessageIdentifier{}, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		var apiErr apiError
		if err := json.NewDecoder(response.Body).Decode(&apiErr); err != nil {
			return v1.MessageIdentifier{}, fmt.Errorf("could not upload file. status code: %d", response.StatusCode)
		}
		return v1.MessageIdentifier{}, fmt.Errorf("could not upload file. status code: %d, error: %s", response.StatusCode, apiErr.Error)
	}
	var fileIdentifier v1.MessageIdentifier
	err = json.NewDecoder(response.Body).Decode(&fileIdentifier)
	if err != nil {
		return v1.MessageIdentifier{}, err
	}
	return fileIdentifier, nil
}

func (rc *RestClient) UploadFileBuffer(chatId int64, fileName string, fileData []byte) (v1.MessageIdentifier, error) {
	reader := bytes.NewReader(fileData)
	return rc.UploadFileReader(chatId, fileName, reader)
}

func (rc *RestClient) DownloadFileReader(identifier v1.MessageIdentifier, copyChat int64) (*v1.DownloadReaderResult, error) {
	params := url.Values{}
	params.Set("chat_id", fmt.Sprintf("%d", identifier.ChatId))
	params.Set("msg_id", fmt.Sprintf("%d", identifier.MessageId))
	params.Set("draft_chat_id", fmt.Sprintf("%d", copyChat))
	params.Set("api-key", rc.apiKey)
	var finalUrl = fmt.Sprintf("%s/%s?%s", rc.getApiUrl(), "files", params.Encode())
	response, err := rc.client.Get(finalUrl)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		var apiErr apiError
		if err := json.NewDecoder(response.Body).Decode(&apiErr); err != nil {
			return nil, fmt.Errorf("could not upload file. status code: %d", response.StatusCode)
		}
		return nil, fmt.Errorf("could not upload file. status code: %d, error: %s", response.StatusCode, apiErr.Error)
	}
	// get the file size in header content-length and convert it to int64
	var fileSize int64 = 0
	if contentLength := response.Header.Get("Content-Length"); contentLength != "" {
		fileSize, _ = strconv.ParseInt(contentLength, 10, 64)
	}

	_, contentDispParams, err := mime.ParseMediaType(response.Header.Get("Content-Disposition"))
	return &v1.DownloadReaderResult{
		Data: response.Body,
		FileInfo: v1.FileInfo{
			Size:        fileSize,
			Name:        contentDispParams["filename"],
			ContentType: response.Header.Get("Content-Type"),
		},
	}, nil
}

func (rc *RestClient) DownloadFileBuffer(identifier v1.MessageIdentifier, copyChat int64) (*v1.DownloadBufferResult, error) {
	result, err := rc.DownloadFileReader(identifier, copyChat)
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
