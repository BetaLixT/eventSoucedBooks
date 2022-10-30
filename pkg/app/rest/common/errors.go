package common

import "github.com/betalixt/gorr"

func NewTokenMissingError() *gorr.Error {
	return gorr.NewError(
		gorr.ErrorCode{
			Code:    1_99_000,
			Message: "TokenMissingError",
		},
		401,
		"",
	)
}

func NewTokenFormatInvalidError() *gorr.Error {
	return gorr.NewError(
		gorr.ErrorCode{
			Code:    1_99_001,
			Message: "TokenFormatInvalidError",
		},
		401,
		"",
	)
}

func NewPropertyMissingError() *gorr.Error {
	return gorr.NewError(
		gorr.ErrorCode{
			Code:    1_99_002,
			Message: "MissingProperties",
		},
		400,
		"one or more required properties were missing",
	)
}

func NewUnbindableError() *gorr.Error {
	return gorr.NewError(
		gorr.ErrorCode{
			Code:    1_99_003,
			Message: "UnbindableBody",
		},
		400,
		"the request body was not bindable",
	)
}
