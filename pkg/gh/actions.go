package gh

import (
	"context"

	"github.com/google/go-github/v33/github"
)

func (c *Client) RecentWorkflowRuns(ctx context.Context, owner, repo string, opts *github.ListWorkflowRunsOptions) (*github.WorkflowRuns, *github.Response, error) {
	workflowRuns, resp, err := c.gh.Actions.ListRepositoryWorkflowRuns(ctx, owner, repo, opts)
	return workflowRuns, resp, err
}

func (c *Client) WorkflowRunById(ctx context.Context, owner, repo string, runId int64) (*github.WorkflowRun, *github.Response, error) {
	run, resp, err := c.gh.Actions.GetWorkflowRunByID(ctx, owner, repo, runId)
	return run, resp, err
}
