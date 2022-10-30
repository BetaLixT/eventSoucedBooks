package base

import "github.com/betalixt/gorr"

func NewTokenMissingError() *gorr.Error {
	return gorr.NewError(
		gorr.ErrorCode{
			Code:    10008,
			Message: "TokenMissingError",
		},
		401,
		"",
	)
}

func NewTokenFormatInvalidError() *gorr.Error {
	return gorr.NewError(
		gorr.ErrorCode{
			Code:    10009,
			Message: "TokenFormatInvalidError",
		},
		401,
		"",
	)
}

func NewPropertyMissingError() *gorr.Error {
	return gorr.NewError(
		gorr.ErrorCode{
			Code:    20001,
			Message: "MissingProperties",
		},
		400,
		"one or more required properties were missing",
	)
}

func NewUnbindableError() *gorr.Error {
	return gorr.NewError(
		gorr.ErrorCode{
			Code:    20002,
			Message: "UnbindableBody",
		},
		400,
		"the request body was not bindable",
	)
}
