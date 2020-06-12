/*
Copyright © 2020 Flanksource
This file is part of Flanksource github-app
*/
package server

import (
	"github.com/flanksource/commons/logger"
	"github.com/flanksource/github-app/config"
	"github.com/flanksource/github-app/server/handlers"
	"github.com/google/go-github/v31/github"
	"github.com/gregjones/httpcache"
	"github.com/palantir/go-baseapp/baseapp"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/palantir/go-githubapp/oauth2"
	"github.com/rs/zerolog"
	"goji.io/pat"
	"net/http"
)

type Server struct {
	config *config.Config
	base   *baseapp.Server
}

// New instantiates a new Server.
// Callers must then invoke Start to run the Server.
func New(config *config.Config, logger zerolog.Logger) (*baseapp.Server, error) {
	server, err := baseapp.NewServer(
		config.Server,
		baseapp.DefaultParams(logger, "flanksource-githubapp.")...,
	)
	if err != nil {
		panic(err)
	}

	cc, err := githubapp.NewDefaultCachingClientCreator(
		config.Github,
		githubapp.WithClientUserAgent("example-app/1.0.0"),
		githubapp.WithClientCaching(false, func() httpcache.Cache { return httpcache.NewMemoryCache() }),
		githubapp.WithClientMiddleware(
			githubapp.ClientMetrics(server.Registry()),
		),
	)
	if err != nil {
		panic(err)
	}

	registerOAuth2Handler(config.Github)

	checkRunHandler := &handlers.StatusHandler{
		ClientCreator: cc,
	}

	checkSuiteHandler := &handlers.CheckSuiteHandler{
		ClientCreator: cc,
	}

	webhookHandler := githubapp.NewEventDispatcher(
		[]githubapp.EventHandler{
			checkRunHandler,
			checkSuiteHandler,
		},
		config.Github.App.WebhookSecret,
	)

	mux := server.Mux()

	// webhook route
	mux.Handle(pat.Post(githubapp.DefaultWebhookRoute), webhookHandler)

	return server, nil
}

func Start() {

}

// Start is blocking and long-running
func (s *Server) Start() error {

	return s.base.Start()
}

func registerOAuth2Handler(c githubapp.Config) {
	http.Handle("/api/auth/github", oauth2.NewHandler(
		oauth2.GetConfig(c, []string{"user:email"}),
		// force generated URLs to use HTTPS; useful if the app is behind a reverse proxy
		oauth2.ForceTLS(true),
		// set the callback for successful logins
		oauth2.OnLogin(func(w http.ResponseWriter, r *http.Request, login *oauth2.Login) {
			// look up the current user with the authenticated client
			client := github.NewClient(login.Client)
			user, _, _ := client.Users.Get(r.Context(), "")
			// handle error, save the user, ...
			logger.Infof("%v",user)

			// redirect the user back to another page
			http.Redirect(w, r, "/dashboard", http.StatusFound)
		}),
	))
}
