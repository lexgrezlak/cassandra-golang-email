package service

import (
	"encoding/base64"
)

func encodeCursor(cursor []byte) string {
	return base64.StdEncoding.EncodeToString(cursor)
}

func decodeCursor(encodedCursor string) ([]byte, error) {
	decodedCursor, err := base64.StdEncoding.DecodeString(encodedCursor)
	if err != nil {
		return nil, err
	}
	return decodedCursor, nil
}

type FetchParams struct {
	Limit int `validate:"omitempty,min=1,max=100,ne=0"`
	// For some reason tag `base64` doesn't work.
	Cursor string `validate:"omitempty"`
}