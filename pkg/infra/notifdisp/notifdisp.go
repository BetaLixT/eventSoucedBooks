package notifdisp

import (
	"eventSourcedBooks/pkg/infra/logger"
	"eventSourcedBooks/pkg/infra/standard"

	"github.com/BetaLixT/rmqevnter"
	"github.com/BetaLixT/usago"
)

func NewNotificationDispatch(
	lgrf *logger.LoggerFactory,
	mngr *usago.ChannelManager,
	tracer rmqevnter.ITracer,
) *rmqevnter.NotificationDispatch {
	disp := rmqevnter.NewNotifDispatch(
		mngr,
		&rmqevnter.RabbitMQBatchPublisherOptions{
			ExchangeName: standard.NOTIFICATIONS_EX_NAME,
			ExchangeType: standard.NOTIFICATIONS_EX_TYPE,
			ServiceName:  standard.NOTIFICATIONS_SERVICE_NAME,
		},
		lgrf.NewLogger(nil),
		tracer,
	)
	return disp
}
