#!/usr/bin/env bash

set -euo pipefail
export DEPLOY_STACK=${DEPLOY_STACK:-kubermatic}
export GIT_HEAD_HASH="$(git rev-parse HEAD|tr -d '\n')"
cd $(dirname $0)/../../..

source ./api/hack/lib.sh

if [[ "${DEPLOY_STACK}" == "kubermatic" ]]; then
  ./api/hack/ci/ci-push-images.sh
fi

echodate "Getting secrets from Vault"
export VAULT_ADDR=https://vault.loodse.com/
retry 5 vault write \
  --format=json auth/approle/login \
  role_id=${VAULT_ROLE_ID} secret_id=${VAULT_SECRET_ID} > /tmp/vault-token-response.json
export VAULT_TOKEN="$(cat /tmp/vault-token-response.json| jq .auth.client_token -r)"
export KUBECONFIG=/tmp/kubeconfig
export VALUES_FILE=/tmp/values.yaml

vault kv get -field=kubeconfig dev/seed-clusters/ci.kubermatic.io > ${KUBECONFIG}
vault kv get -field=values.yaml dev/seed-clusters/ci.kubermatic.io > ${VALUES_FILE}
echodate "Successfully got secrets for dev from Vault"

echodate "Deploying ${DEPLOY_STACK} stack to ci.kubermatic.io"
TILLER_NAMESPACE=kubermatic \
	DEPLOY_NODEPORT_PROXY=false \
	DEPLOY_ALERTMANAGER=false \
	DEPLOY_MINIO=false ./api/hack/deploy.sh master ${VALUES_FILE}
echodate "Successfully deployed ${DEPLOY_STACK} stack to ci.kubermatic.io"
