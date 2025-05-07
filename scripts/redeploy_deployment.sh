#!/bin/bash

set -euo pipefail

: "${HELM_RELEASE:?HELM_RELEASE is required}"
: "${HELM_NAME:?HELM_NAME is required}"
: "${NAMESPACE:?NAMESPACE is required}"

DEPLOYMENT_NAME="${HELM_RELEASE}-${HELM_NAME}"

echo "♻️ Redeploying deployment '$DEPLOYMENT_NAME' in namespace '$NAMESPACE'..."

kubectl rollout restart deployment/"$DEPLOYMENT_NAME" \
  --kubeconfig "$HOME/.kube/config" \
  --namespace "$NAMESPACE"

echo "✅ Deployment '$DEPLOYMENT_NAME' restarted successfully."
