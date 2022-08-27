# Telegram Bot Storage(in development)

![Build Status](https://github.com/dipandaaser/tg-bot-storage/workflows/CI/badge.svg)
[![License](https://img.shields.io/github/license/dipandaaser/tg-bot-storage)](LICENSE)
[![Release](https://img.shields.io/github/release/dipandaaser/tg-bot-storage.svg)](https://github.com/dipandaaser/tg-bot-storage/releases/latest)
[![GitHub Releases Stats of tg-bot-storage](https://img.shields.io/github/downloads/dipandaaser/tg-bot-storage/total.svg?logo=github)](https://somsubhra.github.io/github-release-stats/?username=dipandaaser&repository=tg-bot-storage)

Telegram Bot Storage is a simple library for storing files in Telegram using bot and by passing limits listed
on [telegram bot limits website](https://core.telegram.org/bots/faq#my-bot-is-hitting-limits-how-do-i-avoid-this)

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