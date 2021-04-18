package gh

import (
	"context"

	"github.com/google/go-github/v35/github"
)

type RecentRuns struct {
	Repository string
	Runs       []*github.WorkflowRun
}

func (c *Client) UserRepos(ctx context.Context, user string, opts *github.RepositoryListOptions) ([]*github.Repository, *github.Response, error) {
	repos, resp, err := c.gh.Repositories.List(ctx, user, opts)
	return repos, resp, err
}

func (c *Client) OrganizationRepos(ctx context.Context, user string, opts *github.RepositoryListByOrgOptions) ([]*github.Repository, *github.Response, error) {
	repos, resp, err := c.gh.Repositories.ListByOrg(ctx, user, opts)
	return repos, resp, err
}
