package gh

import (
	"context"
	"fmt"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

type Client struct {
	gh *github.Client
}

/* func NewClient() *Client { */
// client := github.NewClient(nil)
//
// return &Client{
//   gh: client,
// }
/* } */

func NewClientWithToken(ctx context.Context, token string) *Client {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: token,
		},
	)
	tokenContext := oauth2.NewClient(ctx, tokenSource)
	client := github.NewClient(tokenContext)

	repos, _, _ := client.Repositories.List(ctx, "cian911", &github.RepositoryListOptions{Type: "private", Sort: "updated", Direction: "desc", ListOptions: github.ListOptions{Page: 1, PerPage: 100}})
	for _, repo := range repos {
		fmt.Println(*repo.Name)
	}

	return &Client{
		gh: client,
	}
}
