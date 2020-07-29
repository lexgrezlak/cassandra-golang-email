package service

import (
	"fmt"
	"github.com/gocql/gocql"
)

type Message struct {
	Email string `json:"email"`
	Title string `json:"title"`
	Content string `json:"content"`
	MagicNumber int `json:"magic_number"`
}

func (api *api) GetMessagesByEmail(email string) ([]*Message, error) {
	var messages []*Message
	iterable := api.session.Query(
		`SELECT email, title, content, magic_number FROM message WHERE email=?`,
		email).Consistency(gocql.One).Iter()

	m := map[string]interface{}{}
	for iterable.MapScan(m) {
		message := &Message{
			Email:       m["email"].(string),
			Title:       m["title"].(string),
			Content:     m["content"].(string),
			MagicNumber: m["magic_number"].(int),
		}
		messages = append(messages, message)
		m = map[string]interface{}{}
	}
	return messages, nil
}
func (api *api) DeleteMessage(magicNumber int)  error {
	fmt.Printf("deleting messages with magic number %v", magicNumber)
	return nil
}

func (api *api) CreateMessage(i Message) error {
	if err := api.session.Query(
		`INSERT INTO message (id, email, title, content, magic_number) VALUES (?, ?, ?, ?, ?)`,
		gocql.TimeUUID(),i.Email, i.Title, i.Content, i.MagicNumber).Exec(); err != nil {
		return fmt.Errorf("failed to insert a message: %w", err)
	}
	return nil
}