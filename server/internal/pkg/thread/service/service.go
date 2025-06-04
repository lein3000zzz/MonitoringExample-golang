package service

import (
	"server/internal/pkg/domain"
)

type service struct {
	ThreadRepo domain.ThreadRepository
}

func NewService(threadRepo domain.ThreadRepository) domain.ThreadService {
	return service{
		ThreadRepo: threadRepo,
	}
}

func (s service) Create(thread domain.Thread) error {
	return s.ThreadRepo.Create(thread)
}

func (s service) Get(id string) (domain.Thread, error) {
	return s.ThreadRepo.Get(id)
}
