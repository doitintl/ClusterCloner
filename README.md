# Cluster Cloner
Reads a description of Google Kubernetes Engine cluster for later cloning.

# permissions

Permissions:

- `export GOOGLE_APPLICATION_CREDENTIALS='xyz.json'`, where `xyz.json` is a filename for a service account with  Kubernetes Cluster Admin (to create clusters)

# goapp

Based on goapp, the bootstrap project for Go CLI application.

## Docker

The `goapp` uses Docker both as a CI tool and for releasing final `goapp` Docker image (`scratch` with updated `ca-credentials` package).

## Makefile

The `goapp` `Makefile` is used for task automation only: compile, lint, test and etc.

## Continuous Integration

GitHub action `Docker CI` is used for `goapp` CI.

### Required GitHub secrets

Please specify the following GitHub secrets:

1. `DOCKER_USERNAME` - Docker Registry username
2. `DOCKER_PASSWORD` - Docker Registry password or token
3. `DOCKER_REGISTRY` - _optional_; Docker Registry name, default to `docker.io`
4. `DOCKER_REPOSITORY` - _optional_; Docker image repository name, default to `$GITHUB_REPOSITORY` (i.e. `user/repo`)
