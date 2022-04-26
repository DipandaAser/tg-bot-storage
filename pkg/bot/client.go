package bot

import (
	"time"
)

type Client struct {
	botToken            string
	lastMessageSendTime time.Time
}

func NewClient(botToken string) *Client {
	return &Client{
		botToken: botToken,
	}
}

func (c *Client) GetLastMessageSendTime() time.Time {
	return c.lastMessageSendTime
}
