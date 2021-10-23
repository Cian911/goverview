package gh

import (
	"context"
	"time"

	"github.com/google/go-github/v35/github"
)

type RecentRuns struct {
	Repository string
	CreatedAt  time.Time
	Runs       []*github.WorkflowRun
}

type ByCreation []RecentRuns

/* func (a ByCreation) Len() int           { return len(a) } */
// func (a ByCreation) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
/* func (a ByCreation) Less(i, j int) bool { return a[i].CreatedAt < a[j].CreatedAt } */

func (c *Client) UserRepos(ctx context.Context, user string, opts *github.RepositoryListOptions) ([]*github.Repository, *github.Response, error) {
	repos, resp, err := c.gh.Repositories.List(ctx, user, opts)
	return repos, resp, err
}

func (c *Client) OrganizationRepos(ctx context.Context, user string, opts *github.RepositoryListByOrgOptions) ([]*github.Repository, *github.Response, error) {
	repos, resp, err := c.gh.Repositories.ListByOrg(ctx, user, opts)
	return repos, resp, err
}
