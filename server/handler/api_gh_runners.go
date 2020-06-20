/*
Copyright © 2020 Flanksource
This file is part of Flanksource github-app
*/
package handler

import (
	"context"
	"fmt"
	"github.com/flanksource/github-app/config"
	"golang.org/x/oauth2"

	"github.com/google/go-github/v32/github"
	"net/http"
)

// GHRunners is a handler for the API endpoint for Github Runner related functionality
type GHRunners struct {
	config.Config
}

func (h *GHRunners) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	client := getPatClient(ctx, h.Secrets.GhPat)

	token, _, err := client.Actions.CreateRegistrationToken(ctx, h.Runners.Owner, h.Runners.Repo)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "{\"registration token\": \"%s\"}", token.GetToken())

	return nil
}

// getPatClient returns a github client that uses the given
// Personal Access Token to authenticate
// NOTE: this is a workaround for issues experienced with using
//       githubapp.ClientCreator.NewAppClient()
func getPatClient(ctx context.Context, pat string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: pat},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}
