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
	CREATED_AT = "created_at"
)

type Message struct {
	Id          gocql.UUID `json:"id"`
	Email       string `json:"email"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	MagicNumber int    `json:"magic_number"`
	// I think camelCase is better for json,
	// the magic_number is gonna be the exception in this program
	// just to follow the specification
	CreatedAt   time.Time `json:"createdAt"`
}

func (api *api) GetMessagesByEmail(email string, limit int, encodedCursor string) ([]*Message, string, error) {
	var messages []*Message
	var state []byte
	// Defaults to nil in case of an empty cursor.
	state = nil
	if encodedCursor != "" {
		if cursor, err := decodeCursor(encodedCursor); err == nil {
			state = cursor
		}
	}

	q := api.session.Query(
		`SELECT id, email, title, content, magic_number, created_at FROM message WHERE email=?`,
		email).PageState(state)
	if limit > 0 {
		q.PageSize(limit)
	}
	iter := q.Iter()
	endCursor := encodeCursor(iter.PageState())

	m := map[string]interface{}{}
	for iter.MapScan(m) {
		message := &Message{
			Id: 		 m[ID].(gocql.UUID),
			Email:       m[EMAIL].(string),
			Title:       m[TITLE].(string),
			Content:     m[CONTENT].(string),
			MagicNumber: m[MAGIC_NUMBER].(int),
			CreatedAt:   m[CREATED_AT].(time.Time),
		}
		messages = append(messages, message)
		m = map[string]interface{}{}
	}
	return messages, endCursor, nil
}

func (api *api) DeleteMessage(magicNumber int) error {
	// Pull the ids of messages with the specified magicNumber
	// to delete them later.
	iter := api.session.Query(
		`SELECT id FROM message WHERE magic_number=?`, magicNumber).Iter()
	var ids []gocql.UUID
	m := map[string]interface{}{}
	for iter.MapScan(m) {
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
