/*
Copyright © 2020 Flanksource
This file is part of Flanksource github-app
*/
package handler

import (
	"fmt"
	"github.com/flanksource/github-app/config"
	"github.com/flanksource/github-app/pkg/util"
	"net/http"
)

// GHRunners is a handler for the API endpoint for Github Runner related functionality
type GHRunners struct {
	config.Config
}

func (h *GHRunners) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	client := util.GetPatClient(h.Secrets.GhPat)

	token, _, err := client.Actions.CreateRegistrationToken(ctx, h.Runners.Owner, h.Runners.Repo)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "{\"registration token\": \"%s\"}", token.GetToken())

	return nil
}

