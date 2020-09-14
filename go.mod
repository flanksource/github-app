module github.com/flanksource/github-app

go 1.14

require (
	github.com/alexedwards/scs v1.4.1
	github.com/bluekeyes/hatpear v0.0.0-20180714193905-ffb42d5bb417
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/flanksource/build-tools v0.9.10
	github.com/flanksource/commons v1.2.0
	github.com/go-oauth2/oauth2/v4 v4.1.0
	github.com/google/go-github/v32 v32.0.0
	github.com/palantir/go-baseapp v0.2.0

	//TODO: switch to next release of go-githubapp when it happens
	//      we need go-githubapp v32 for the fixes to ListRepositoryWorkflowRuns options
	//      https://github.com/google/go-github/issues/1497
	//      but latest go-githubapp release still using v31
	github.com/palantir/go-githubapp v0.3.1-0.20200530154104-bd812e979e03
	github.com/pkg/errors v0.9.1
	github.com/pquerna/cachecontrol v0.0.0-20180517163645-1555304b9b35 // indirect
	github.com/rs/zerolog v1.19.0
	github.com/spf13/cobra v1.0.0
	goji.io v2.0.2+incompatible
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	gopkg.in/flanksource/yaml.v3 v3.1.1
)

replace gopkg.in/hairyhenderson/yaml.v2 => github.com/maxaudron/yaml v0.0.0-20190411130442-27c13492fe3c
