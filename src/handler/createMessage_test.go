package handler

import (
	"bytes"
	"errors"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"net/http/httptest"
	"request-golang/src/service"
	"testing"
)

func TestCreateMessage(t *testing.T) {
	validInput := `{"email":"hello@world.com","title":"test-title","content":"test-content","magic_number":343}`

	invalidInput := `{"email":"john@he.com","title":"title-2"}`

	testCases := []struct {
		name          string
		wantCode      int
		input 		  string
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
			req := httptest.NewRequest("POST", "/api/message", bytes.NewBufferString(tc.input))
			h := CreateMessage(api)
			h(res, req)

			gotCode := res.Code
			if diff := cmp.Diff(tc.wantCode, gotCode); diff != "" {
				t.Errorf("mismatch (-wantCode, +got): \n%s", diff)
			}
		})
	}
}
