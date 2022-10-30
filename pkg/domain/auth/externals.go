package auth

import "context"

type IRepository interface {
	Validate(
		ctx context.Context,
		token string,
	) (claims map[string]string, err error)
}
