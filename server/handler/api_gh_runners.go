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

	if !Authorize(ctx,r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}


	//client := getPatClient(ctx, h.Secrets.GhPat)
	//
	//token, _, err := client.Actions.CreateRegistrationToken(ctx, h.Runners.Owner, h.Runners.Repo)
	//if err != nil {
	//	return err
	//}

	//fmt.Fprintf(w, "{\"registration token\": \"%s\"}", token.GetToken())
	fmt.Fprintf(w, "{\"registration token\": \"%s\"}", "dummy")

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

type user struct {
	email  string
	groups []string
}

// Authorize verifies a bearer token exists
// and has a valid JWT
func Authorize(ctx context.Context, r *http.Request) bool {
	auth := r.Header.Get("Authorization")
	if auth == "" || !strings.HasPrefix(auth,"Bearer ") {
		return false
	}
	tokenString := strings.TrimPrefix(auth,"Bearer ")
	if !ValidateJwt(tokenString){
		return false
	}
	return true
}

// ValidateJwt
func ValidateJwt(tokenString string) bool {
	token, err := jwt.Parse(tokenString, findTokenKey)
	if err != nil {
		return false
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["foo"], claims["nbf"])
		return true
	}
	return false

}

// findTokenKey examines the JWT and using the header Key ID finds a key
func findTokenKey(token *jwt.Token) (interface{}, error) {
	// Don't forget to validate the alg is what you expect:
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}

	// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
	return []byte("22222222"), nil
}

