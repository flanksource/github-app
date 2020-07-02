/*
Copyright © 2020 Flanksource
This file is part of Flanksource github-app
*/
package handler

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/flanksource/github-app/config"
	"github.com/go-oauth2/oauth2/v4"
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
	Cfg *config.Config
	// base oath2 server
	srv server.Server
}

// Init initializes the auth server, setting its manager, client and token store
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

	// HS512 is an Abbreviation of HMAC using SHA-512
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("kid", []byte("22222222"),jwt.SigningMethodHS512))

	// ES512 is ECDSA using P-521 and SHA-512 - public/private
	//manager.MapAccessGenerate(generates.NewJWTAccessGenerate("kid", []byte("22222222"),jwt.SigningMethodES512))

	//srv.ClientAuthorizedHandler = clientAuthorizedHandler
	srv.ClientScopeHandler = clientScopeHandler
	srv.SetExtensionFieldsHandler(func(ti oauth2.TokenInfo) (fieldsValue map[string]interface{}) {
		dummy := map[string]interface{}{
			"test": "test",
		}
		return dummy
	})


	as.srv = *srv

	return nil

}

func (as *AuthServer) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	err := as.srv.HandleTokenRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func (as *AuthServer) generateClientAuthorizedHandler() server.ClientAuthorizedHandler {

	clientAuthorizedHandler := func(clientID string, grant oauth2.GrantType) (allowed bool, err error) {

		if grant == oauth2.ClientCredentials {
			return true, nil
		}
		return false, nil
	}
	return clientAuthorizedHandler
}


func clientScopeHandler(clientID, scope string) (allowed bool, err error) {
	if scope=="runner" {
		return true, nil
	}
	return false, nil
}



