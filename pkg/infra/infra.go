package infra

import (
	"eventSourcedBooks/pkg/domain/auth"
	"eventSourcedBooks/pkg/domain/base"
	"eventSourcedBooks/pkg/domain/courtroom"
	"eventSourcedBooks/pkg/infra/clients/naga"
	"eventSourcedBooks/pkg/infra/config"
	"eventSourcedBooks/pkg/infra/db"
	"eventSourcedBooks/pkg/infra/http"
	"eventSourcedBooks/pkg/infra/insights"
	"eventSourcedBooks/pkg/infra/logger"
	"eventSourcedBooks/pkg/infra/msgbrker"
	"eventSourcedBooks/pkg/infra/notifdisp"
	"eventSourcedBooks/pkg/infra/repos"

	trace "github.com/BetaLixT/appInsightsTrace"
	"github.com/BetaLixT/gottp"
	"github.com/BetaLixT/rmqevnter"
	"github.com/BetaLixT/tsqlx"
	"github.com/google/wire"
)

var DependencySet = wire.NewSet(
	config.NewConfig,
	config.NewDatabaseOptions,
	config.NewInsightsOptions,
	config.NewNagaOptions,
	config.NewUsagoOptions,

	msgbrker.NewUsagoManager,
	insights.NewInsightsCore,
	logger.NewLoggerFactory,
	wire.Bind(
		new(base.ILoggerFactory),
		new(*logger.LoggerFactory),
	),

	repos.NewAuthRepository,
	wire.Bind(
		new(auth.IRepository),
		new(*repos.AuthRepository),
	),
	repos.NewCourtroomRepository,
	wire.Bind(
		new(courtroom.IRepository),
		new(*repos.CourtroomRepository),
	),

	db.NewDatabaseContext,
	wire.Bind(
		new(tsqlx.ITracer),
		new(*trace.AppInsightsCore),
	),
	http.NewHttpClient,
	wire.Bind(
		new(gottp.ITracer),
		new(*trace.AppInsightsCore),
	),
	naga.NewNagaClient,
	notifdisp.NewNotificationDispatch,
	wire.Bind(
		new(rmqevnter.ITracer),
		new(*trace.AppInsightsCore),
	),
	NewInfrastructure,
)

type Infrastructure struct {
	insightsCore  *trace.AppInsightsCore
	loggerFactory *logger.LoggerFactory
}

func NewInfrastructure(
	insightsCore *trace.AppInsightsCore,
	loggerFactory *logger.LoggerFactory,
) *Infrastructure {
	return &Infrastructure{
		insightsCore:  insightsCore,
		loggerFactory: loggerFactory,
	}
}

// Startup services required by infrastructure
func (infra *Infrastructure) Start() {

}

// Cleaup services required by infrastructure
func (infra *Infrastructure) Stop() {
	infra.insightsCore.Close()
	infra.loggerFactory.Close()
}
