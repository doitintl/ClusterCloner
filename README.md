# Cluster Cloner
Reads a description of Google Kubernetes Engine cluster for later cloning.

# permissions

Permissions:

- `export GOOGLE_APPLICATION_CREDENTIALS='xyz.json'`, where `xyz.json` is a filename for a service account with  Kubernetes Cluster Admin (to create clusters)
- Create `.env` with Azure credentials. See `.env.tpl` for needed keys
- Makes sure needed keys are in `.aws` as usual



# goapp
This was based on [goapp](https://github.com/alexei-led/goapp), the bootstrap project for Go CLI application.

## Docker
Build with `DOCKER_BUILDKIT=1 docker build -t <TAG> .`

Based on `goapp`, this uses Docker both as a CI tool and for releasing a final Docker image 
(based on `scratch` with updated `ca-credentials` package).

## Makefile

The `Makefile` is used for task automation: compile, lint, test and etc.

## Continuous Integration

GitHub action `Docker CI` is used.

### Required GitHub secrets

For GitHub CI, please specify the following GitHub secrets:

1. `DOCKER_USERNAME` - Docker Registry username
2. `DOCKER_PASSWORD` - Docker Registry password or token
3. `DOCKER_REGISTRY` - _optional_; Docker Registry name, default to `docker.io`
4. `DOCKER_REPOSITORY` - _optional_; Docker image repository name, default to `$GITHUB_REPOSITORY` (i.e. `user/repo`)
