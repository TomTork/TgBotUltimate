package Database

import "time"

type ChatMessage struct {
	TgId    uint64
	Message string
}

type Message struct {
	Id        uint64
	CreatedAt time.Time
	ChatMessage
}
