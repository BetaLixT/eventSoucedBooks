package book

import (
	"context"
	"eventSourcedBooks/pkg/domain/base"
)

type IRepository interface {
	// - Commands
	Create(
		ctx context.Context,
		tx base.ITransaction,
		sagaId *uint64,
		data BookData,
	) error
	Update(
		ctx context.Context,
		tx base.ITransaction,
		sagaId *uint64,
		id uint64,
		ver uint64,
		data BookData,
	) error
	Delete(
		ctx context.Context,
		tx base.ITransaction,
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
