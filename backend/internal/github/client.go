package github

import (
	"context"
	"fmt"

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

func (c *Client) GetRepository(
	ctx context.Context,
	owner string,
	repo string,
) (*githubapi.Repository, error) {
	result, _, err := c.GitHub.Repositories.Get(
		ctx,
		owner,
		repo,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"get github repository %s/%s: %w",
			owner,
			repo,
			err,
		)
	}

	return result, nil
}

func (c *Client) GetUser(
	ctx context.Context,
) (*githubapi.User, error) {
	user, _, err := c.GitHub.Users.Get(
		ctx,
		"",
	)
	if err != nil {
		return nil, fmt.Errorf(
			"get github user: %w",
			err,
		)
	}

	return user, nil
}