package book

// - Commands
type CreateBookCommand struct {
	SagaId      *int
	Title       string
	Description string
	Author      string
	Genres      []string
}

type DeleteBookCommand struct {
	Id uint64
}

type UpdatePagePositionCommand struct {
	Id       int64
	Position float32
}

type UpdateBookCompletedCommand struct {
	Id        int64
	Completed bool
}

type ToggleCompletionCommand struct {
	Id int64
}

// - Queries
type ListBooks struct {
	CountPerPage int
	PageNumber   int
	OrderBy      *string
	Descending   bool
}
