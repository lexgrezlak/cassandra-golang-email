package handler

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"net/http/httptest"
	"request-golang/service"
	"testing"
)

func TestGetMessagesByEmail(t *testing.T) {
	testCases := []struct {
		name               string
		want               int
		getMessagesByEmail func(email string, limit int, cursor string) ([]*service.Message, string, error)
	}{
		{
			"valid input",
			http.StatusOK,
			nil,
		},
		{
			"valid input, api returns an error",
			http.StatusInternalServerError,
			func(email string, limit int, cursor string) ([]*service.Message, string, error) {
				return nil, "", fmt.Errorf("failed to get messages")
			}},
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
			got := res.Code
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("mismatch (-want, +got): \n%s", diff)
			}
		})
	}
}
