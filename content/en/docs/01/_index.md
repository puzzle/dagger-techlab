---
title: "1. Function Calls from the CLI"
weight: 1
sectionnumber: 1
---

## The Dagger CLI

The Dagger CLI is the connection to your [Dagger Engine](https://dagger.io/dagger-engine).

On a common setup the Dagger CLI manages the Dagger Engine over the Docker API.

The Dagger CLI lets you create, inspect and invoke Dagger Functions from the command line.


## Dagger CLI commands

Once installed, the `dagger` CLI offers you these commands:

```
DAGGER CLOUD COMMANDS
  login         Log in to Dagger Cloud
  logout        Log out from Dagger Cloud

DAGGER MODULE COMMANDS
  call          Call one or more functions, interconnected into a pipeline
  config        Get or set module configuration
  core          Call a core function
  develop       Prepare a local module for development
  functions     List available functions
  init          Initialize a new module
  install       Install a dependency

EXECUTION COMMANDS
  query         Send API queries to a dagger engine
  run           Run a command in a Dagger session

ADDITIONAL COMMANDS
  completion    Generate the autocompletion script for the specified shell
  help          Help about any command
  version       Print dagger version
```

{{% alert title="Note" color="primary" %}}
Checkout the autocompletion by tipping `dagger`, followed by some `Tab` keystrokes.\
Or visit the official documentation: https://docs.dagger.io/reference/cli/
{{% /alert %}}


## Function Calls from the CLI

The easiest way to get to a function is to use a Dagger module.

{{% alert title="Note" color="primary" %}}
Dagger Functions are packaged, shared and reused using Dagger Modules.
{{% /alert %}}

Dagger Modules are published at the [Daggerverse](https://daggerverse.dev/).
It is similar to the [MvnRepository](https://mvnrepository.com/). The MvnRepository provides Java libraries and the Daggerverse provides Dagger Modules.

The most common way to call Dagger Functions is using the `dagger` CLI:

```bash
dagger call --mod ./mod hello
```

The `dagger` CLI first loads a `hello` module directly from its [GitHub repository](https://github.com/shykes/daggerverse/tree/main/hello) and then executes the `Hello()` function from that module.

{{% alert title="Note" color="primary" %}}
Explanation to the dagger CLI call.\
`dagger call` : execute the dagger CLI `call` command\
`--mod ./mod` : `call` command option to use the `hello` module (load its functions)\
`hello` : execute the `hello` function
{{% /alert %}}

After a while you should see:

```
hello, world!
```

{{% alert title="Note" color="primary" %}}
Due to Daggers caching mechanism, subsequent calls will be executed much faster!
{{% /alert %}}


### Exploring Modules and Functions

If you are curious, what other [Functions](https://docs.dagger.io/api/reference/#definition-Function) are available on this [Module](https://docs.dagger.io/api/reference/#definition-Module), you can either have a look at its [source code](https://github.com/shykes/daggerverse/blob/main/hello/main.go)
or you can explore its functions using:

```bash
dagger functions --mod ./mod
```

In this particular case, there aren't any other functions :( - but what about additional arguments of the `Hello()` function?
Let's find out:

```bash
dagger call --mod ./mod hello --help
```

{{% alert title="Note" color="primary" %}}
Additional to the available arguments, this often also shows you the type of value a particular argument expects.
{{% /alert %}}


### Function Arguments

Dagger Functions can accept arguments. In addition to basic types (string, boolean, integer, array...),
Dagger also defines powerful core types: [Container](https://docs.dagger.io/api/reference/#definition-Container), [Directory](https://docs.dagger.io/api/reference/#definition-Directory), [File](https://docs.dagger.io/api/reference/#definition-File), [Service](https://docs.dagger.io/api/reference/#definition-Service) and [Secret](https://docs.dagger.io/api/reference/#definition-Secret).


#### String Arguments

To pass a String argument to a Dagger Function, append the corresponding flag to the dagger call command, followed by the string value:

```bash
dagger call --mod ./mod hello --name=sun
```


#### Boolean Arguments

To pass a Boolean argument to a Dagger Function, simply add the corresponding flag:

* To set the argument to true: `--foo=true`, or simply `--foo`
* To set the argument to false: `--foo=false`, or just omit the argument as the default is usually false

True:

```bash
# explicit
dagger call --mod ./mod hello --shout=true
```

```bash
# implicit
dagger call --mod ./mod hello --shout
```

False:

```bash
# explicit
dagger call --mod ./mod hello --shout=false
```

```bash
# implicit
dagger call --mod ./mod hello
```


#### Directory Arguments

You can also pass a Directory argument. To do so, add the corresponding flag, followed by a local filesystem path **or** a remote Git reference.

In **both** cases, the `dagger` CLI will convert it to an object referencing the contents of that filesystem path or Git repository location,
and pass the resulting `Directory` object as argument to the Dagger Function.

Filesystem path:
```bash
dagger call --mod ./mod ls --dir .
```

Git repository:
```bash
dagger call --mod ./mod ls --dir https://github.com/puzzle/action-owasp-dependecy-track-check
```


#### Container Arguments

Same as directories, you can pass a Container argument. To do so, add the corresponding flag, followed by the address of an OCI image.

The CLI will dynamically pull the image, and pass the resulting `Container` object as argument to the Dagger Function.

```bash
dagger --mod ./mod call os --ctr=alpine:latest
```

{{% alert title="Note" color="primary" %}}
It is important to know that in Dagger, a `Container` object is not merely a string referencing an image on a remote registry.
It is the **actual state of a container**, managed by the Dagger Engine, and passed to a Dagger Function's code as if it were just another variable!
{{% /alert %}}


#### Secret Arguments

Dagger allows you to use confidential information, such as passwords, tokens, etc., in your Dagger Functions, without exposing them in plaintext logs,
writing them into the filesystem of containers you're building, or inserting them into the cache.

To pass a Secret to a Dagger Function, source it from a host environment variable `env:`, the host filesystem `file:`, or a host command `cmd:`.

Here is an example of passing a GitHub access token from an environment variable named `GITHUB_TOKEN` to a Dagger Function.

The Dagger Function uses the token to query the GitHub CLI for a list of issues in the Dagger open-source repository:

```bash
dagger call --mod github.com/aweris/daggerverse/gh@v0.0.2 run --token=env:GITHUB_TOKEN --cmd="issue list --repo=dagger/dagger"
```

{{% alert title="Note" color="primary" %}}
This is only an example, you don't have to make it run.
{{% /alert %}}


## Task {{% param sectionnumber %}}.1: Explore a module

Explore the `github.com/purpleclay/daggerverse/ponysay@v0.1.0` module.

Make it return the phrase `Dagger puts a smile on my face!`.

{{% details title="show hint" mode-switcher="normalexpertmode" %}}
```bash
dagger functions --mod github.com/purpleclay/daggerverse/ponysay@v0.1.0
```
{{% /details %}}

{{% details title="show hint" mode-switcher="normalexpertmode" %}}
```bash
dagger call --mod github.com/purpleclay/daggerverse/ponysay@v0.1.0 say --help
```
{{% /details %}}

{{% details title="show solution" mode-switcher="normalexpertmode" %}}
```bash
dagger call --mod github.com/purpleclay/daggerverse/ponysay@v0.1.0 say --msg="Dagger puts a smile on my face!"
```
{{% /details %}}


## Task {{% param sectionnumber %}}.2: Make use of multiple arguments

Call the `Hello()` function of `github.com/shykes/daggerverse/hello@v0.3.0` so that it returns the phrase `Welcome, sunshine!` in ASCII-art (giant letters).

{{% details title="show solution" mode-switcher="normalexpertmode" %}}
```bash
dagger call --mod ./mod hello --giant --greeting=Welcome --name=sunshine
```
{{% /details %}}


## Task {{% param sectionnumber %}}.3: Pass a secret

Set the `--password` value in the following call with a secret, using an environment variable, containing the password "MySuperSecret".

```bash
dagger call --mod ./mod unlock --password=visible
```

{{% details title="show solution" mode-switcher="normalexpertmode" %}}
```bash
export SECRET=MySuperSecret
dagger call --mod ./mod unlock --password env:SECRET
```
{{% /details %}}

or using a file

{{% details title="show solution" mode-switcher="normalexpertmode" %}}
```bash
echo $SECRET > secret.txt
dagger call --mod ./mod unlock --password file:./secret.txt
```
{{% /details %}}

or using a command

{{% details title="show solution" mode-switcher="normalexpertmode" %}}
```bash
dagger call --mod ./mod unlock --password cmd:"echo $SECRET"
```
{{% /details %}}

{{% alert title="Note" color="primary" %}}
Unless you provide the right password, the function execution will fail with an error.
{{% /alert %}}
