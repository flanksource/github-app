/*
Copyright © 2020 Flanksource
This file is part of Flanksource github-app
*/
package server

import (
	"fmt"
	"github.com/alexedwards/scs"
	"github.com/bluekeyes/hatpear"
	//"github.com/die-net/lrucache"
	"github.com/flanksource/github-app/server/handlers"
	"github.com/flanksource/github-app/version"

	"github.com/palantir/go-baseapp/baseapp"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/palantir/go-githubapp/oauth2"
	"github.com/pkg/errors"
	//"github.com/rs/zerolog"
	goji "goji.io"
	"goji.io/pat"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	DefaultSessionLifetime = 24 * time.Hour
)

type Server struct {
	config *Config
	base   *baseapp.Server
}

// New instantiates a new Server.
// Callers must then invoke Start to run the Server.
func New(c *Config) (*Server, error) {
	logger := baseapp.NewLogger(baseapp.LoggingConfig{
		Level:  c.Logging.Level,
		Pretty: c.Logging.Text,
	})

	lifetime, _ := time.ParseDuration(c.Sessions.Lifetime)
	if lifetime == 0 {
		lifetime = DefaultSessionLifetime
	}

	publicURL, err := url.Parse(c.Server.PublicURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed parse public URL")
	}

	basePath := strings.TrimSuffix(publicURL.Path, "/")
	forceTLS := publicURL.Scheme == "https"

	sessions := scs.NewCookieManager(c.Sessions.Key)
	sessions.Name("policy-bot")
	sessions.Lifetime(lifetime)
	sessions.Persist(true)
	sessions.HttpOnly(true)
	sessions.Secure(forceTLS)

	base, err := baseapp.NewServer(c.Server, baseapp.DefaultParams(logger, "flanksource-githubapp.")...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize base server")
	}

	//maxSize := int64(50 * datasize.MB)
	//if c.Cache.MaxSize != 0 {
	//	maxSize = int64(c.Cache.MaxSize)
	//}
	//
	userAgent := fmt.Sprintf("%s/%s", "flanksource-github-app", version.GetVersion())
	//cc, err := githubapp.NewDefaultCachingClientCreator(
	//	c.Github,
	//	githubapp.WithClientUserAgent(userAgent),
	//	githubapp.WithClientCaching(true, func() httpcache.Cache {
	//		return lrucache.New(maxSize, 0)
	//	}),
	//	githubapp.WithClientMiddleware(
	//		githubapp.ClientLogging(zerolog.DebugLevel),
	//		githubapp.ClientMetrics(base.Registry()),
	//	),
	//)
	ghc := c.Github
	cc := githubapp.NewClientCreator(ghc.V3APIURL, ghc.V4APIURL, ghc.App.IntegrationID, []byte(ghc.App.PrivateKey), githubapp.WithClientUserAgent(userAgent))

	//if err != nil {
	//	return nil, errors.Wrap(err, "failed to initialize client creator")
	//}

	appClient, err := cc.NewAppClient()
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize Github app client")
	}

	//TODO: how to do this
	//registerOAuth2Handler(c.Github)

	//checkRunHandler := &handlers.StatusHandler{
	//	ClientCreator: cc,
	//}

	checkSuiteHandler := &handler.CheckSuiteHandler{
		ClientCreator: cc,
		Installations: githubapp.NewInstallationsService(appClient),
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
	mux.Handle(pat.Get("/api/health"), handler.Health())
	mux.Handle(pat.Get(oauth2.DefaultRoute), oauth2.NewHandler(
		oauth2.GetConfig(c.Github, nil),
		oauth2.ForceTLS(forceTLS),
		oauth2.WithStore(&oauth2.SessionStateStore{
			Sessions: sessions,
		}),
		oauth2.OnLogin(handler.Login(c.Github, sessions)),
	))

	//mux.Handle(pat.Get("/dispense/github-runner-token"), apiDispenseGithubRunnerToken)
	gh_runners := goji.SubMux()
	//gh_runners.Use(handler.RequireLogin(sessions, basePath))
	gh_runners.Handle(pat.Get("/github-runner-token"), hatpear.Try(&handler.GHRunners{
		ClientCreator: cc,
		Installations: githubapp.NewInstallationsService(appClient),
		Sessions:      sessions,
	}))
	mux.Handle(pat.New("/dispense/*"), gh_runners)

	return &Server{
		config: c,
		base:   base,
	}, nil
}

func Start() {

}

// Start is blocking and long-running
func (s *Server) Start() error {

	return s.base.Start()
}

func apiDispenseGithubRunnerToken(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, string("{\"echo\": \"hello\""))
}

//func registerOAuth2Handler(c githubapp.Config) {
//	http.Handle("/api/auth/github", oauth2.NewHandler(
//		oauth2.GetConfig(c, []string{"user:email"}),
//		// force generated URLs to use HTTPS; useful if the app is behind a reverse proxy
//		oauth2.ForceTLS(true),
//		// set the callback for successful logins
//		oauth2.OnLogin(func(w http.ResponseWriter, r *http.Request, login *oauth2.Login) {
//			// look up the current user with the authenticated client
//			client := github.NewClient(login.Client)
//			user, _, _ := client.Users.Get(r.Context(), "")
//			// handle error, save the user, ...
//			logger.Infof("%v",user)
//
//			// redirect the user back to another page
//			http.Redirect(w, r, "/dashboard", http.StatusFound)
//		}),
//	))
//}
