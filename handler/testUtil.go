package handler

import "request-golang/service"

type mockApi struct {
	MockCreateMessage      func(i service.Message) error
	MockSendMessages       func(magicNumber int) error
	MockGetMessagesByEmail func(email string, limit int, cursor string) ([]*service.Message, string, error)
}

func (api *mockApi) CreateMessage(i service.Message) error {
	if api.MockCreateMessage != nil {
		return api.MockCreateMessage(i)
	}
	return nil
}

func (api *mockApi) GetMessagesByEmail(email string, limit int, cursor string) ([]*service.Message, string, error) {
	if api.MockGetMessagesByEmail != nil {
		return api.MockGetMessagesByEmail(email, limit, cursor)
	}
	return nil, "", nil
}

func (api *mockApi) SendMessages(magicNumber int) error {
	if api.MockSendMessages != nil {
		return api.MockSendMessages(magicNumber)
	}
	return nil
}

