package repos

import (
	"context"

	"eventSourcedBooks/pkg/domain/auth"
)

type AuthRepository struct {
}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{}
}

var _ auth.IRepository = (*AuthRepository)(nil)

func (repo *AuthRepository) Validate(
	ctx context.Context,
	token string,
) (claims map[string]string, err error) {
	return make(map[string]string), nil
}
