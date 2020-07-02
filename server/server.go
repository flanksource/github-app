/*
Copyright © 2020 Flanksource
This file is part of Flanksource github-app
*/
package server

import (
	"fmt"
	"github.com/bluekeyes/hatpear"
	"github.com/flanksource/commons/logger"
	cfg "github.com/flanksource/github-app/config"
	"github.com/flanksource/github-app/server/handler"
	"github.com/flanksource/github-app/version"
	"net/http"

	"github.com/rs/zerolog"

	"github.com/google/go-github/v32/github"
	"github.com/palantir/go-baseapp/baseapp"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/palantir/go-githubapp/oauth2"
	"github.com/pkg/errors"

	goji "goji.io"
	"goji.io/pat"
	"net/url"
	"strings"
	"time"
)

const (
	DefaultSessionLifetime = 24 * time.Hour
)

type Server struct {
	config *cfg.Config
	base   *baseapp.Server
}

// New instantiates a new Server.
// Callers must then invoke Start to run the Server.
func New(c *cfg.Config) (*Server, error) {
	logger := baseapp.NewLogger(baseapp.LoggingConfig{
		Level:  c.Logging.Level,
		Pretty: c.Logging.Text,
	})

	publicURL, err := url.Parse(c.Server.PublicURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed parse public URL")
	}

	basePath := strings.TrimSuffix(publicURL.Path, "/")

	base, err := baseapp.NewServer(c.Server, baseapp.DefaultParams(logger, "flanksource-githubapp.")...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize base server")
	}

	userAgent := fmt.Sprintf("%s/%s", "flanksource-github-app", version.GetVersion())
	cc, err := githubapp.NewDefaultCachingClientCreator(
		c.Github,
		githubapp.WithClientUserAgent(userAgent),
		githubapp.WithClientMiddleware(
			githubapp.ClientLogging(zerolog.DebugLevel),
			githubapp.ClientMetrics(base.Registry()),
		),
	)

	checkSuiteHandler := &handler.CheckSuiteHandler{
		ClientCreator: cc,
	}

	queueSize := c.Workers.QueueSize
	if queueSize < 1 {
		queueSize = 100
	}

	workers := c.Workers.Workers
	if workers < 1 {
		workers = 10
	}

	dispatcher := githubapp.NewEventDispatcher(
		[]githubapp.EventHandler{
			//checkRunHandler,
			checkSuiteHandler,
		},
		c.Github.App.WebhookSecret,
		githubapp.WithScheduler(
			githubapp.QueueAsyncScheduler(queueSize, workers, githubapp.WithSchedulingMetrics(base.Registry())),
		),
	)

	var mux *goji.Mux
	if basePath == "" {
		mux = base.Mux()
	} else {
		mux = goji.SubMux()
		base.Mux().Handle(pat.New(basePath+"/*"), mux)
	}

	// webhook route
	mux.Handle(pat.Post(githubapp.DefaultWebhookRoute), dispatcher)

	// additional API routes
	mux.Handle(pat.Get("/health"), hatpear.Try(&handler.HealthCheck{}))



	as := handler.AuthServer{Cfg: c}
	err = as.Init()
	if err != nil {
		return nil, err
	}
	mux.Handle(pat.Post("/token"), hatpear.Try(&as))
	mux.Handle(pat.Get("/token"), hatpear.Try(&as))
	//registerOAuth2Handler(c.Github)

	gh_runners := goji.SubMux()
	gh_runners.Handle(pat.Get("/github-runner-token"), &handler.GHRunners{Config: *c})
	mux.Handle(pat.New("/dispense/*"), gh_runners)

	return &Server{
		config: c,
		base:   base,
	}, nil
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
			logger.Infof("%v", user)

			// redirect the user back to another page
			http.Redirect(w, r, "/health", http.StatusFound)
		}),
	))
}
