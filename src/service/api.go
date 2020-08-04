package service

import (
	"github.com/gocql/gocql"
	"request-golang/src/config"
)

type api struct {
	session *gocql.Session
}

type MessageDatastore interface {
	CreateMessage(i Message) error
	SendMessages(magicNumber int, c *config.SmtpConfig) error
	GetMessagesByEmail(email string, limit int, cursor string) ([]*Message, string, error)
}

func NewAPI(session *gocql.Session) *api {
	return &api{session: session}
}
