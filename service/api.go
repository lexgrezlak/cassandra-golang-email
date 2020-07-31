package service

import (
	"github.com/gocql/gocql"
)

type api struct {
	session *gocql.Session
}

type MessageDatastore interface {
	CreateMessage(i Message) error
	SendMessage(magicNumber int) error
	GetMessagesByEmail(email string, limit int, cursor string) ([]*Message, string, error)
}

func NewAPI(session *gocql.Session) *api {
	return &api{session: session}
}
