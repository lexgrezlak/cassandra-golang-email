package service

import (
	"github.com/gocql/gocql"
)

type api struct {
	session *gocql.Session
}

type MessageDatastore interface {
	CreateMessage(i Message) error
	DeleteMessage(magicNumber int) error
	GetMessagesByEmail(email string) ([]*Message, error)
}


func NewAPI(session *gocql.Session) *api {
	return &api{session: session}
}