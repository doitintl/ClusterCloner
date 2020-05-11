# Cluster Cloner
Reads the Kubernetes clusters in one location and clones them into another (or just output JSON as a dry run), to/from GKE and Azure.

For usage, run  `clustercloner --help`

# permissions
Permissions:

- `export GOOGLE_APPLICATION_CREDENTIALS='xyz.json'`, where `xyz.json` is a filename for a service account with  Kubernetes Cluster Admin (to create clusters)
- Create `.env` with Azure credentials. Base it on `.env.tpl` for needed keys

# goapp
This project used  [goapp](https://github.com/alexei-led/goapp), the bootstrap project for Go CLI application, as a template.

## Docker
Build with `make build` (which also runs formatting, linting, and testing) or 
for just building it, run `DOCKER_BUILDKIT=1 docker build -t <TAG> .`

This uses Docker both as a CI tool and for releasing a final Docker image 
(which is based on `scratch` with updated `ca-credentials` package).

## Makefile
The `Makefile` is used for task automation: build, format, lint, test etc.

## Continuous Integration

GitHub action `Docker CI` is used.

### Required GitHub secrets

For GitHub CI, please specify the following GitHub secrets:

1. `DOCKER_USERNAME` - Docker Registry username
2. `DOCKER_PASSWORD` - Docker Registry password or token
3. `DOCKER_REGISTRY` - _optional_; Docker Registry name, default to `docker.io`
4. `DOCKER_REPOSITORY` - _optional_; Docker image repository name, default to `$GITHUB_REPOSITORY` (i.e. `user/repo`)
