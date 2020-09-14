/*
Copyright © 2020 Flanksource
This file is part of Flanksource github-app
*/
package handler

import (
	"context"
	"fmt"
	"github.com/coreos/go-oidc"
	"github.com/dgrijalva/jwt-go"
	"github.com/flanksource/github-app/config"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
	"net/http"
	"strings"
)

// GHRunners is a handler for the API endpoint for Github Runner related functionality
type GHRunners struct {
	config.Config
	verifier oidc.IDTokenVerifier
}

func (h *GHRunners) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if !Authorize(ctx,r, []byte(h.Auth.SymmetricKey)) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	client := getPatClient(ctx, h.Secrets.GhPat)
	token, _, err := client.Actions.CreateRegistrationToken(ctx, h.Runners.Owner, h.Runners.Repo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "{\"registration token\": \"%s\"}", token.GetToken())
	return
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

// Authorize verifies a bearer token exists
// and has a valid JWT
func Authorize(ctx context.Context, r *http.Request, key []byte) bool {
	auth := r.Header.Get("Authorization")
	if auth == "" || !strings.HasPrefix(auth,"Bearer ") {
		return false
	}
	tokenString := strings.TrimPrefix(auth,"Bearer ")
	if !ValidateJwt(tokenString, key){
		return false
	}
	return true
}

// ValidateJwt with a given key
func ValidateJwt(tokenString string, key []byte) bool {
	findTokenKey := func(token *jwt.Token) (interface{}, error) {
		return key, nil
	}
	token, err := jwt.Parse(tokenString, findTokenKey)
	if err != nil {
		return false
	}
	if token.Valid {
		return true
	}
	return false
}

