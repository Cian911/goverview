package gh

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-github/v33/github"
)

func TestListRecentWorkflows(t *testing.T) {
	var (
		ctx   = context.Background()
		owner = "Cian911"
		repo  = "gomerge"
		token = os.Getenv("GITHUB_TOKEN")
		opts  = &github.ListWorkflowRunsOptions{ListOptions: github.ListOptions{Page: 1, PerPage: 25}}
		c     = NewClientWithToken(ctx, token)
	)

	t.Run("It returns a list of recent workflow runs", func(t *testing.T) {
		runs, _ := c.RecentWorkflowRuns(ctx, owner, repo, opts)

		if len(runs.WorkflowRuns) == 0 {
			t.Errorf("Expected runs count to be greater than 0, but it wasn't: %v", runs)
		}
	})

	t.Run("It returns a single workflow run", func(t *testing.T) {
		runId := int64(474347037)
		run, _ := c.WorkflowRunById(ctx, owner, repo, runId)

		if *run.ID != runId {
			fmt.Errorf("Expected 1 individual run, got something else: %v", run)
		}
	})
}
