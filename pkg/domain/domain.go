package domain

import (
	"eventSourcedBooks/pkg/domain/auth"
	"eventSourcedBooks/pkg/domain/book"

	"github.com/betalixt/gorr"
	"github.com/google/wire"
)

var DependencySet = wire.NewSet(
	auth.NewAuthService,
	book.NewBookService,
)

// - Constants
// All domain level constants
const (
	APP_CLAIMS_KEY             = "claims"
	NOTIFICATIONS_SERVICE_NAME = "eventSourcedBooks"
	NOTIFICATIONS_EX_NAME      = "notifications"
	NOTIFICATIONS_EX_TYPE      = "topic"

	DOMAIN_CREATE_EVENT = "create"
	DOMAIN_UPDATE_EVENT = "update"
	DOMAIN_DELETE_EVENT = "delete"
)

// - Errors
// All of the domain level errors
// first digit identifies the layer (2 = domain)
// the first two digit identify the domain the error
// was created for 99 refers to a non domain specific error

func NewBookMissingError() *gorr.Error {
	return gorr.NewError(
		gorr.ErrorCode{
			Code:    2_10_000,
			Message: "BookMissingError",
		},
		404,
		"Book does not exist or deleted",
	)
}

func NewTokenMissingError() *gorr.Error {
	return gorr.NewError(
		gorr.ErrorCode{
			Code:    2_99_000,
			Message: "TokenMissingError",
		},
		401,
		"",
	)
}

func NewTokenFormatInvalidError() *gorr.Error {
	return gorr.NewError(
		gorr.ErrorCode{
			Code:    2_99_001,
			Message: "TokenFormatInvalidError",
		},
		401,
		"",
	)
}

func NewPropertyMissingError() *gorr.Error {
	return gorr.NewError(
		gorr.ErrorCode{
			Code:    2_99_002,
			Message: "MissingProperties",
		},
		400,
		"one or more required properties were missing",
	)
}

func NewUnbindableError() *gorr.Error {
	return gorr.NewError(
		gorr.ErrorCode{
			Code:    2_99_003,
			Message: "UnbindableBody",
		},
		400,
		"the request body was not bindable",
	)
}
