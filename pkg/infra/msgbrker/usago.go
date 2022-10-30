package msgbrker

import (
	"context"

	"eventSourcedBooks/pkg/infra/logger"

	"github.com/BetaLixT/usago"
)

func NewUsagoManager(
	lgrf *logger.LoggerFactory,
	optn *UsagoOptions,
) *usago.ChannelManager {
	return usago.NewChannelManager(
		optn.Url,
		lgrf.NewLogger(context.TODO()),
	)
}
