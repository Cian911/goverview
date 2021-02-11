package github

import (
	"context"
	"log"

	. "github.com/google/go-github/github"
)

func ListRecentWorkflows(ctx context.Context, owner, repo string, opts *gh.ListOptions) (*gh.Workflows, *gh.Response, error) {
	workflows, res, err := gh.Actions.ListWorkflows(ctx, owner, repo, opts)
	if err != nil {
		log.Fatal(err)
	}

	return workflows, res, nil
}
