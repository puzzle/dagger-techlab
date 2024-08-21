---
title: "Installation for macOS"
weight: 2
type: docs
sectionnumber: 1
---

## Installation for macOS

Install the latest `dagger` version using `homebrew`:

```bash
brew install dagger/tap/dagger
```

Or using `sh` and `curl`:

```bash
cd /usr/local
curl -L https://dl.dagger.io/dagger/install.sh | sh
```

## Verification

Verify the installation:

```bash
which dagger
dagger version
```

This should output something similar to:

```
/opt/homebrew/bin/dagger
dagger v0.12.5 (registry.dagger.io/engine:v0.12.5) darwin/arm64
```


## Update

To update your `dagger` CLI to the latest version using `homebrew`, use this commands:

```bash
brew update
brew upgrade dagger
```

To update using `sh` and `curl` instead, use the same commands as for a fresh installation, i.e.:

```bash
cd /usr/local
curl -L https://dl.dagger.io/dagger/install.sh | sh
```