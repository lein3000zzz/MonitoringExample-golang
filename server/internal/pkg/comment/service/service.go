package service

import (
	"server/internal/pkg/domain"
)

type service struct {
	CommentRepo domain.CommentRepository
	ThreadRepo  domain.ThreadRepository
}

func NewService(commentRepo domain.CommentRepository, threadRepo domain.ThreadRepository) domain.CommentService {
	return service{
		CommentRepo: commentRepo,
		ThreadRepo:  threadRepo,
	}
}

func (s service) Create(threadID string, comment domain.Comment) error {
	if err := s.checkThread(threadID); err != nil {
		return err
	}

	return s.CommentRepo.Create(comment)
}

func (s service) Like(threadID string, commentID string) error {
	if err := s.checkThread(threadID); err != nil {
		return err
	}

	return s.CommentRepo.Like(commentID)
}

func (s service) checkThread(threadID string) error {
	_, err := s.ThreadRepo.Get(threadID)
	return err
}
