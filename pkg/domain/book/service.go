package book

import (
	"context"
	"eventSourcedBooks/pkg/domain/base"

	"go.uber.org/zap"
)

type BookService struct {
	lgrf base.ILoggerFactory
	txfc base.ITransactionFactory
	repo IRepository
}

func NewBookService(
	lgrf base.ILoggerFactory,
	txfc base.ITransactionFactory,
	repo IRepository,
) *BookService {
	return &BookService{
		lgrf: lgrf,
		txfc: txfc,
		repo: repo,
	}
}

func (s *BookService) CreateBook(
	ctx context.Context,
	cmd CreateBookCommand,
) error {
	fls := false
	zero := float32(0.0)

	tx, err := s.txfc.Create(ctx)
	if err != nil {
		lgr := s.lgrf.NewLogger(ctx)
		lgr.Error("failed to create transaction")
		return err
	}

	err = tx.Begin(ctx)
	if err != nil {
		lgr := s.lgrf.NewLogger(ctx)
		lgr.Error("failed to begin transaction")
		return err
	}

	err = s.repo.Create(
		ctx,
		tx,
		cmd.SagaId,
		BookData{
			Title:           &cmd.Title,
			Description:     &cmd.Description,
			Author:          &cmd.Author,
			Genres:          cmd.Genres,
			Completed:       &fls,
			CurrentPosition: &zero,
		},
	)
	if err != nil {
		lgr := s.lgrf.NewLogger(ctx)
		lgr.Error("failed to create book")
		tx.Rollback(ctx)
	} else {
		err = tx.Commit(ctx)
		if err != nil {
			lgr := s.lgrf.NewLogger(ctx)
			lgr.Error("failed to commit transaction")
			tx.Rollback(ctx)	
		}
	}
	return err
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

	err = s.repo.Delete(
		ctx,
		cmd.SagaId,
		cmd.Id,
		last.Event.Version+1,
	)
	if err != nil {
		lgr := s.lgrf.NewLogger(ctx)
		lgr.Error("failed to delete book")
	}
	return err
}

func (s *BookService) UpdateBookPosition(
	ctx context.Context,
	cmd UpdatePagePositionCommand,
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

	err = s.repo.Update(ctx,
		cmd.SagaId,
		cmd.Id,
		last.Event.Version+1,
		BookData{
			CurrentPosition: &cmd.Position,
		},
	)
	if err != nil {
		lgr := s.lgrf.NewLogger(ctx)
		lgr.Error("failed to update book")
	}
	return err
}

func (s *BookService) UpdateBookCompletion(
	ctx context.Context,
	cmd UpdateBookCompletedCommand,
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

	err = s.repo.Update(ctx,
		cmd.SagaId,
		cmd.Id,
		last.Event.Version+1,
		BookData{
			Completed: &cmd.Completed,
		},
	)
	if err != nil {
		lgr := s.lgrf.NewLogger(ctx)
		lgr.Error("failed to update book")
	}
	return err
}

func (s *BookService) ToggleBookCompletion(
	ctx context.Context,
	cmd ToggleCompletionCommand,
) error {
	evnts, err := s.repo.ListEvents(ctx, cmd.Id)
	if err != nil {
		lgr := s.lgrf.NewLogger(ctx)
		lgr.Error("failed to fetch events")
		return err
	}
	if len(evnts) == 0 {
		lgr := s.lgrf.NewLogger(ctx)
		lgr.Error("no events found")
		return base.NewBookMissingError()
	}
	if evnts[len(evnts)-1].Event.Event == base.DOMAIN_DELETE_EVENT {
		lgr := s.lgrf.NewLogger(ctx)
		lgr.Error("book was deleted")
		return base.NewBookMissingError()
	}
	state := BuildState(evnts)

	cmp := !state.Completed
	err = s.repo.Update(ctx,
		cmd.SagaId,
		cmd.Id,
		state.Version+1,
		BookData{
			Completed: &cmp,
		},
	)
	if err != nil {
		lgr := s.lgrf.NewLogger(ctx)
		lgr.Error("failed to update book")
	}
	return err
}

func BuildState(events []BookEvent) BookState {
	book := BookState{
		Version:         events[len(events)-1].Event.Version,
		DateTimeCreated: events[0].Event.EventTime,
		DateTimeUpdated: events[len(events)-1].Event.EventTime,
	}
	for _, evnt := range events {
		if evnt.Data.Title != nil {
			book.Title = *evnt.Data.Title
		}
		if evnt.Data.Description != nil {
			book.Description = *evnt.Data.Description
		}
		if evnt.Data.Author != nil {
			book.Author = *evnt.Data.Author
		}
		if evnt.Data.Genres != nil {
			book.Genres = evnt.Data.Genres
		}
		if evnt.Data.Completed != nil {
			book.Completed = *evnt.Data.Completed
		}
		if evnt.Data.CurrentPosition != nil {
			book.CurrentPosition = *evnt.Data.CurrentPosition
		}
	}
	return book
}
