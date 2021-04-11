package gh

import (
	"context"

	"github.com/google/go-github/v33/github"
)

func (c *Client) OrganizationRepos(ctx context.Context, user string, opts *github.RepositoryListOptions) ([]*github.Repository, *github.Response, error) {
	repos, resp, err := c.gh.Repositories.List(ctx, user, opts)
	return repos, resp, err
}
