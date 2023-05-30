# Telegram Bot Storage(in development)

![Build Status](https://github.com/dipandaaser/tg-bot-storage/workflows/CI/badge.svg)
[![License](https://img.shields.io/github/license/dipandaaser/tg-bot-storage)](LICENSE)
[![Release](https://img.shields.io/github/release/dipandaaser/tg-bot-storage.svg)](https://github.com/dipandaaser/tg-bot-storage/releases/latest)
[![GitHub Releases Stats of tg-bot-storage](https://img.shields.io/github/downloads/dipandaaser/tg-bot-storage/total.svg?logo=github)](https://somsubhra.github.io/github-release-stats/?username=dipandaaser&repository=tg-bot-storage)

Telegram Bot Storage is a simple library for storing files in Telegram using bot and by passing limits listed
on [telegram bot limits website](https://core.telegram.org/bots/faq#my-bot-is-hitting-limits-how-do-i-avoid-this)

### Why we created this project
Telegram bot API has some limits for uploading and downloading files. 
  - With telegram bot API you are limited to 50MB when uploading a file [see here](https://core.telegram.org/bots/faq#how-do-i-upload-a-large-file)
  - With telegram bot API you are limited to 20MB when downloading a file [see here](https://core.telegram.org/bots/faq#how-do-i-download-files)
  - With telegram bot API you are limited to 20 message per minutes in a same chat and 30 message per second in all chats [see here](https://core.telegram.org/bots/faq#my-bot-is-hitting-limits-how-do-i-avoid-this)

All these limitations can be stressful if you plan to used telegram bot API to store files. 
One way to avoid this is to used multiple bots, but even here telegram file storage is designed to be unique for each bot [see here](https://core.telegram.org/bots/api#video); 
a file identifier returning by telegram API to bot A after uploading a file can't be used by bot B to download the file.
The only file identifier that persist across bots is ``file_unique_id`` but we can't use it to download file; its only purpose is to avoid re-uploading a file that has been already uploaded.

The only things left that is persistent across bots is the ``message_id`` containing the uploaded file.

💡 We can use this ``message_id`` field to find the message containing the file and download it. (bye bye file identifier limitations 🚮)

💡 We can use multiple bot to upload files and multiple bot to download files. (bye bye file message per minutes limitations 🚮)


### How it works

  - We used multiple bots to upload and download files.
  - We used a manager to manage the bots.
  - After a bot make an action(upload or download) it will be put in a queue for a specific time 4 seconds by default before being used again. To avoid message / second limitations.
  - After uploading a file, the bot will return the message id of the uploaded file
  - When downloading a file, the bot will make a copy([see here](https://core.telegram.org/bots/api#copymessage)) of the message in order to retrieve the message and get the file


### How to use

#### In your project

```shell
  go get github.com/DipandaAser/tg-bot-storage
```

```go
package main

import (
	"fmt"
	"github.com/DipandaAser/tg-bot-storage/pkg/manager"
	"os"
)

func main() {
	myBotsManager, err := manager.NewManager("MyAwesomeBotToken1", "MyAwesomeBotToken2", "MyAwesomeBotTokenX")
	if err != nil {
		panic(err)
	}

	// Start the uploader manager in a goroutine because it's blocking
	go myBotsManager.StartUploaderManager()

	// Upload a file
	//group or channel or private chat where to upload the file
	chatId := int64(123456789)
	myFile, err := os.Open("path/to/file")
	if err != nil {
		panic(err)
	}
	
	fileIdentifier, err := myBotsManager.UploadFileReader(chatId, "fileName.extension", myFile)
	if err != nil {
		panic(err)
	}

	fmt.Printf("File saved in chat: %v, message identifier: %v\n", fileIdentifier.ChatId, fileIdentifier.MessageId)
	
	// Download a file
	// you can used the same channel as a draft chat
	myDownloadedFile, err := myBotsManager.DownloadFileReader(fileIdentifier, chatId)
	if err != nil {
		panic(err) 
	}
	
	// Do something with the downloaded file
}

```

#### Self hosted service with docker

##### TL;DR

```shell
   docker run --name bot-storage ghcr.io/dipandaaser/bot-storage:latest -v /path/to/config:/app/config -p 80:7000
````

##### Config file

```yaml
tokens:
    - "myBotsToken1" 
    - "myBotsToken2" 
    - "myBotsToken3" 
api_key: "my api key used to access the service"
```

##### Docker Compose

```yaml
version: '3'

services:
  bot-storage:
    image: ghcr.io/dipandaaser/bot-storage:latest
    restart: on-failure
    ports:
      - "7000:7000"
    volumes:
      - "./config:/app/config"
```

### Development

#### Dependencies

- Golang 1.16+
- GNU Make

#### How to run

Create a `config.yaml` file under the [config](config) folder.Use [config.example.yaml](config.example.yaml) as a
template to fill the file.

Run the following command: `make run` or `go run ./cmd/main.go`


#### How to test

Create a `config.yaml` file under the [config-test](config-test)
folder. Use [config.test.example.yaml](config.example.yaml) as a template to fill the file

Run the following commands: `make test` and `make docker-e2etest`