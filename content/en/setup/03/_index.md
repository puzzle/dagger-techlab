---
title: "Installation for Windows"
weight: 3
type: docs
sectionnumber: 1
---

## Installation for Windows

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


## Verification

Verify the installation:

```bash
where.exe dagger
```

This should output something similar to:

```
C:\<your home folder>\dagger\dagger.exe
```


## Update

To update your `dagger` CLI to the latest version, use the same command as for a fresh installation, i.e.:

```bash
Invoke-WebRequest -UseBasicParsing -Uri https://dl.dagger.io/dagger/install.ps1 | Invoke-Expression; Install-Dagger
```

