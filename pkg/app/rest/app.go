package rest

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	v1 "eventSourcedBooks/pkg/app/rest/controllers/v1"
	"eventSourcedBooks/pkg/app/rest/middlewares"
	"eventSourcedBooks/pkg/domain"
	"eventSourcedBooks/pkg/domain/auth"
	"eventSourcedBooks/pkg/infra/common"
	"eventSourcedBooks/pkg/infra"
	"eventSourcedBooks/pkg/infra/db"
	"eventSourcedBooks/pkg/infra/logger"

	trace "github.com/BetaLixT/appInsightsTrace"
	"github.com/BetaLixT/tsqlx"
	"github.com/betalixt/gingorr"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/soreing/trex"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

var dependencySet = wire.NewSet(
	domain.DependencySet,
	infra.DependencySet,
	v1.NewCourtroomController,
	NewApp,
)

var authDisabled = false

func Start(authDisb bool) {
	authDisabled = authDisb

	app, err := InitializeApp()
	if err != nil {
		panic(err)
	}
	fmt.Println(app)
	app.startService()
}

type app struct {
	lgrf    *logger.LoggerFactory
	authsvc *auth.AuthService
	v1croom *v1.CourtroomController

	insights *trace.AppInsightsCore
	inf      *infra.Infrastructure
	dbctx    *tsqlx.TracedDB
}

func NewApp(
	lgrf *logger.LoggerFactory,
	authsvc *auth.AuthService,
	v1croom *v1.CourtroomController,
	insights *trace.AppInsightsCore,
	inf *infra.Infrastructure,
	dbctx *tsqlx.TracedDB,
) *app {
	return &app{
		lgrf:     lgrf,
		authsvc:  authsvc,
		v1croom:  v1croom,
		insights: insights,
		inf:      inf,
		dbctx:    dbctx,
	}
}

func (a *app) startService() {

	// - Empty Context
	ctx := context.TODO()

	// - Setting up logger
	baseLgr := a.lgrf.NewLogger(ctx)

	// - Setting up gin router
	router := gin.New()
	// gin.SetMode(gin.ReleaseMode)
	router.SetTrustedProxies(nil)

	// - Running migrations
	db.RunMigrations(
		ctx,
		baseLgr,
		a.dbctx,
		db.GetMigrationScripts(),
	)

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status": "alive",
		})
	})
	// - Setting up middlewares
	router.Use(gingorr.RootRecoveryMiddleware(baseLgr))
	router.Use(trex.TxContextMiddleware(common.TRACE_INFO_KEY))
	router.Use(trex.RequestTracerMiddleware(a.traceRequest))
	router.Use(gingorr.RecoveryMiddleware(a.lgrf, common.TRACE_INFO_KEY))
	router.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)
	router.Use(gingorr.ErrorHandlerMiddleware(a.lgrf, common.TRACE_INFO_KEY))
	if !authDisabled {
		router.Use(middlewares.AuthMiddleware(a.authsvc))
	}

	// - Setting up routes

	v1g := router.Group("api/v1")
	a.v1croom.RegisterRoutes(v1g.Group("courtrooms"))

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(404, "Not Found")
	})

	a.inf.Start()
	go func() {
		if err := router.Run(":8080"); err != nil && err != http.ErrServerClosed {
			baseLgr.Error(
				"failed running router",
				zap.Error(err),
			)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	a.inf.Stop()
	baseLgr.Info("Server exiting")
}

func (a *app) traceRequest(
	context context.Context,
	method,
	path,
	query,
	agent,
	ip string,
	status,
	bytes int,
	start,
	end time.Time) {
	latency := end.Sub(start)

	lgr := a.lgrf.NewLogger(context)
	a.insights.TraceRequest(
		context,
		method,
		path,
		query,
		status,
		bytes,
		ip,
		agent,
		start,
		end,
		map[string]string{},
	)
	lgr.Info(
		"Request",
		zap.Int("status", status),
		zap.String("method", method),
		zap.String("path", path),
		zap.String("query", query),
		zap.String("ip", ip),
		zap.String("userAgent", agent),
		zap.Time("mvts", end),
		zap.String("pmvts", end.Format("2006-01-02T15:04:05-0700")),
		zap.Duration("latency", latency),
		zap.String("pLatency", latency.String()),
	)
}
