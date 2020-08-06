package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/google/go-cmp/cmp"
	"github.com/gorilla/mux"
	"net/http/httptest"
	"request-golang/src/service"
	"testing"
)

func TestGetMessagesByEmail(t *testing.T) {
	email := "john@doe.com"
	jsonMessages := fmt.Sprintf(
		`{"messages":[{"id":"%s","email":"test-email","title":"test-title","content":"test-content","magic_number":12313,"createdAt":"%s"}],"endCursor":"hello"}`,
		gocql.TimeUUID(), "2019-02-02T15:04:05Z")
	var getMessagesByEmailResult struct{
		Messages []*service.Message `json:"messages"`
		EndCursor string `json:"endCursor"`
	}
	err := json.Unmarshal([]byte(jsonMessages), &getMessagesByEmailResult)
	if err != nil {
		t.Fatalf("failed to unmarshal messages: %v", err)
	}

	testCases := []struct{
		name               string
		getMessagesByEmail func(i service.GetMessagesByEmailInput) ([]*service.Message, string)
		wantCode           int
		wantBody           string
	}{
		// API can't return an error, so there's no test for that
		{
			"api returns no messages",
			nil,
			200,
			`{"messages":null,"endCursor":""}`,
		},
		{
			"api returns messages",
			func(i service.GetMessagesByEmailInput) ([]*service.Message, string) {
				return getMessagesByEmailResult.Messages, getMessagesByEmailResult.EndCursor
			},
			200,
			jsonMessages,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			api := &mockApi{}
			if tc.getMessagesByEmail != nil {
				api.MockGetMessagesByEmail = tc.getMessagesByEmail
			}
			res := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/messages/" + email, nil)
			// We have to set url vars for unit testing, otherwise gorilla mux won't register
			// our vars, so the email would be an empty string.
			vars := map[string]string{
				"email": email,
			}
			req = mux.SetURLVars(req, vars)
			h := GetMessagesByEmail(api)
			h(res, req)
			gotCode := res.Code

			// Check for the status code
			if diff := cmp.Diff(tc.wantCode, gotCode); diff != "" {
				t.Errorf("mismatch (-wantCode, +gotCode): \n%s", diff)
			}
			// Check for the response body
			if diff := cmp.Diff(tc.wantBody, res.Body.String()); diff != "" {
				t.Errorf("mismatch (-wantCode, +gotCode): \n%s", diff)
			}

		})
	}
}