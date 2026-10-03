package auth

import (
	"context"

	"github.com/sowncns/k3s-deploy-platform/internal/auth/dto"
)

type AuthService interface {
	CreateOrUpdateUser(
		ctx context.Context,
		input dto.CreateAuthUser,
	) (*GitHubUser, error)
}

type Service struct {
	repo AuthRepository
}

func NewService(repo AuthRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateOrUpdateUser(
	ctx context.Context,
	input dto.CreateAuthUser,
) (*GitHubUser, error) {
	user := &GitHubUser{
		ID:          input.ID,
		Login:       input.Login,
		Name:        input.Name,
		Email:       input.Email,
		AvatarURL:   input.AvatarURL,
		AccessToken: input.AccessToken,
	}

	if err := s.repo.Upsert(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
