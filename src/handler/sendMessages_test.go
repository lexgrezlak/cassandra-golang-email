package handler

import (
	"bytes"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"net/http/httptest"
	"request-golang/src/config"
	"testing"
)

func TestSendMessages(t *testing.T) {
	validInput := `{"magic_number": 232}`
	invalidInput := `{"magic_number": "hello"}`
	invalidEmptyInput := ``
	testCases := []struct {
		name         string
		// Input is a json string.
		input        string
		want         int
		sendMessages func(magicNumber int, c *config.SmtpConfig) error
	}{
		{
			"valid input",
			validInput,
			http.StatusOK,
			nil,
		},
		{
			"valid input and api returns an error",
			validInput,
			http.StatusInternalServerError,
			func(magicNumber int, c *config.SmtpConfig) error {
				return fmt.Errorf("failed to send messages")
			}},
			{
			"invalid non-empty input",
			invalidInput,
			http.StatusBadRequest,
			nil,
			},
			{
			"invalid empty input",
			invalidEmptyInput,
			http.StatusBadRequest,
			nil,
			},

	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			api := &mockApi{}
			if tc.sendMessages != nil {
				api.MockSendMessages = tc.sendMessages
			}
			res := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/send", bytes.NewBufferString(tc.input))

			h := SendMessages(api, &config.SmtpConfig{})
			h(res, req)
			got := res.Code
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("mismatch (-wantCode, +got): \n%s", diff)
			}
		})
	}
}
