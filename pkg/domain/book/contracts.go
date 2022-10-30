package book

// - Commands
type CreateBookCommand struct {
	SagaId      *uint64
	Title       string
	Description string
	Author      string
	Genres      []string
}

type DeleteBookCommand struct {
	SagaId *uint64
	Id     uint64
}

type UpdatePagePositionCommand struct {
	SagaId   *uint64
	Id       uint64
	Position float32
}

type UpdateBookCompletedCommand struct {
	SagaId    *uint64
	Id        uint64
	Completed bool
}

type ToggleCompletionCommand struct {
	Id     uint64
	SagaId *uint64
}

// - Queries
type ListBooks struct {
	CountPerPage int
	PageNumber   int
	OrderBy      *string
	Descending   bool
}
