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
	testCases := []struct {
		name          string
		wantCode      int
		createMessage func(i service.CreateMessageInput) error
	}{
		{"valid input", http.StatusCreated, nil},
		{"valid input, api returns an error", http.StatusInternalServerError, func(i service.CreateMessageInput) error {
			return errors.New("failed to create message")
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			api := &mockApi{}
			if tc.createMessage != nil {
				api.MockCreateMessage = tc.createMessage
			}
			res := httptest.NewRecorder()

			i := service.CreateMessageInput{
				Email:       "hello@world.com",
				Title:       "Hello World",
				Content:     "Content 111.",
				MagicNumber: 324,
			}
			buf := new(bytes.Buffer)
			if err := json.NewEncoder(buf).Encode(i); err != nil {
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
