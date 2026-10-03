package github

import (
	"context"
	"fmt"
	githubapi "github.com/google/go-github/v60/github"
)

type GithubRepository interface {
	GetRepository(
		ctx context.Context,
		owner string,
		repo string,
	) (*githubapi.Repository, error)

	GetUser(
		ctx context.Context,
	) (*githubapi.User, error)
	ListRepositories(
		ctx context.Context,
	) ([]*githubapi.Repository, error)
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

func (c *Client) ListRepositories(
	ctx context.Context,
) ([]*githubapi.Repository, error) {
	opts := &githubapi.RepositoryListOptions{
		Sort:      "updated",
		Direction: "desc",
		ListOptions: githubapi.ListOptions{
			PerPage: 100,
		},
	}

	repos, _, err := c.GitHub.Repositories.List(ctx, "", opts)
	if err != nil {
		return nil, fmt.Errorf("list github repositories: %w", err)
	}

	return repos, nil
}
