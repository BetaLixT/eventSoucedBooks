package base

import (
	"context"

	"go.uber.org/zap"
)

type ILoggerFactory interface {
	NewLogger(
		ctx context.Context,
	) *zap.Logger
}

type ITransaction interface {
	Begin(ctx context.Context) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context)
}

type ITransactionFactory interface {
	Create(ctx context.Context) (ITransaction, error)
}
