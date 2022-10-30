package book

import (
)

type BookService struct {
	repo IRepository
}

func NewBookService(
	repo IRepository,
) *BookService {
	return &BookService{
		repo: repo,
	}
}
