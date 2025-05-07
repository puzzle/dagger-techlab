#!/bin/bash

set -euo pipefail

: "${TRAINING_HELM_RELEASE:?Missing release name}"
: "${TRAINING_NAMESPACE:?Missing namespace}"
: "${KUBECONFIG:?Missing KUBECONFIG path}"

echo "⛏️  Uninstalling Helm release: $TRAINING_HELM_RELEASE from namespace: $TRAINING_NAMESPACE"

helm uninstall "$TRAINING_HELM_RELEASE" \
  --namespace "$TRAINING_NAMESPACE" \
  --kubeconfig "$KUBECONFIG" \
  --ignore-not-found

echo "✅ Successfully cleaned up Helm release."