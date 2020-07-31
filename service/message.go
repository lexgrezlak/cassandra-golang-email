package service

import (
	"fmt"
	"github.com/gocql/gocql"
	"net/smtp"
	"time"
)

const (
	ID           = "id"
	EMAIL        = "email"
	TITLE        = "title"
	CONTENT      = "content"
	MAGIC_NUMBER = "magic_number"
	CREATED_AT   = "created_at"
)

type sendEmailInput struct {
	Email string
	Subject string
	Content string
}

type Message struct {
	Id          gocql.UUID `json:"id"`
	Email       string     `json:"email"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	MagicNumber int        `json:"magic_number"`
	// I think camelCase is better for json,
	// the magic_number is gonna be the exception in this program
	// just to follow the specification
	CreatedAt time.Time `json:"createdAt"`
}

func (api *api) GetMessagesByEmail(email string, limit int, encodedCursor string) ([]*Message, string, error) {
	var messages []*Message
	var state []byte
	state = nil
	// Defaults to nil in case of an empty cursor.
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
			Id:          m[ID].(gocql.UUID),
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



func (api *api) SendMessages(magicNumber int) error {
	// Pull the ids of messages with the specified magicNumber
	// to delete them later.
	iter := api.session.Query(
		`SELECT id, title, content, email FROM message WHERE magic_number=?`, magicNumber).Iter()

	// We're allocating the length of array so that there's no
	// "indexing may cause panic because of nil slice"
	// when using the index of `sendEmailInputs` to access `ids` in the for loop below
	ids := make([]*gocql.UUID, iter.NumRows())
	sendEmailInputs := make([]sendEmailInput, iter.NumRows())

	m := map[string]interface{}{}

	for iter.MapScan(m) {
		// Get the ids for delete query
		id := m[ID].(gocql.UUID)
		ids = append(ids, &id)

		// Get the inputs for sending emails
		input := sendEmailInput{
			Email:       m[EMAIL].(string),
			Subject:       m[TITLE].(string),
			Content:     m[CONTENT].(string),
		}
		sendEmailInputs = append(sendEmailInputs, input)

		m = map[string]interface{}{}
	}

	// Send an email, and then delete it on each iteration.
	for i, input := range sendEmailInputs {
		err := sendEmail(input)
		if err != nil {
			return fmt.Errorf("failed to send an email: %v", err)
		// If email has been successfully sent, delete the message
		} else {
			err := api.deleteMessage(ids[i])
			if err != nil {
				return fmt.Errorf("failed to delete messages: %v", err)
			}
		}
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

func (api *api) deleteMessage(id *gocql.UUID) error {
	if err := api.session.Query(`DELETE FROM message WHERE id=?`, id).Exec(); err != nil {
		return err
	}
	return nil
}

func sendEmail(i sendEmailInput) error {
	to := []string{i.Email}
	msg := []byte("To: " + i.Email + "\r\n" +
		"Subject: " + i.Subject + "\r\n" +
		"\r\n" +
		i.Content + "\r\n")
	smtpConfig := getSmtpConfig()
	err := smtp.SendMail(smtpConfig.Address, smtpConfig.Auth, smtpConfig.From, to, msg)
	if err != nil {
		return err
	}
	return nil
}

