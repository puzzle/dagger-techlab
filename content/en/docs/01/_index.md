---
title: "1. Functions and Chaining "
weight: 1
sectionnumber: 1
---

## 1. Functions and Chaining

### Function Calls from the CLI

The most common way to call Dagger Functions is using the `dagger` CLI:

```bash
dagger -m github.com/shykes/daggerverse/hello@v0.3.0 call hello
```

The `dagger` CLI is first loads a `hello` module directly from its [GitHub repository](https://github.com/shykes/daggerverse/tree/main/hello) and then executes the `Hello()` function from that module.

After a while you should see:

```
hello, world!
```

{{% alert title="Note" color="primary" %}}
Due to Daggers caching mechanism, subsequent calls will be executed much faster!
{{% /alert %}}


### Exploring Modules and Functions

If you are curious, what other functions are available on this module, you can either have a look at its [source code](https://github.com/shykes/daggerverse/blob/main/hello/main.go) 
or you can explore its functions using:

```bash
dagger -m github.com/shykes/daggerverse/hello@v0.3.0 functions
```

In this particular case, there aren't any other functions :( - but what about additional arguments of the `Hello()` function?
Let's find out:

```bash
dagger -m github.com/shykes/daggerverse/hello@v0.3.0 call hello --help
```

{{% alert title="Note" color="primary" %}}
Additional to the available arguments, this often also shows you the type of value a particular argument expects.
{{% /alert %}}


### Function Arguments

Dagger Functions can accept arguments. In addition to basic types (string, boolean, integer, array...), 
Dagger also defines powerful core types: Directory, File, Container, Service, and Secret.

#### String Arguments

To pass a string argument to a Dagger Function, append the corresponding flag to the dagger call command, followed by the string value:

```bash
dagger -m github.com/shykes/daggerverse/hello@v0.3.0 call hello --name=sun
```

#### Boolean Arguments

To pass a boolean argument to a Dagger Function, simply add the corresponding flag:
    
- To set the argument to true: `--foo=true`, or simply `--foo`
- To set the argument to false: `--foo=false`

#### Directory Arguments
You can also pass a directory argument. To do so, add the corresponding flag, followed by a local filesystem path 
**or** a remote Git reference. 

In **both** cases, the `dagger` CLI will convert it to an object referencing the contents of that filesystem path or Git repository location, 
and pass the resulting `Directory` object as argument to the Dagger Function.

#### Container Arguments

Same as directories, you can pass a container argument. To do so, add the corresponding flag, followed by the address of an OCI image. 

The CLI will dynamically pull the image, and pass the resulting `Container` object as argument to the Dagger Function.

```bash
dagger -m github.com/jpadams/daggerverse/trivy@v0.4.0 call scan-container --ctr=alpine:latest
```
{{% alert title="Note" color="primary" %}}
It is important to know that in Dagger, a `Container` object is not merely a string referencing an image on a remote registry. 
It is the **actual state of a container**, managed by the Dagger Engine, and passed to a Dagger Function's code as if it were just another variable!
{{% /alert %}}

#### Secret Arguments

Dagger allows you to use confidential information, such as passwords, tokens, etc., in your Dagger Functions, without exposing them in plaintext logs, 
writing them into the filesystem of containers you're building, or inserting them into the cache.

To pass a secret to a Dagger Function, source it from a host environment variable `env:`, the host filesystem `file:`, or a host command `cmd:`.

Here is an example of passing a GitHub access token from an environment variable named `GITHUB_TOKEN` to a Dagger Function. 
The Dagger Function uses the token to query the GitHub CLI for a list of issues in the Dagger open-source repository:

```bash
dagger -m github.com/aweris/daggerverse/gh@99a1336f8091ff43bf833778a324de1cadcf25ac call run --token=env:GITHUB_TOKEN --cmd="issue list --repo=dagger/dagger"
```


### Task 1.1: Make use of arguments

Call the `Hello()` function so that it returns the string `Welcome, sunshine!` in ASCII-art.

```bash
dagger -m github.com/shykes/daggerverse/hello@v0.3.0 call hello --giant --greeting=Welcome --name=sunshine
```

