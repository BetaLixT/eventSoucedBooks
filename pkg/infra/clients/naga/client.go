package naga

import (
	"context"
	"fmt"

	"eventSourcedBooks/pkg/domain/base"
	"eventSourcedBooks/pkg/domain/courtroom"
	"eventSourcedBooks/pkg/infra/logger"

	"github.com/BetaLixT/gottp"
	"go.uber.org/zap"
)

type NagaClient struct {
	optn *NagaOptions
	gttp *gottp.HttpClient
	lgrf *logger.LoggerFactory
}

var _ courtroom.IGraphClient = (*NagaClient)(nil)

func (repo *NagaClient) CreateMeeting(
	ctx context.Context,
) (base.JoinInfo, error) {
	lgr := repo.lgrf.NewLogger(ctx)
	res, err := repo.gttp.Get(
		ctx,
		map[string]string{
			"x-api-key": repo.optn.ApiKey,
		},
		repo.optn.BaseUrl,
		map[string][]string{},
	)
	if err != nil {
		lgr.Error(
			"failed to run request for creating online meeting",
			zap.Error(err),
		)
		return base.JoinInfo{}, err
	}
	if res.StatusCode != 200 {
		return base.JoinInfo{}, fmt.Errorf("unexpected status code")
	}
	joinInfo := JoinInfo{}
	res.Unmarshal(&joinInfo)
	if joinInfo.JoinWebURL == "" {
		lgr.Error("join url was missing from response")
		return base.JoinInfo{}, fmt.Errorf("join url was missing from response")
	}

	ji := base.JoinInfo(joinInfo)
	return ji, nil
}

func NewNagaClient(
	gttp *gottp.HttpClient,
	lgrf *logger.LoggerFactory,
	optn *NagaOptions,
) *NagaClient {
	return &NagaClient{
		gttp: gttp,
		lgrf: lgrf,
		optn: optn,
	}
}
