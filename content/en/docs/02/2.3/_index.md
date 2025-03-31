---
title: "2.3 The Dagger Shell"
weight: 23
sectionnumber: 2.3
description: >
  A fast way to use Dagger.
---

## The Dagger Shell

The Dagger shell is an alternative to the `dagger call` command of the [Dagger CLI](../2.1/). An other way to interact with the Dagger Engine.

Not to be confused with this dagger shell:

![the other dagger shell](dagger-shell.png)

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

Open your terminal at any location and execute the `dagger` command of the [Dagger CLI](../2.1/).

The output should look like this:

```bash
$ dagger
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

Download the [Puzzle cosign](https://daggerverse.dev/mod/github.com/puzzle/dagger-module-cosign/cosign@v0.1.1).

```bash
github.com/puzzle/dagger-module-cosign/cosign@v0.1.1
```

See the extended API by running `.help` again:

```bash
⋈ github.com/puzzle/dagger-module-cosign/cosign@v0.1.1 | .help
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
github.com/puzzle/dagger-module-cosign/cosign@v0.1.1 | .help sign
```
{{% /details %}}


## Task {{% param sectionnumber %}}.4: Define a Build

With the knowledge of the Dagger Shell we go to define a Container build.

Make sure that there is no default module set.
Exit the shell with `Ctrl +d` and start a new shell.

{{% details title="show hint" mode-switcher="normalexpertmode" %}}
```bash
dagger
```
{{% /details %}}

Start with the standard command to get a container:

```bash
container
```

This will return a container entity of the Dagger core API.

```bash
defaultArgs: []
entrypoint: []
mounts: []
platform: linux/amd64
user: ""
workdir: ""
```

Run help to get the available functions on the container object.

{{% details title="show hint" mode-switcher="normalexpertmode" %}}
```bash
container | .help
```
{{% /details %}}

There we find the `from` function. It is known from the Dockerfile.

The base image of the container should be `alpine`.
Use the `from` function to extend the container with alpine.

What will be returned?

{{% details title="show hint" mode-switcher="normalexpertmode" %}}
```bash
container | from alpine
```

The new container (state) is returned.
{{% /details %}}

The container should run the `echo "Daggernaut"` command.

This will be achieved using the `with-exec`.

{{% details title="show hint" mode-switcher="normalexpertmode" %}}
```bash
container | from alpine | with-exec echo "Daggernaut"
```
{{% /details %}}

The updated container is returned.

To get the output of the `echo` command inside the container, the Dagger container `stdout` function is needed.

Extend the build with the `stdout` function.


{{% details title="show hint" mode-switcher="normalexpertmode" %}}
```bash
container | from alpine | with-exec echo "Daggernaut" | stdout
```
{{% /details %}}

There is our build chain!

* it constructed a DAG in the background
* the executions are cached, such that only the newly added functions have to be executed

Create an artefact, add a transformation, get a new artefact, add a transformation, ...

<!---
TODO:
* [ ] variables, sub-shells
* [ ] tab tab autocompletion
* [ ] replay examples from 1.x labs
* [ ] new API: with-xyz
* [ ] replace the first example, such that it is not the same as the Task 4
-->
