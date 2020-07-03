/*
Copyright © 2020 Flanksource
This file is part of Flanksource github-app
*/
package handler

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/flanksource/github-app/config"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"log"
	"net/http"
)

// AuthServer provides a basic OAuth 2.0 auth server implementing the
// Client Credentials grant type.
type AuthServer struct {
	// Cfg contains the app config
	Cfg *config.Config
	// srv is base oath2 server that implements functionality
	srv server.Server
}

// Init initializes the auth server, setting up its manager, client and token store
// and initializing the signing key and allowed clients from the config.
func (as *AuthServer) Init() error {
	clientStore := store.NewClientStore()
	for _, cspec := range as.Cfg.Auth.Clients {
		c := cspec.GetClient()
		clientStore.Set(c.ID, c)
	}
	manager := manage.NewDefaultManager()
	manager.MapClientStorage(clientStore)
	manager.MustTokenStorage(store.NewMemoryTokenStore())
	srv := server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})
	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})
	// HS512 is an Abbreviation of HMAC using SHA-512 - i.e. symmetric encryption
	// used in this case because it's common, fast enough and we will be verifying
	// tokens that we generated ourselves. So we don't have a key distribution
	// problem.
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte(as.Cfg.Auth.SymmetricKey),jwt.SigningMethodHS512))
	srv.ClientScopeHandler = clientScopeHandler
	as.srv = *srv
	return nil
}

// ServeHTTP is the actual handler for OAuth 2.0 Client Credentials grant type
// requests
func (as *AuthServer) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	err := as.srv.HandleTokenRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

// clientScopeHandler decides if a given client/scope is allowed access.
func clientScopeHandler(clientID, scope string) (allowed bool, err error) {
	// for now we are only handling the runner scope and nothing else
	if scope=="runner" {
		return true, nil
	}
	return false, nil
}




