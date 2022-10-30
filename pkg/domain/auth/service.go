package auth

import "context"

type AuthService struct {
	repo IRepository
}

func NewAuthService(repo IRepository) *AuthService {
	return &AuthService{repo}
}

func (auth *AuthService) ValidateToken(
	ctx context.Context,
	token string,
) (claims map[string]string, err error) {
	claims, err = auth.repo.Validate(ctx, token)
	return
}
