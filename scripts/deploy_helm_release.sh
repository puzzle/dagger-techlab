
#!/bin/bash

set -euo pipefail

: "${HELM_RELEASE:?HELM_RELEASE is required}"
: "${NAMESPACE:?NAMESPACE is required}"
: "${TRAINING_VERSION:?TRAINING_VERSION is required}"

echo "🚀 Deploying Helm release '$HELM_RELEASE' into namespace '$NAMESPACE'..."

helm upgrade "$HELM_RELEASE" acend-training-chart \
  --install \
  --wait \
  --kubeconfig "$HOME/.kube/config" \
  --namespace "$NAMESPACE" \
  --set=app.name="$HELM_RELEASE" \
  --set=app.version="$TRAINING_VERSION" \
  --repo=https://acend.github.io/helm-charts/ \
  --values=helm-chart/values.yaml \
  --set-string=acendTraining.deployments[0].ingress.labels.public=true \
  --atomic

echo "✅ Helm release '$HELM_RELEASE' deployed successfully."