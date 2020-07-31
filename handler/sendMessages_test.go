package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendMessages(t *testing.T) {
	testCases := []struct {
		name         string
		want         int
		sendMessages func(magicNumber int) error
	}{
		{
			"valid input",
			http.StatusOK,
			nil,
		},
		{
			"valid input and api returns an error",
			http.StatusInternalServerError,
			func(magicNumber int) error {
				return fmt.Errorf("failed to send messages")
			}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			api := &mockApi{}
			if tc.sendMessages != nil {
				api.MockSendMessages = tc.sendMessages
			}
			res := httptest.NewRecorder()

			i := SendMessagesInput{MagicNumber: 390}
			buf := new(bytes.Buffer)
			if err := json.NewEncoder(buf).Encode(i); err != nil {
				log.Fatal(err)
			}
			req := httptest.NewRequest("POST", "/api/send", buf)

			h := SendMessages(api)
			h(res, req)
			got := res.Code
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("mismatch (-want, +got): \n%s", diff)
			}
		})
	}
}
