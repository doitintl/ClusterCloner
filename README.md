# Cluster Cloner
Reads the Kubernetes clusters in one location and clones them into another (or just outputs JSON as a dry run), to/from GKE and Azure.

For usage, run  `clustercloner --help`

# Setup
## Supply credentials
- `export GOOGLE_APPLICATION_CREDENTIALS='xyz.json'`, where `xyz.json` is a filename for a service account with the Kubernetes Cluster Admin role (to read and create clusters)
- Create `.env` with Azure credentials. Base the file on `.env.tpl` for needed keys

## Required GitHub secrets
(Not needed for local build.)
For GitHub CI, please specify the following GitHub secrets:
1. `DOCKER_USERNAME` - Docker Registry username
2. `DOCKER_PASSWORD` - Docker Registry password or token
3. `DOCKER_REGISTRY` - _optional_; Docker Registry name, default to `docker.io`
4. `DOCKER_REPOSITORY` - _optional_; Docker image repository name, default to `$GITHUB_REPOSITORY` (i.e. `user/repo`)

## Build
Use naked `make` command. You can also run `make help` to see more build targets.

For just the Docker build, run  `DOCKER_BUILDKIT=1 docker build -t <TAG> .`

## Docker
This uses Docker both as a CI tool and for releasing a final Docker image.

## Continuous Integration
Th GitHub `Docker CI` action is used.

# Credits
This project was started from the [goapp](https://github.com/alexei-led/goapp) template, a bootstrap project for Go CLI applications.
