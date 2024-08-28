---
title: "1. Function Calls from the CLI"
weight: 1
sectionnumber: 1
---

## {{% param sectionnumber %}}. Function Calls from the CLI


Once installed, the `dagger` CLI offers you these functions:

```
call        (Call one or more functions, interconnected into a pipeline)
completion  (Generate the autocompletion script for the specified shell)
config      (Get or set module configuration)
core        (Call a core function)
develop     (Prepare a local module for development)
functions   (List available functions)
help        (Help about any command)
init        (Initialize a new module)
install     (Install a dependency)
login       (Log in to Dagger Cloud)
logout      (Log out from Dagger Cloud)
query       (Send API queries to a dagger engine)
run         (Run a command in a Dagger session)
version     (Print dagger version)
```

{{% alert title="Note" color="primary" %}}
Checkout the autocompletion by tipping `dagger`, followed by some `Tab` keystrokes.
{{% /alert %}}

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
Dagger also defines powerful core types: [Container](https://docs.dagger.io/api/reference/#definition-Container), [Directory](https://docs.dagger.io/api/reference/#definition-Directory), [File](https://docs.dagger.io/api/reference/#definition-File), [Service](https://docs.dagger.io/api/reference/#definition-Service) and [Secret](https://docs.dagger.io/api/reference/#definition-Secret).


#### String Arguments

To pass a String argument to a Dagger Function, append the corresponding flag to the dagger call command, followed by the string value:

```bash
dagger -m github.com/shykes/daggerverse/hello@v0.3.0 call hello --name=sun
```


#### Boolean Arguments

To pass a Boolean argument to a Dagger Function, simply add the corresponding flag:

* To set the argument to true: `--foo=true`, or simply `--foo`
* To set the argument to false: `--foo=false`


#### Directory Arguments

You can also pass a Directory argument. To do so, add the corresponding flag, followed by a local filesystem path **or** a remote Git reference.

In **both** cases, the `dagger` CLI will convert it to an object referencing the contents of that filesystem path or Git repository location,
and pass the resulting `Directory` object as argument to the Dagger Function.


#### Container Arguments

Same as directories, you can pass a Container argument. To do so, add the corresponding flag, followed by the address of an OCI image.

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

To pass a Secret to a Dagger Function, source it from a host environment variable `env:`, the host filesystem `file:`, or a host command `cmd:`.

Here is an example of passing a GitHub access token from an environment variable named `GITHUB_TOKEN` to a Dagger Function.
The Dagger Function uses the token to query the GitHub CLI for a list of issues in the Dagger open-source repository:

```bash
dagger -m github.com/aweris/daggerverse/gh@v0.0.2 call run --token=env:GITHUB_TOKEN --cmd="issue list --repo=dagger/dagger"
```


### Task {{% param sectionnumber %}}.1: Explore a module

Explore the `github.com/purpleclay/daggerverse/ponysay@v0.1.0` module.
Make it return the phrase `Dagger puts a smile on my face!`.

{{% details title="show hint" mode-switcher="normalexpertmode" %}}
```bash
dagger -m github.com/purpleclay/daggerverse/ponysay@v0.1.0 functions
```
{{% /details %}}

{{% details title="show hint" mode-switcher="normalexpertmode" %}}
```bash
dagger -m github.com/purpleclay/daggerverse/ponysay@v0.1.0 call say --help
```
{{% /details %}}

{{% details title="show solution" mode-switcher="normalexpertmode" %}}
```bash
dagger -m github.com/purpleclay/daggerverse/ponysay@v0.1.0 call say --msg="Dagger puts a smile on my face!"
```
{{% /details %}}


### Task {{% param sectionnumber %}}.2: Make use of multiple arguments

Call the `Hello()` function of `github.com/shykes/daggerverse/hello@v0.3.0` so that it returns the phrase `Welcome, sunshine!` in ASCII-art.

{{% details title="show solution" mode-switcher="normalexpertmode" %}}
```bash
dagger -m github.com/shykes/daggerverse/hello@v0.3.0 call hello --giant --greeting=Welcome --name=sunshine
```
{{% /details %}}


### Task {{% param sectionnumber %}}.3: Pass a secret

Set and replace the `--token` value in the following call with a secret using an environment variable

```bash
dagger -m github.com/aweris/daggerverse/gh@v0.0.2 call run --token=visible --cmd="issue list --repo=dagger/dagger"
```


{{% details title="show solution" mode-switcher="normalexpertmode" %}}
```bash
export SECRET=invisible
dagger -m github.com/aweris/daggerverse/gh@v0.0.2 call run --token=env:SECRET --cmd="issue list --repo=dagger/dagger"
```
{{% /details %}}

or using a file

{{% details title="show solution" mode-switcher="normalexpertmode" %}}
```bash
echo $SECRET > secret.txt
dagger -m github.com/aweris/daggerverse/gh@v0.0.2 call run --token=file:./secret.txt --cmd="issue list --repo=dagger/dagger"
```
{{% /details %}}

or using a command

{{% details title="show solution" mode-switcher="normalexpertmode" %}}
```bash
dagger -m github.com/aweris/daggerverse/gh@v0.0.2 call run --token=cmd:"head -c10 /dev/random | base64" --cmd="issue list --repo=dagger/dagger"
```
{{% /details %}}
