package service

import (
	"github.com/gocql/gocql"
	"log"
	"net/smtp"
	"request-golang/src/config"
	"strings"
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

type deleteMessageInput struct {
	Id gocql.UUID
	CreatedAt time.Time
	Email string
}

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
	// Magic number can't be 0 (it's a default "undefined" value),
	// so that we can detect if magic number has been specified in the
	// body of a request.
	MagicNumber int        `json:"magic_number" validate:"required"`
}

type Message struct {
	Id          gocql.UUID `json:"id"`
	Email       string `json:"email"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	MagicNumber int `json:"magic_number"`
	CreatedAt time.Time `json:"createdAt"`
}

// Gets the messages from the database and returns them along with an end cursor
// for pagination purposes.
func (api *api) GetMessagesByEmail(i GetMessagesByEmailInput) ([]*Message, string) {
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
	defer iter.Close()

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

	endCursor := encodeCursor(iter.PageState())
	return messages, endCursor
}

// Gets the messages from the database, then sends an email,
// and on success deletes it from the database,
// one by one.
func (api *api) SendMessages(magicNumber int, c *config.SmtpConfig) error {
	// Pull the messages with the specified magicNumber
	// to delete them later.
	iter := api.session.Query(
		`SELECT id, created_at, title, content, email FROM message WHERE magic_number=?`, magicNumber).Iter()

	deleteMessageInputs := make([]*deleteMessageInput, 0)
	sendEmailInputs := make([]*sendEmailInput, 0)
	m := map[string]interface{}{}

	for iter.MapScan(m) {
		// Get the input for delete query
		deleteInput := deleteMessageInput{
			Id:        m[ID].(gocql.UUID),
			CreatedAt: m[CREATED_AT].(time.Time),
			Email: m[EMAIL].(string),
		}
		deleteMessageInputs = append(deleteMessageInputs, &deleteInput)

		// Get the inputs for sending emails
		input := sendEmailInput{
			Email:   m[EMAIL].(string),
			Subject: m[TITLE].(string),
			Content: m[CONTENT].(string),
		}
		sendEmailInputs = append(sendEmailInputs, &input)
		m = map[string]interface{}{}
	}

	// Send an email, and then delete it on each iteration.
	for index, input := range sendEmailInputs {
		err := sendEmail(input, c)
		if err != nil {
			log.Printf("failed to send email: %v",err)
			return err
		} else {
			// If the email has been successfully sent, delete the message
			err := api.deleteMessage(deleteMessageInputs[index])
			if err != nil {
				log.Printf("failed to delete messages: %v", err)
				return err
			}
		}
	}

	return nil
}

func (api *api) CreateMessage(i CreateMessageInput) error {
	id := gocql.TimeUUID()
	createdAt := time.Now()
	if err := api.session.Query(
		`INSERT INTO message (id, email, title, content, magic_number, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
		id, i.Email, i.Title, i.Content, i.MagicNumber, createdAt).Exec(); err != nil {
		log.Printf("failed to insert message: %v", err)
		return err
	}
	return nil
}


// Deletes one message from the database with given id.
func (api *api) deleteMessage(i *deleteMessageInput) error {
	// I'm not sure if that's the best way to delete a message (3 parameters).
	// I would really appreciate your feedback on it.
	if err := api.session.Query(`DELETE FROM message WHERE id=? AND created_at=? AND email=?`, i.Id, i.CreatedAt, i.Email).Exec(); err != nil {
		log.Printf("failed to delete message: %v",err)
		return err
	}
	return nil
}

//
func sendEmail(i *sendEmailInput, c *config.SmtpConfig) error {
	// I'm not sure if I should keep the auth here or keep it higher up
	// in the code (efficiency?, to not have to make a connection every time?).
	// I'd appreciate if you could comment on that.
	hostAndPort := strings.Split(c.Address, ":")
	host := hostAndPort[0]
	auth := smtp.PlainAuth("", c.From, c.Password, host)

	// Email data.
	to := []string{i.Email}
	msg := []byte("To: " + i.Email + "\r\n" +
		"Subject: " + i.Subject + "\r\n" +
		"\r\n" +
		i.Content + "\r\n")

	// Send the mail.
	err := smtp.SendMail(c.Address, auth, c.From, to, msg)
	if err != nil {
		log.Printf("failed to send mail: %v", err)
		return err
	}
	return nil
}
