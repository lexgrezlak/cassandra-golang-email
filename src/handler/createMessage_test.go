package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/go-cmp/cmp"
	"log"
	"net/http"
	"net/http/httptest"
	"request-golang/src/service"
	"testing"
)

func TestCreateMessage(t *testing.T) {
	validInput := service.CreateMessageInput{
		Email:       "hello@world.com",
		Title:       "Hello World",
		Content:     "Content 111.",
		MagicNumber: 324,
	}

	invalidInput := service.CreateMessageInput{
		Email:       "hello@world.com",
		MagicNumber: 324,
	}

	testCases := []struct {
		name          string
		wantCode      int
		input 		  service.CreateMessageInput
		createMessage func(i service.CreateMessageInput) error
	}{
		{"valid input", http.StatusCreated, validInput,nil},
		{"valid input, api returns an error", http.StatusInternalServerError,validInput, func(i service.CreateMessageInput) error {
			return errors.New("failed to create message")
		}},
		{"invalid input", http.StatusBadRequest, invalidInput,nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			api := &mockApi{}
			if tc.createMessage != nil {
				api.MockCreateMessage = tc.createMessage
			}
			res := httptest.NewRecorder()


			buf := new(bytes.Buffer)
			if err := json.NewEncoder(buf).Encode(tc.input); err != nil {
				log.Fatal(err)
			}
			req := httptest.NewRequest("POST", "/api/message", buf)
			h := CreateMessage(api)
			h(res, req)

			got := res.Code
			if diff := cmp.Diff(tc.wantCode, got); diff != "" {
				t.Errorf("mismatch (-wantCode, +got): \n%s", diff)
			}
		})
	}
}
