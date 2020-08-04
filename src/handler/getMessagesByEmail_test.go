package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/google/go-cmp/cmp"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"request-golang/src/service"
	"testing"
	"time"
)

func TestGetMessagesByEmail(t *testing.T) {
	msg1 := service.Message{
		Id:          gocql.TimeUUID(),
		Email:       "john@doe.com",
		Title:       "Title1",
		Content:     "Content1",
		MagicNumber: 32,
		CreatedAt:   time.Now(),
	}
	msg2 := service.Message{
		Id:          gocql.TimeUUID(),
		Email:       "qwerty@gmail.com",
		Title:       "test-title",
		Content:     "test-content",
		MagicNumber: 11,
		CreatedAt:   time.Now(),
	}

	testCases := []struct {
		name               string
		wantCode           int
		getMessagesByEmail func(email string, limit int, cursor string) ([]*service.Message, string, error)
		wantRes            getMessagesByEmailResponse
	}{
		{
			"valid input",
			http.StatusOK,
			nil,
			getMessagesByEmailResponse{
				Messages:  nil,
				EndCursor: "",
			},
		},
		{
			"valid input and api returns an error",
			http.StatusInternalServerError,
			func(email string, limit int, cursor string) ([]*service.Message, string, error) {
				return nil, "", fmt.Errorf("failed to get messages")
			},
			getMessagesByEmailResponse{},
		},
		{
			"valid input and api returns an array of messages with cursor",
			http.StatusOK,
			func(email string, limit int, cursor string) ([]*service.Message, string, error) {
				return []*service.Message{&msg1, &msg2}, "encoded-cursor", nil
			},
			getMessagesByEmailResponse{
				Messages:  []*service.Message{&msg1, &msg2},
				EndCursor: "encoded-cursor",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			api := &mockApi{}
			if tc.getMessagesByEmail != nil {
				api.MockGetMessagesByEmail = tc.getMessagesByEmail
			}
			res := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/messages/john@doe.com", nil)
			h := GetMessagesByEmail(api)
			h(res, req)
			gotCode := res.Code
			if diff := cmp.Diff(tc.wantCode, gotCode); diff != "" {
				t.Errorf("mismatch (-wantCode, +gotCode): \n%s", diff)
			}

			// If the status code isn't 200, finish here, if we didn't skip then we'd get,
			// for example, `Internal Server Error` instead of the object we want.
			// We could as well wrap all the latter part with an if statement
			if gotCode != http.StatusOK {
				t.SkipNow()
			}

			// Only applies to status code 200, check if messages and end cursor are there.
			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Fatal(err)
			}
			expectedData, err := json.Marshal(tc.wantRes)
			if err != nil {
				log.Fatal(err)
			}
			if diff := cmp.Diff(string(expectedData), string(data)); diff != "" {
				t.Errorf("mismatch (-wantCode, +gotCode): \n%s", diff)
			}
		})
	}
}
