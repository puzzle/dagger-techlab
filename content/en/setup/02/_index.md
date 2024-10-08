---
title: "Installation for macOS"
weight: 2
type: docs
sectionnumber: 1
---

## Installation for macOS

Install the latest `dagger` CLI version using `homebrew`:

```bash
brew install dagger/tap/dagger
```

Or using `sh` and `curl`:

```bash
cd /usr/local
curl -L https://dl.dagger.io/dagger/install.sh | sh
```

{{% alert title="Note" color="primary" %}}
To install a different `dagger` CLI version, you can specify it using this param, added to the previous command:

`DAGGER_VERSION=x.y.z`
{{% /alert %}}


## Verification

Verify the installation:

```bash
which dagger
dagger version
```

This should output something similar to:

```
/opt/homebrew/bin/dagger
dagger v0.12.7 (registry.dagger.io/engine:v0.12.7) darwin/arm64
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

