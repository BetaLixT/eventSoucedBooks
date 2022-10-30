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
