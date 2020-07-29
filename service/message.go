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
	fmt.Println("getting all messages")
	messages := []*Message{{Email: "hiell", Title: "title", Content: "asddas", MagicNumber: 123}}
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