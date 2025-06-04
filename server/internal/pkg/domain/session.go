package domain

import (
	"errors"
	"net/http"
)

var ErrNoSession = errors.New("no session")

type Session struct {
	UserID string
}

type SessionService interface {
	CheckSession(headers http.Header) (Session, error)
}
