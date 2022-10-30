// All domain level errors are to be defined in this file

package common

import "github.com/betalixt/gorr"

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
