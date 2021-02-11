package github

import (
	"context"
	"fmt"
	"testing"

	gh "github.com/google/go-github/github"
)

func TestListRecentWorkflows(t *testing.T) {
	var (
		ctx   = context.Background()
		owner = "Cian911"
		repo  = "gomerge"
	)

	t.Run("It returns a list of recent workflows", func(t *testing.T) {
		opts := &gh.ListOptions{Page: 1, PerPage: 5}
		workflows, res, err := ListRecentWorkflows(ctx, owner, repo, opts)
		fmt.Println(*workflows)
	})
}
