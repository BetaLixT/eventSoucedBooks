package config

import (
	"strings"

	"eventSourcedBooks/pkg/infra/clients/naga"
	"eventSourcedBooks/pkg/infra/db"
	"eventSourcedBooks/pkg/infra/logger"
	"eventSourcedBooks/pkg/infra/msgbrker"

	trace "github.com/BetaLixT/appInsightsTrace"
	"github.com/spf13/viper"
)

func NewInsightsOptions(cfg *viper.Viper) *trace.AppInsightsOptions {
	opt := &trace.AppInsightsOptions{
		ServiceName:        "Pinedule.Courtrooms",
		InstrumentationKey: cfg.GetString("InsightsOptions.InstrumentationKey"),
	}
	if opt.InstrumentationKey == "" {
		panic("Instrumentation key not provided")
	}
	return opt
}

func NewDatabaseOptions(cfg *viper.Viper) *db.DatabaseOptions {
	opt := &db.DatabaseOptions{
		DatabaseServiceName: "main-database",
		ConnectionString:    cfg.GetString("DatabaseOptions.ConnectionString"),
	}
	if opt.ConnectionString == "" {
		panic("ConnectionString key not provided")
	}
	return opt
}

func NewUsagoOptions(cfg *viper.Viper) *msgbrker.UsagoOptions {
	opt := &msgbrker.UsagoOptions{
		Url: cfg.GetString("UsagoOptions.Url"),
	}
	if opt.Url == "" {
		panic("rmq url not provided")
	}
	return opt
}

func NewNagaOptions(cfg *viper.Viper) *naga.NagaOptions {
	opt := &naga.NagaOptions{
		BaseUrl: cfg.GetString("NagaOptions.BaseUrl"),
		ApiKey:  cfg.GetString("NagaOptions.ApiKey"),
	}
	if opt.BaseUrl == "" {
		panic("naga url not provided")
	}
	if opt.BaseUrl == "" {
		panic("naga apikey provided")
	}
	return opt
}

func InitializeConfig(
	logger *logger.LoggerFactory,
	pth string,
) error {
	viper.SetConfigName("config")
	viper.KeyDelimiter("__")
	viper.AddConfigPath(pth)
	viper.SetEnvPrefix("PNDLCR")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.AutomaticEnv()
	viper.BindEnv("PORT", "PORT")

	lgr := logger.NewLogger(nil)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			lgr.Warn("No config file found")
		} else {
			lgr.Error("Failed loading config")
			return err
		}
	}
	return nil
}

func NewConfig(lgr *logger.LoggerFactory) *viper.Viper {
	if err := InitializeConfig(lgr, "./cfg"); err != nil {
		panic(err)
	}

	return viper.GetViper()
}
