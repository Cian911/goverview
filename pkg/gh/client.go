package gh

import (
	"context"

	"github.com/google/go-github/v35/github"
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

	/* repos, resp, err := client.Repositories.List(ctx, "Cian911", &github.RepositoryListOptions{Type: "private", Sort: "updated", Direction: "desc"}) */
	// fmt.Println(resp)
	// fmt.Println(err)
	// for _, repo := range repos {
	//   fmt.Println(*repo.Name)
	/* } */

	return &Client{
		gh: client,
	}
}
