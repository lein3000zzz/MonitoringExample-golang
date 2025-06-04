package session

import (
	"errors"
	"net/http"
	"server/internal/pkg/domain"
)

type service struct{}

func NewService() domain.SessionService {
	return service{}
}

func (s service) CheckSession(headers http.Header) (domain.Session, error) {
	req, err := http.NewRequest(http.MethodGet, "http://vk-golang.ru:17000/int/CheckSession", nil)
	if err != nil {
		return domain.Session{}, err
	}

	req.Header = headers

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return domain.Session{}, err
	}

	switch resp.StatusCode {
	case 500:
		return domain.Session{}, errors.New("failed to request check session")
	case 200:
		return domain.Session{}, nil
	default:
		return domain.Session{}, domain.ErrNoSession
	}
}
