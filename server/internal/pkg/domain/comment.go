package domain

type Comment struct {
	ID   string
	Text string
}

type CommentService interface {
	Create(threadID string, comment Comment) error
	Like(threadID string, commentID string) error
}

type CommentRepository interface {
	Create(comment Comment) error
	Like(commentID string) error
}
