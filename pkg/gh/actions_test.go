package gh

import (
	"context"
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
	)

	t.Run("It returns a list of recent workflows", func(t *testing.T) {
		opts := &github.ListOptions{Page: 1, PerPage: 5}
		c := NewClientWithToken(ctx, token)

		c.ListRecentWorkflows(ctx, owner, repo, opts)
	})
}
