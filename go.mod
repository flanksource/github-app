module github.com/flanksource/github-app

go 1.14

require (
	github.com/flanksource/build-tools v0.9.10
	github.com/flanksource/commons v1.2.0
	github.com/google/go-github/v30 v30.1.0
	github.com/google/go-github/v31 v31.0.0
	github.com/google/go-github/v32 v32.0.0
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79
	github.com/palantir/go-baseapp v0.2.0

	//TODO: switch to next release of go-githubapp when it happens
	//      we need go-githubapp v32 for the fixes to ListRepositoryWorkflowRuns options
	//      https://github.com/google/go-github/issues/1497
	//      but latest go-githubapp release still using v31
	github.com/palantir/go-githubapp v0.3.1-0.20200530154104-bd812e979e03
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.19.0
	github.com/spf13/cobra v1.0.0
	goji.io v2.0.2+incompatible
	google.golang.org/appengine v1.6.6 // indirect
	gopkg.in/flanksource/yaml.v3 v3.1.1
)

replace gopkg.in/hairyhenderson/yaml.v2 => github.com/maxaudron/yaml v0.0.0-20190411130442-27c13492fe3c
