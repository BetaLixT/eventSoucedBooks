package book

import (
	"context"
)

type IRepository interface {
	// - Commands
	CreateBook(
		ctx context.Context,
		data BookData,
	) error
	UpdateBook(
		ctx context.Context,
		id string,
		data BookData,
	)
	DeleteBook(
		ctx context.Context,
		id string,
	)
	ListEvents(
		ctx context.Context,
		id string,
	) ([]BookEvent, error)
}
