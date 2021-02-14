package gh

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-github/v33/github"
)

func TestWorkflowRuns(t *testing.T) {
	var (
		ctx   = context.Background()
		owner = "Cian911"
		repo  = "gomerge"
		token = os.Getenv("GITHUB_TOKEN")
		opts  = &github.ListWorkflowRunsOptions{ListOptions: github.ListOptions{Page: 1, PerPage: 25}}
		c     = NewClientWithToken(ctx, token)
	)

	t.Run("It returns a list of recent workflow runs", func(t *testing.T) {
		runs, _, _ := c.RecentWorkflowRuns(ctx, owner, repo, opts)

		if len(runs.WorkflowRuns) == 0 {
			t.Errorf("Expected runs count to be greater than 0, but it wasn't: %v", runs)
		}
	})

	t.Run("It throws an error when wrong data has been provided", func(t *testing.T) {
		cc := NewClientWithToken(ctx, "")
		_, resp, err := cc.RecentWorkflowRuns(ctx, owner, repo, opts)
		want := 401

		if want != resp.StatusCode {
			t.Errorf("Expected 401 error code but got something else: %d - %v", resp.StatusCode, err)
		}
	})

	t.Run("It returns a single workflow run", func(t *testing.T) {
		runId := int64(474347037)
		run, _, _ := c.WorkflowRunById(ctx, owner, repo, runId)

		if *run.ID != runId {
			fmt.Errorf("Expected 1 individual run, got something else: %v", run)
		}
	})

	t.Run("It returns a 404 when run does not exist", func(t *testing.T) {
		runId := int64(11111111)
		_, resp, err := c.WorkflowRunById(ctx, owner, repo, runId)
		want := 404

		if resp.StatusCode != want {
			t.Errorf("Expected a 404, got something else: %d - %v", resp.StatusCode, err)
		}
	})
}
