---
title: "2.3 The Dagger Shell"
weight: 23
sectionnumber: 2.3
description: >
  A fast way to use Dagger.
---

## The Dagger Shell

The Dagger shell is an alternative to the `dagger call` command of the [Dagger CLI](../2.1/_index.md). An other way to interact with the Dagger Engine.

It is a simpler setup to get started to use Dagger. There is no need to initialize modules or to use a Dagger SDK.

This interpreter looks like a shell pipeline: `command, | (pipe), command, | (pipe), ...`

Under the hood it calls the Dagger API to run functions. It is based on a typed and discoverable API.

The Dagger Shell builds on top of the full Dagger power:

* Caching
* Module support
* Dynamic API extension
* ...


### Usage Example

```sh
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
