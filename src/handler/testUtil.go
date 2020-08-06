package handler

import (
	"request-golang/src/config"
	"request-golang/src/service"
)

type mockApi struct {
	MockCreateMessage      func(i service.CreateMessageInput) error
	MockSendMessages       func(magicNumber int, c *config.SmtpConfig) error
	MockGetMessagesByEmail func(i service.GetMessagesByEmailInput) ([]*service.Message, string)
}

func (api *mockApi) CreateMessage(i service.CreateMessageInput) error {
	if api.MockCreateMessage != nil {
		return api.MockCreateMessage(i)
	}
	return nil
}

func (api *mockApi) GetMessagesByEmail(i service.GetMessagesByEmailInput) ([]*service.Message, string) {
	if api.MockGetMessagesByEmail != nil {
		return api.MockGetMessagesByEmail(i)
	}
	return nil, "test-cursor"
}

func (api *mockApi) SendMessages(magicNumber int, c *config.SmtpConfig) error {
	if api.MockSendMessages != nil {
		return api.MockSendMessages(magicNumber, c)
	}
	return nil
}
