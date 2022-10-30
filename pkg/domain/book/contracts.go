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
	Id int64
}

type UpdatePagePosition struct {
	Id       int64
	Position float32
}

type UpdateBookCompleted struct {
	Id        int64
	Completed bool
}

// - Queries
type ListBooks struct {
	CountPerPage int
	PageNumber   int
	OrderBy      *string
	Descending   bool
}
