AZURE_BASE_GROUP_NAME=
AZURE_LOCATION_DEFAULT=westus2

# create with:
# `az ad sp create-for-rbac --name 'my-sp' --output json --sdk-auth`
# sp must have Contributor role on subscription
AZURE_TENANT_ID=
AZURE_CLIENT_ID=
AZURE_CLIENT_SECRET=
AZURE_SUBSCRIPTION_ID=

# create with:
# `az ad sp create-for-rbac --name 'my-sp' --sdk-auth > $HOME/.azure/sdk_auth.json`
# sp must have Contributor role on subscription
AZURE_AUTH_LOCATION=$HOME/.azure/sdk_auth.json

AZURE_STORAGE_ACCOUNT_NAME=
AZURE_STORAGE_ACCOUNT_GROUP_NAME=
