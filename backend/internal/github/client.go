package github

import (
	"context"
	githubapi "github.com/google/go-github/v60/github"
	"golang.org/x/oauth2"
)

type Client struct {
	GitHub *githubapi.Client
}

func NewClient(token string) *Client {
	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: token,
		},
	)

	httpClient := oauth2.NewClient(ctx, ts)

	return &Client{
		GitHub: githubapi.NewClient(httpClient),
	}
}

