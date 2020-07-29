package service

type api struct {

}

type MessageDatastore interface {
	CreateMessage() error
	DeleteMessage() error
	GetAllMessagesByEmail(email string) ([]*Message, error)
}
