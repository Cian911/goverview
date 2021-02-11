package github

import (
	gh "github.com/google/go-github/github"
)

type Client struct {
	ghClient *gh.Client
}
