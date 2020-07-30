package service

import (
	"fmt"
	"github.com/gocql/gocql"
	"time"
)

const (
	ID           = "id"
	EMAIL        = "email"
	TITLE        = "title"
	CONTENT      = "content"
	MAGIC_NUMBER = "magic_number"
)

type Message struct {
	Email       string `json:"email"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	MagicNumber int    `json:"magic_number"`
}

func (api *api) GetMessagesByEmail(email string, limit int, cursor string) ([]*Message, error) {
	var messages []*Message
	iterable := api.session.Query(
		`SELECT email, title, content, magic_number FROM message WHERE email=?`,
		email).Consistency(gocql.One).Iter()

	m := map[string]interface{}{}
	for iterable.MapScan(m) {
		message := &Message{
			Email:       m[EMAIL].(string),
			Title:       m[TITLE].(string),
			Content:     m[CONTENT].(string),
			MagicNumber: m[MAGIC_NUMBER].(int),
		}
		messages = append(messages, message)
		m = map[string]interface{}{}
	}
	return messages, nil
}

func (api *api) DeleteMessage(magicNumber int) error {
	// Pull the ids of messages with the specified magicNumber
	// to delete them later.
	iterable := api.session.Query(
		`SELECT id FROM message WHERE magic_number=?`, magicNumber).Iter()
	var ids []gocql.UUID
	m := map[string]interface{}{}
	for iterable.MapScan(m) {
		id := m[ID].(gocql.UUID)
		ids = append(ids, id)
		m = map[string]interface{}{}
	}

	// Delete the messages with the specified magicNumber.
	if err := api.session.Query(`DELETE FROM message WHERE id IN ?`, ids).Exec(); err != nil {
		return fmt.Errorf("failed to delete messages: %w", err)
	}
	return nil
}

func (api *api) CreateMessage(i Message) error {
	if err := api.session.Query(
		`INSERT INTO message (id, email, title, content, magic_number, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
		gocql.TimeUUID(), i.Email, i.Title, i.Content, i.MagicNumber, time.Now()).Exec(); err != nil {
		return fmt.Errorf("failed to insert a message: %w", err)
	}
	return nil
}
