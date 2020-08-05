package service

import (
	"fmt"
	"github.com/gocql/gocql"
	"net/smtp"
	"request-golang/src/config"
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
	Email   string
	Subject string
	Content string
}

type GetMessagesByEmailInput struct {
	Email string `validate:"email"`
	Params FetchParams
}

type CreateMessageInput struct {
	Email       string     `json:"email" validate:"email"`
	Title       string     `json:"title" validate:"required,max=200"`
	Content     string     `json:"content" validate:"required,max=5000"`
	// We assume camelCase is better for json,
	// the magic_number is gonna be the exception in this program
	// just to follow the specification.
	MagicNumber int        `json:"magic_number" validate:"required"`
}

type Message struct {
	Id          gocql.UUID
	Email       string
	Title       string
	Content     string
	MagicNumber int
	CreatedAt time.Time
}

func (api *api) GetMessagesByEmail(i GetMessagesByEmailInput) ([]*Message, string, error) {
	var messages []*Message
	var state []byte
	state = nil
	// Defaults to nil in case of an empty cursor.
	if i.Params.Cursor != "" {
		if cursor, err := decodeCursor(i.Params.Cursor); err == nil {
			state = cursor
		}
	}
	q := api.session.Query(
		`SELECT id, email, title, content, magic_number, created_at FROM message WHERE email=?`,
		i.Email).PageState(state)
	if i.Params.Limit > 0 {
		q.PageSize(i.Params.Limit)
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

func (api *api) SendMessages(magicNumber int, c *config.SmtpConfig) error {
	// Pull the ids of messages with the specified magicNumber
	// to delete them later.
	iter := api.session.Query(
		`SELECT id, title, content, email FROM message WHERE magic_number=?`, magicNumber).Iter()

	ids := make([]*gocql.UUID, 0)
	sendEmailInputs := make([]*sendEmailInput, 0)
	m := map[string]interface{}{}

	for iter.MapScan(m) {
		// Get the ids for delete query
		id := m[ID].(gocql.UUID)
		ids = append(ids, &id)

		// Get the inputs for sending emails
		input := sendEmailInput{
			Email:   m[EMAIL].(string),
			Subject: m[TITLE].(string),
			Content: m[CONTENT].(string),
		}
		sendEmailInputs = append(sendEmailInputs, &input)
		fmt.Printf("INPUT %v INPUT", sendEmailInputs)

		m = map[string]interface{}{}
	}

	// Send an email, and then delete it on each iteration.
	for i, input := range sendEmailInputs {
		err := sendEmail(input, c)
		if err != nil {
			return fmt.Errorf("failed to send an email: %v", err)
			// If the email has been successfully sent, delete the message
		} else {
			err := api.deleteMessage(ids[i])
			if err != nil {
				return fmt.Errorf("failed to delete messages: %v", err)
			}
		}
	}

	return nil
}

//
func (api *api) CreateMessage(i CreateMessageInput) error {
	id := gocql.TimeUUID()
	createdAt := time.Now()
	if err := api.session.Query(
		`INSERT INTO message (id, email, title, content, magic_number, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
		id, i.Email, i.Title, i.Content, i.MagicNumber, createdAt).Exec(); err != nil {
		return fmt.Errorf("failed to insert a message: %w", err)
	}
	return nil
}

// Deletes one message from the database with given id.
func (api *api) deleteMessage(id *gocql.UUID) error {
	if err := api.session.Query(`DELETE FROM message WHERE id=?`, id).Exec(); err != nil {
		return err
	}
	return nil
}

func sendEmail(i *sendEmailInput, c *config.SmtpConfig) error {
	// I'm not sure if I should keep the auth here or keep it higher up
	// in the code (efficiency?, to not have to make a connection every time?).
	// I'd appreciate if you could comment on that.
	auth := smtp.PlainAuth("", c.From, c.Password, "smtp.gmail.com")

	// Email data.
	to := []string{i.Email}
	msg := []byte("To: " + i.Email + "\r\n" +
		"Subject: " + i.Subject + "\r\n" +
		"\r\n" +
		i.Content + "\r\n")

	err := smtp.SendMail(c.Address, auth, c.From, to, msg)
	if err != nil {
		return err
	}
	return nil
}
