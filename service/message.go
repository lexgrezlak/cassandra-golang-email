package service

import "fmt"

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
	fmt.Println("creating message")
	return nil
}