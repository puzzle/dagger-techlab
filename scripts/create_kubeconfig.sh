#!/bin/bash

set -euo pipefail

: "${KUBE_CONFIG:?KUBE_CONFIG environment variable is required}"

echo "📄 Writing kubeconfig to \$HOME/.kube/config..."
mkdir -p "$HOME/.kube"
echo "$KUBE_CONFIG" > "$HOME/.kube/config"
echo "✅ Kubeconfig written successfully."
