package github

import "context"

type GitHubService interface {
	CloneRepository(ctx context.Context, url, branch, path string) error
}
