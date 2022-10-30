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
		id uint64,
		ver uint64,
		data BookData,
	)
	DeleteBook(
		ctx context.Context,
		id uint64,
		ver uint64,
	) error
	ListEvents(
		ctx context.Context,
		id string,
	) ([]BookEvent, error)
	LastEvent(
		ctx context.Context,
		id uint64,
	) (BookEvent, error)
}
