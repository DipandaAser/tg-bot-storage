# Telegram Bot Storage(in development)
![Build Status](https://github.com/dipandaaser/tg-bot-storage/workflows/CI/badge.svg)
[![License](https://img.shields.io/github/license/dipandaaser/tg-bot-storage)](LICENSE)
[![Release](https://img.shields.io/github/release/dipandaaser/tg-bot-storage.svg)](https://github.com/dipandaaser/tg-bot-storage/releases/latest)
[![GitHub Releases Stats of tg-bot-storage](https://img.shields.io/github/downloads/dipandaaser/tg-bot-storage/total.svg?logo=github)](https://somsubhra.github.io/github-release-stats/?username=dipandaaser&repository=tg-bot-storage)

Telegram Bot Storage is a simple library for storing files in Telegram using bot and by passing limits listed on [telegram bot limits website](https://core.telegram.org/bots/faq#my-bot-is-hitting-limits-how-do-i-avoid-this)

### Development

#### Dependencies
- Golang 1.16+
- GNU Make

#### How to test
Create a `.env.test` file with [.env.test.example](.env.test.example) as a template OR just set env var with the same name.
Run the following command: `make test`