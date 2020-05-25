# Cluster Cloner
Reads the Kubernetes clusters in one location (optionally filtering by labels) and
clones them into another (or just outputs JSON as a dry run), to/from GKE and Azure.

# Usage
For usage, run  `clustercloner --help`

# Setup
## Supply credentials
- Add a file `credentials-cluster-manager.json` with credentials for a service account with the Kubernetes Cluster Admin role (to read and create clusters).
This is loaded through the `GOOGLE_APPLICATION_CREDENTIALS` environment variable; see the `Dockerfile`.
- Add a file `.env` with Azure credentials. Use `.env.tpl` as a template. The user should have the  Azure Kubernetes Service Cluster Admin Role.
- Add a file `awscredentials` with AWS credentials. The user should have the policy
 discussed [here](https://docs.aws.amazon.com/eks/latest/userguide/security_iam_id-based-policy-examples.html). Specific example [here](https://github.com/weaveworks/eksctl/issues/204#issuecomment-631630355)

## Required secrets for Building in GitHub Continuous Integration
Not needed for local build.
Store the base-64 encodings, for example echo `my-credential.json |base64`
For GitHub CI, please specify the following GitHub secrets:
- `AZ_ENV_BASE64  `.env` file with Azure credentials, following  `.env.tpl` as a template.
- `GCP_CLUSTER_MANAGER_KEYJSON_BASE64` Google credentials file (JSON)  with role Kubernetes Cluster Manager. Used at runtime.
- `GCR_PUSHER_KEYJSON_BASE64` Google credentials file (JSON) with role Storage Admin for pushing  to your GCR registry
- `DOCKER_REGISTRY` - Point this to GCR
- `DOCKER_REPOSITORY` - _optional_; Docker image name including repository, default to `$GITHUB_REPOSITORY` (in the form `user/repo`)

## Build
The Docker image is built in Github Workflows. In development, you can run  `DOCKER_BUILDKIT=1 docker build -t <TAG> .`

# Credits
This project was started from the [goapp](https://github.com/alexei-led/goapp) template, a bootstrap project for Go CLI applications.
