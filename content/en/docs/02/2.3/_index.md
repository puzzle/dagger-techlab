---
title: "2.3 The Dagger Shell"
weight: 23
sectionnumber: 2.3
description: >
  A fast way to use Dagger.
---

## The Dagger Shell

The Dagger shell is an alternative to the `dagger call` command of the [Dagger CLI](../2.1/). An other way to interact with the Dagger Engine.

It is a simpler setup to get started to use Dagger. There is no need to initialize modules or to use a Dagger SDK.

This interpreter looks like a shell pipeline (with commands and pipes `|`): `command | command | ...`

Under the hood it calls the Dagger API to run functions. It is based on a typed and discoverable API.

The Dagger Shell builds on top of the full Dagger power:

* Caching
* Module support
* Dynamic API extension
* ...


### Usage Example

```bash
⋈ container | from alpine | with-exec echo "Daggernaut" | stdout
Daggernaut
```

This is what happens:

1. call a Dagger core function: `container`
1. get an immutable artifact: container
1. call a container function: `from alpine`
1. get an immutable artifact: container with alpine as base image
1. call a container function: `with-exec`
1. get an immutable artifact: container with executed command
1. call a container function: `stdout`
1. get an immutable artifact: string (output of the container command)

The string output of the executed command `echo "Daggernaut"` is returned and printed to the shell.

Basically this builds up the DAG (directed acyclic graph) which will run on Dagger, the **open-source runtime for composable workflows**.


## Task {{% param sectionnumber %}}.1: Start the shell

Open your terminal at any location and execute the `dagger shell` command of the [Dagger CLI](../2.1/).

The output should look like this:

```bash
$ dagger shell
Experimental Dagger interactive shell. Type ".help" for more information. Press Ctrl+D to exit.
⋈
```

The special prompt `⋈` is ready to get and execute commands.


## Task {{% param sectionnumber %}}.2: Explore the shell functionality

Run the `.help` command. This will show you the available commands.

{{% details title="show hint" mode-switcher="normalexpertmode" %}}
```bash
.help
```
{{% /details %}}

Try to get the help of the `.install` command.

{{% details title="show hint" mode-switcher="normalexpertmode" %}}
```bash
.help .install
```
{{% /details %}}

You can get even more details using the `help` option:

```bash
.install --help
```


## Task {{% param sectionnumber %}}.3: Use a Dagger module

The Dagger Shell session can have a default module. This is set with the `.use` command.

Set the [Puzzle cosign](https://daggerverse.dev/mod/github.com/puzzle/dagger-module-cosign/cosign@v0.1.1) module as default.

```bash
⋈ .use github.com/puzzle/dagger-module-cosign/cosign@v0.1.1
```

The prompt shows now the name of the default module:

```bash
github.com/puzzle/dagger-module-cosign/cosign@v0.1.1 ⋈
```

See the extended API by running `.help` again:

```bash
github.com/puzzle/dagger-module-cosign/cosign@v0.1.1 ⋈ .help
BUILTIN COMMANDS
  .core         Load any core Dagger type
  ...

AVAILABLE MODULE FUNCTIONS
  sign             Sign will run cosign sign from the image, as defined by the cosignImage
  sign-keyless     SignKeyless will run cosign sign (keyless) from the image, as defined by the cosignImage
  attest           Attest will run cosign attest from the image, as defined by the cosignImage
  attest-keyless   AttestKeyless will run cosign attest (keyless) from the image, as defined by the cosignImage
  clean            Clean will run cosign clean from the image, as defined by the cosignImage

STANDARD COMMANDS
  cache-volume   Constructs a cache volume for a given cache key.
  ...
```


### Explore the sign function

Explore the `sign` function of the default module.

Find out, which arguments are required.

{{% details title="show result" mode-switcher="normalexpertmode" %}}
The required arguments are:

* private-key
* password
* digest
{{% /details %}}

{{% details title="show hint" mode-switcher="normalexpertmode" %}}
Run the help command:
```bash
.help sign
```
{{% /details %}}
