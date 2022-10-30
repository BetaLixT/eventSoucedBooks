package book

import (
	"context"
)

type IRepository interface {
	// - Commands
	Create(
		ctx context.Context,
		sagaId *uint64,
		data BookData,
	) error
	Update(
		ctx context.Context,
		sagaId *uint64,
		id uint64,
		ver uint64,
		data BookData,
	) error
	Delete(
		ctx context.Context,
		sagaId *uint64,
		id uint64,
		ver uint64,
	) error
	ListEvents(
		ctx context.Context,
		id uint64,
	) ([]BookEvent, error)
	LastEvent(
		ctx context.Context,
		id uint64,
	) (BookEvent, error)
}
