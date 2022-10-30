package domain

import (
	"eventSourcedBooks/pkg/domain/auth"
	"eventSourcedBooks/pkg/domain/book"
	"github.com/google/wire"
)

var DependencySet = wire.NewSet(
	auth.NewAuthService,
	book.NewBookService,
)
