package github

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
)

type CloneConfig struct {
	URL    string
	Branch string
	Dir    string
}

func (c *Client) CloneRepository(
	ctx context.Context,
	cfg CloneConfig,
) error {
	args := []string{
		"clone",
		"--depth", "1",
	}

	if cfg.Branch != "" {
		args = append(
			args,
			"--branch",
			cfg.Branch,
		)
	}

	args = append(
		args,
		cfg.URL,
		cfg.Dir,
	)

	cmd := exec.CommandContext(
		ctx,
		"git",
		args...,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"clone repository: %w: %s",
			err,
			string(output),
		)
	}

	return nil
}

func TempRepoDir(baseDir, projectID string) string {
	return filepath.Join(
		baseDir,
		projectID,
	)
}