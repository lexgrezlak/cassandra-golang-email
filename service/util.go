package service

import (
	"encoding/base64"
	"net/smtp"
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

type smtpConfig struct {
	Host    string
	Address string
	Auth    smtp.Auth
	From    string
}

func getSmtpConfig() *smtpConfig {
	// Again, you should never put your API keys, credentials, urls, etc.
	// into your code, even for testing or dev. We're doing it for the purposes
	// of program specification, because it's supposed to work out of the box
	// when you pull it from Docker Hub, so we assume we're not gonna make you configure it.
	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port
	from := "***REMOVED***"
	password := "***REMOVED***"
	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")
	return &smtpConfig{
		Host:    host,
		Address: address,
		Auth:    auth,
		From:    from,
	}
}
