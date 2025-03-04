---
title: "Setup"
weight: 1
type: docs
menu:
  main:
    weight: 1
---

## Installation Guide

This guide shows you how to install the `dagger` CLI.

Follow the instructions on the subsequent pages to complete the setup on your operating system of choice.


### Prerequisits

In order to use Dagger you need a few prerequisits. You need an text editor to write and update your Dagger Pipeline. Then you need a container tool such as Docker, which Dagger calls on oin the background, and then of course the Dagger CLI.

* **Editor** You need a text editor such as vi, vi, gedit, nano, ...
* **Container Tool** The dagger engine is run in a container tool, we highly recommend Docker, as it needs docker compose and a rootless setup (so it might run bugy with Podman).
* **Dagger CLI** You have to install the Dagger-CLI on your machine, see below how to achieve this.

In the following you can find installation guides for [Linux](#installation-for-linux), [MacOS](installation-for-macos) and [Windows](installation-for-windows). 


### Installation for Linux

Download and install the latest `dagger` CLI version:

```bash
curl -L https://dl.dagger.io/dagger/install.sh | BIN_DIR=$HOME/.local/bin sh
```

{{% alert title="Note" color="primary" %}}
To install a different `dagger` CLI version, you can specify it using this param, added to the previous command:

`DAGGER_VERSION=x.y.z`
{{% /alert %}}


#### Verification

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


#### Update

To update your `dagger` CLI to the latest version, use the same command as for a fresh installation, i.e.:

```bash
curl -L https://dl.dagger.io/dagger/install.sh | BIN_DIR=$HOME/.local/bin sh
```


### Installation for macOS

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


#### Verification

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


#### Update

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


### Installation for Windows

The `dagger` CLI can be installed on Windows using a PowerShell 7.0 script:

```bash
Invoke-WebRequest -UseBasicParsing -Uri https://dl.dagger.io/dagger/install.ps1 | Invoke-Expression; Install-Dagger
```

{{% alert title="Note" color="primary" %}}
To install a different `dagger` CLI version, you can specify it using this param, added to the previous command:

`-DaggerVersion x.y.z`
{{% /alert %}}

For further customizations, such as adding Dagger to your system's PATH or using the interactive installation process,
additional parameters are available. To view all available options:

```bash
Invoke-WebRequest -UseBasicParsing -Uri https://dl.dagger.io/dagger/install.ps1 | Invoke-Expression;
Get-Command -Name Install-Dagger -Syntax
```


#### Verification

Verify the installation:

```bash
where.exe dagger
```

This should output something similar to:

```
C:\<your home folder>\dagger\dagger.exe
```


#### Update

To update your `dagger` CLI to the latest version, use the same command as for a fresh installation, i.e.:

```bash
Invoke-WebRequest -UseBasicParsing -Uri https://dl.dagger.io/dagger/install.ps1 | Invoke-Expression; Install-Dagger
```



