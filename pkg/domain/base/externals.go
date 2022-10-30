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
