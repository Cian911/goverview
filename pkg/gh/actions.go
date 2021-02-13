package gh

import (
	"context"
	"fmt"

	"github.com/google/go-github/v33/github"
)

func (c *Client) ListRecentWorkflows(ctx context.Context, owner, repo string, opts *github.ListOptions) {
	fmt.Println(c.gh.Actions.ListWorkflows(ctx, owner, repo, opts))
	// workflows, res, err := *c.ghClient.Actions.ListWorkflows(ctx, owner, repo, opts)
	// if err != nil {
	// log.Fatal(err)
	// }

	// fmt.Println(workflows)
}
