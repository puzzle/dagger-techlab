---
title: "Installation for Linux"
weight: 1
type: docs
sectionnumber: 1
---

## Installation for Linux

Download and install the latest `dagger` CLI version:

```bash
curl -L https://dl.dagger.io/dagger/install.sh | BIN_DIR=$HOME/.local/bin sh
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
~/.local/bin/dagger
dagger v0.12.7 (registry.dagger.io/engine) linux/amd64
```


## Update

To update your `dagger` CLI to the latest version, use the same command as for a fresh installation, i.e.:

```bash
curl -L https://dl.dagger.io/dagger/install.sh | BIN_DIR=$HOME/.local/bin sh
```

