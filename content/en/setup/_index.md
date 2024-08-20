---
title: "Setup"
weight: 1
type: docs
menu:
  main:
    weight: 1
---

## Install the Dagger CLI

### Linux

Download and install the latest `dagger` version:

```bash
curl -L https://dl.dagger.io/dagger/install.sh | BIN_DIR=$HOME/.local/bin sh
```

Verify:

```bash
which dagger
~/.local/bin/dagger

dagger version
dagger v0.12.3 (registry.dagger.io/engine) linux/amd64
```

### macOS

Install the latest `dagger` version using `homebrew`:

```bash
brew install dagger/tap/dagger
```

Or using `sh` and `curl`:

```bash
cd /usr/local
curl -L https://dl.dagger.io/dagger/install.sh | sh
```

Verify:

```bash
which dagger
/opt/homebrew/bin/dagger

dagger version
dagger v0.12.3 (registry.dagger.io/engine:v0.12.3) darwin/arm64
```




