package auth

import (
	"context"
	"fmt"

	"github.com/sowncns/k3s-deploy-platform/internal/auth/dto"
	githubclient "github.com/sowncns/k3s-deploy-platform/internal/github"
)

type AuthService interface{
	LoginWithGitHub(
	ctx context.Context,
	code string,
) (*GitHubUser, error)
GetGitHubLoginURL() string 
}

type Service struct {
	GitHubOAuth *githubclient.OAuthConfig
	repo  AuthRepository
}

func NewService(
	githubOAuth *githubclient.OAuthConfig,
	repo AuthRepository,
) *Service {
	return &Service{
		GitHubOAuth: githubOAuth,
		repo : repo,
	}
}


func (s *Service) LoginWithGitHub(
	ctx context.Context,
	code string,
) (*GitHubUser, error) {

	oauthConfig := githubclient.NewOAuthConfig()

	token, err := githubclient.ExchangeCode(
		ctx,
		oauthConfig,
		code,
	)
	if err != nil {
		return nil, err
	}

	client := githubclient.NewClient(token.AccessToken)

	ghUser, err := client.GetUser(ctx)
	if err != nil {
		return nil, fmt.Errorf(
			"get github user: %w",
			err,
		)
	}

	user, err := s.repo.getUserGithub(ctx, ghUser.GetID())

	if err != nil {
		user, err = s.repo.Create(ctx,
			&dto.CreateAuthUser{
				ID:        ghUser.GetID(),
				Login:     ghUser.GetLogin(),
				Name:      ghUser.GetName(),
				Email:     ghUser.GetEmail(),
				AvatarURL: ghUser.GetAvatarURL(),
			},
		)

		if err != nil {
			return nil, fmt.Errorf(
				"create github user: %w",
				err,
			)
		}
	}

	return user, nil
}

func (s *Service) GetGitHubLoginURL() string {
	config := githubclient.NewOAuthConfig()

	return config.AuthCodeURL("github-login")
}