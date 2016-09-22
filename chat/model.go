package chat

import "time"

type Conversation struct {
	Messages []*Message `json:"messages"`
}

type Message struct {
	FromUserID int       `json:"from_user_id"`
	ToUserID   int       `json:"to_user_id"`
	Text       string    `json:"text"`
	CreatedAt  time.Time `json:"created_at"`
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
