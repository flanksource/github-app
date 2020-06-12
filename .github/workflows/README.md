# Workflows

## docker

Builds the `Dockerfile` and pushes to a docker hub image 
named: *github_owner*/*repo_name*:*tag_name*

(i.e. assumes that the docker hub account and github account match)

Usage:

e.g.
```bash
git tag v0.1.1
git push --tags
```

### Configs
Create repository secrets for:
* `DOCKER_USERNAME`
* `DOCKER_PASSWORD`

This user/credential needs access on docker hub to *github_owner*/*repo_name*

## binary

Builds the executable binary named `github-app` for linux and `github-app_osx` for OSX 
using the **release** target in the `Makefile` (i.e. `make release`) and 
releases to the Github repo with a release named *tag_name*.

Usage:

e.g.
```bash
git tag v0.1.1
git push --tags
```

### Configs
*none*