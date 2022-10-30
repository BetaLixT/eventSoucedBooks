package book

import (
	"context"
	"eventSourcedBooks/pkg/domain/base"
)

type BookService struct {
	lgrf base.ILoggerFactory
	repo IRepository
}

func NewBookService(
	lgrf base.ILoggerFactory,
	repo IRepository,
) *BookService {
	return &BookService{
		lgrf: lgrf,
		repo: repo,
	}
}

func (s *BookService) CreateBook(
	ctx context.Context,
	cmd CreateBookCommand,
) error {
	fls := false
	zero := float32(0.0)
	return s.repo.CreateBook(
		ctx,
		BookData{
			Title:           &cmd.Title,
			Description:     &cmd.Description,
			Author:          &cmd.Author,
			Genres:          cmd.Genres,
			Completed:       &fls,
			CurrentPosition: &zero,
		},
	)
}

func (s *BookService) DeleteBook(
	ctx context.Context,
	cmd DeleteBookCommand,
) error {
	last, err := s.repo.LastEvent(ctx, cmd.Id)
	if err != nil {
		lgr := s.lgrf.NewLogger(ctx)
		lgr.Error("failed to fetch last event")
		return err
	}
	if last.Event.Event == base.DOMAIN_DELETE_EVENT {
		lgr := s.lgrf.NewLogger(ctx)
		lgr.Error("book was deleted")
		return base.NewBookMissingError()
	}
	err = s.repo.DeleteBook(
		ctx,
		cmd.Id,
		last.Event.Version+1,
	)
	if err != nil {
		lgr := s.lgrf.NewLogger(ctx)
		lgr.Error("failed to delete book")
	}
	return err
}

func (s *BookService) UpdatePagePosition(
	ctx context.Context,
	cmd UpdateBookCompletedCommand,
) error {

}

func (s *BookService) ToggleComletion(
	ctx context.Context,
	cmd ToggleCompletionCommand,
) error {

}
