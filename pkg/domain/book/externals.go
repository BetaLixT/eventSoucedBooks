package book

import (
	"context"
)

type IRepository interface {
	CreateBook(
		ctx context.Context,
		title string,
		description string,
		author string,
		genres []string,
	) error
}
