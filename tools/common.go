package tools

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

var SessionPrefix = "session_"

func GetRandomToken(length int) string {
	r := make([]byte, length)
	io.ReadFull(rand.Reader, r)
	return base64.URLEncoding.EncodeToString(r)
}

func CreateTokenId(sessionId string) string {
	return SessionPrefix + sessionId
}
