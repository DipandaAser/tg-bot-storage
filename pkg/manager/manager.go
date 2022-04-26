package manager

import "github.com/DipandaAser/tg-bot-storage/pkg/bot"

type Manager struct {
	bots []bot.Client
}

func NewManager(botsTokens ...string) *Manager {
	return &Manager{}
}
