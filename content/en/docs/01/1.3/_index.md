---
title: "1.3 Getting started"
weight: 13
sectionnumber: 1.3
description: >
  Get familiar with the lab setup.
---

## Installation Guide

Prerequisits:

* **Container Tool** The dagger engine is run in a container tool, we highly recommend Docker, as it needs docker compose and a rootless setup (so it might run bugy with Podman). 
* **Dagger CLI** You have to install the Dagger-CLI on your machine, see below how to achieve this.


### Linux


1. Make sure, that the installation path `$HOME/.local/bin` is declared in the variable PATH
```
export PATH="$HOME/.local/bin:$PATH"
```

1. Run the following command to install Dagger
```
curl -fsSL https://dl.dagger.io/dagger/install.sh | BIN_DIR=$HOME/.local/bin sh
```

1.Check your installation by run
```
dagger version
```


### Windows


1. Make sure, that the installation path is declared in the variable PATH

1. Download the installation script and run it in the PowerShell
```
iwr https://dl.dagger.io/dagger/install.ps1 -useb | iex
```

1.Check your installation by running
```
dagger version
```


### Mac


1. Install Dagger
```
brew install dagger/tap/dagger
```
or, alternatively use an installation step
```
curl -fsSL https://dl.dagger.io/dagger/install.sh | BIN_DIR=/usr/local/bin sh
```

1.Check your installation by run
```
dagger version
```
