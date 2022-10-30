package book

import (
	"context"
	"time"

	"eventSourcedBooks/pkg/domain/base"
	"eventSourcedBooks/pkg/infra/standard"

	"go.uber.org/zap"
)

type CourtroomService struct {
	lgrf base.ILoggerFactory
	repo IRepository
}

func NewCourtroomService(
	lgrf base.ILoggerFactory,
	repo IRepository,
) *CourtroomService {
	return &CourtroomService{
		lgrf: lgrf,
		repo: repo,
	}
}

func (svc *CourtroomService) CreateCourtroom(
	ctx context.Context,
	title string,
	description string,
	dateTimeStart time.Time,
	dateTimeEnd time.Time,
	location string,
	categories []string,
	metadata map[string]interface{},
) (*Courtroom, error) {
	externalCmsId := (*string)(nil)
	if data, ok := metadata[standard.EXTERNAL_CMSID_TAG]; ok {
		if val, ok := data.(string); ok {
			externalCmsId = &val
		}
	}

	lgr := svc.lgrf.NewLogger(ctx)
	lgr.Info("creating courtroom", zap.String("title", title),
		zap.String("description", description),
		zap.Time("dateTimeStart", dateTimeStart),
		zap.Time("dateTimeEnd", dateTimeEnd),
		zap.String("location", location),
		zap.Strings("categories", categories),
		zap.Stringp("externalCmsId", externalCmsId),
		zap.Any("metadata", metadata))

	courtroom, ver, err := svc.repo.Create(
		ctx,
		title,
		description,
		dateTimeStart,
		dateTimeEnd,
		categories,
		location,
		externalCmsId,
		metadata,
	)
	if err != nil {
		lgr.Error("failed to create courtroom")
		return courtroom, err
	}

	err = svc.repo.NotifyCreated(
		ctx,
		*courtroom,
		ver,
	)

	return courtroom, err
}

func (svc *CourtroomService) GetCourtroom(
	ctx context.Context,
	id string,
) (*Courtroom, error) {
	lgr := svc.lgrf.NewLogger(ctx)
	lgr.Info("getting courtroom by id", zap.String("courtroomId", id))

	courtroom, err := svc.repo.Get(ctx, id)
	if err != nil {
		lgr.Error("failed to get courtroom")
		return courtroom, err
	}
	return courtroom, nil
}

func (svc *CourtroomService) DeleteCourtroom(
	ctx context.Context,
	id string,
) (*Courtroom, error) {
	lgr := svc.lgrf.NewLogger(ctx)
	lgr.Info("deleting courtroom", zap.String("courtroomId", id))

	courtroom, ver, err := svc.repo.Delete(ctx, id)
	if err != nil {
		lgr.Error("failed to delete courtroom")
		return courtroom, err
	}

	err = svc.repo.NotifyDeleted(
		ctx,
		*courtroom,
		ver,
	)

	return courtroom, err
}

func (svc *CourtroomService) UpdateCourtroom(
	ctx context.Context,
	id string,
	title *string,
	description *string,
	dateTimeStart *time.Time,
	dateTimeEnd *time.Time,
	location *string,
	categories []string,
	metadata map[string]interface{},
) (*Courtroom, error) {
	lgr := svc.lgrf.NewLogger(ctx)
	lgr.Info("updating courtroom", zap.String("id", id),
		zap.Stringp("title", title),
		zap.Stringp("description", description),
		zap.Timep("dateTimeStart", dateTimeStart),
		zap.Timep("dateTimeEnd", dateTimeEnd),
		zap.Stringp("location", location),
		zap.Strings("categories", categories),
		zap.Any("metadata", metadata))

	ol, nw, ver, err := svc.repo.Update(
		ctx,
		id,
		title,
		description,
		dateTimeStart,
		dateTimeEnd,
		categories,
		location,
		metadata,
	)

	if err != nil {
		lgr.Error("failed to update courtroom")
		return nil, err
	}

	err = svc.repo.NotifyUpdated(
		ctx,
		*ol,
		*nw,
		title,
		description,
		dateTimeStart,
		dateTimeEnd,
		categories,
		location,
		metadata,
		ver,
	)

	return nw, err
}

func (svc *CourtroomService) QueryCourtroom(
	ctx context.Context,
	page int,
	items int,
	search []base.SearchQuery,
	orderby string,
	desc bool,
) ([]Courtroom, error) {
	lgr := svc.lgrf.NewLogger(ctx)
	lgr.Info("querying courtroom", zap.Int("page", page),
		zap.Int("items", items),
		zap.Strings("search", base.FlattenSearchQueryList(search)),
		zap.String("orderby", orderby),
		zap.Bool("desc", desc))

	courtrooms, err := svc.repo.Query(
		ctx,
		page,
		items,
		search,
		orderby,
		desc,
	)
	if err != nil {
		lgr.Error("failed to query courtroom")
		return courtrooms, err
	}
	return courtrooms, nil
}

func (svc *CourtroomService) UpdateCourtroomMeta(
	ctx context.Context,
	id string,
	keyvals map[string]interface{},
) (*Courtroom, error) {
	lgr := svc.lgrf.NewLogger(ctx)
	lgr.Info("updating courtroom meta", zap.String("id", id),
		zap.Any("keyvals", keyvals))

	courtroom, ver, err := svc.repo.UpdateMeta(
		ctx,
		id,
		keyvals,
	)

	if err != nil {
		lgr.Error("failed to update courtroom meta")
		return courtroom, err
	}

	err = svc.repo.NotifyMetadataUpdated(
		ctx,
		*courtroom,
		ver,
	)

	return courtroom, err
}

func (svc *CourtroomService) DeleteCourtroomMeta(
	ctx context.Context,
	id string,
	keys []string,
) (*Courtroom, error) {
	lgr := svc.lgrf.NewLogger(ctx)
	lgr.Info("deleting courtroom meta", zap.String("id", id),
		zap.Strings("keys", keys))

	courtroom, ver, err := svc.repo.DeleteMeta(
		ctx,
		id,
		keys,
	)

	if err != nil {
		lgr.Error("failed to create courtroom meta")
		return courtroom, err
	}

	err = svc.repo.NotifyMetadataUpdated(
		ctx,
		*courtroom,
		ver,
	)

	return courtroom, err
}
