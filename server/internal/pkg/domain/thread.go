package domain

type Thread struct {
	ID   string
	Name string
}

type ThreadService interface {
	Create(thread Thread) error
	Get(id string) (Thread, error)
}

type ThreadRepository interface {
	Create(thread Thread) error
	Get(id string) (Thread, error)
}
