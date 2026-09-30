package github

import (
	"context"
	"fmt"
	"github.com/sowncns/k3s-deploy-platform/internal/config"

	"golang.org/x/oauth2"
)

type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

func NewOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.Load().GitHubClientID,
		ClientSecret: config.Load().GitHubClientSecret,
		RedirectURL:  config.Load().GitHubRedirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
		Scopes: []string{
			"read:user",
			"user:email",
			"repo",
		},
	}
}

func ExchangeCode(
	ctx context.Context,
	config *oauth2.Config,
	code string,
) (*oauth2.Token, error) {
	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("exchange github oauth code: %w", err)
	}

	return token, nil
}