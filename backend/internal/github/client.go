package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v60/github"
	"golang.org/x/oauth2"
)

type RepositoryItem struct {
	FullName      string `json:"full_name"`
	CloneURL      string `json:"clone_url"`
	DefaultBranch string `json:"default_branch"`
	Private       bool   `json:"private"`
}

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) getGhClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

// ListRepositories lấy danh sách repo của user
func (c *Client) ListRepositories(ctx context.Context, token string) ([]RepositoryItem, error) {
	gh := c.getGhClient(ctx, token)
	repos, _, err := gh.Repositories.ListByAuthenticatedUser(ctx, &github.RepositoryListByAuthenticatedUserOptions{
		Visibility:  "all",
		Sort:        "updated",
		ListOptions: github.ListOptions{PerPage: 50},
	})
	if err != nil {
		return nil, err
	}

	var res []RepositoryItem
	for _, r := range repos {
		res = append(res, RepositoryItem{
			FullName:      r.GetFullName(),
			CloneURL:      r.GetCloneURL(),
			DefaultBranch: r.GetDefaultBranch(),
			Private:       r.GetPrivate(),
		})
	}
	return res, nil
}

// ListBranches lấy danh sách branch của 1 repo
func (c *Client) ListBranches(ctx context.Context, token, fullName string) ([]string, error) {
	gh := c.getGhClient(ctx, token)
	parts := strings.Split(fullName, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("fullName khong hop le: %s", fullName)
	}

	branches, _, err := gh.Repositories.ListBranches(ctx, parts[0], parts[1], &github.BranchListOptions{})
	if err != nil {
		return nil, err
	}

	var list []string
	for _, b := range branches {
		list = append(list, b.GetName())
	}
	return list, nil
}

// CreateWebhook tự gắn webhook push event vào repo
func (c *Client) CreateWebhook(ctx context.Context, token, fullName, webhookURL, secret string) (int64, error) {
	gh := c.getGhClient(ctx, token)
	parts := strings.Split(fullName, "/")
	if len(parts) != 2 {
		return 0, fmt.Errorf("fullName khong hop le")
	}

	hook := &github.Hook{
		Name:   github.String("web"),
		Active: github.Bool(true),
		Events: []string{"push"},
		Config: &github.HookConfig{
			URL:         github.String(webhookURL),
			ContentType: github.String("json"),
			Secret:      github.String(secret),
		},
	}

	created, _, err := gh.Repositories.CreateHook(ctx, parts[0], parts[1], hook)
	if err != nil {
		return 0, err
	}
	return created.GetID(), nil
}