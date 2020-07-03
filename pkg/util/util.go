package util

import (
	"github.com/google/go-github/v32/github"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

// GetPatClient returns a github client that uses the given
// Personal Access Token to authenticate
func GetPatClient(pat string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: pat},
	)
	tc := oauth2.NewClient(context.TODO(), ts)
	return github.NewClient(tc)
}

