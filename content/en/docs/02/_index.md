---
title: "2. Functions Chaining"
weight: 2
sectionnumber: 2
---

## Functions Chaining

Dagger Functions can return either basic types or objects. Objects can define their own functions.\
So when calling a Dagger Function that returns an object, the Dagger API lets you follow up by calling one of that object's
functions, which itself can return another object, and so on.\
This is called "function chaining", and is a core feature of Dagger.

Dagger's core types ([Container](https://docs.dagger.io/api/reference/#definition-Container), [Directory](https://docs.dagger.io/api/reference/#definition-Directory), [File](https://docs.dagger.io/api/reference/#definition-File), [Service](https://docs.dagger.io/api/reference/#definition-Service), ...)
are all objects. They each define various functions for interacting with their respective objects.

Let's explore them step by step:

```bash
dagger call --mod github.com/sagikazarmark/daggerverse/openssh-server@v0.1.0 --help
```

{{% details title="show available 'module' functions" mode-switcher="normalexpertmode" %}}
```
    USAGE                                                                            
      dagger call [options] [arguments] <function>

    FUNCTIONS                           
      container             
      host-keys             Returns the SSH host keys.
      known-hosts           Return a formatted SSH known_hosts file.              
 ---> service               Return a service that runs the OpenSSH server.
      with-authorized-key   Authorize a public key.
      with-config           Mount a custom SSH configuration file (with .conf extension). 
```
{{% /details %}}

```bash
dagger call --mod github.com/sagikazarmark/daggerverse/openssh-server@v0.1.0 service --help
```

{{% details title="show available 'service' object functions" mode-switcher="normalexpertmode" %}}
```
    USAGE
      dagger call service [arguments] <function>

    FUNCTIONS
      endpoint      Retrieves an endpoint that clients can use to reach this container.
      hostname      Retrieves a hostname which can be used by clients to reach this container.
      ports         Retrieves the list of ports provided by the service.
      start         Start the service and wait for its health checks to succeed.
      stop          Stop the service.
 ---> up            Creates a tunnel that forwards traffic from the callers network to this service.
```
{{% /details %}}

```bash
dagger call --mod github.com/sagikazarmark/daggerverse/openssh-server@v0.1.0 service up --help
```

{{% details title="show available 'up' function arguments" mode-switcher="normalexpertmode" %}}
```
    USAGE
      dagger call service up [arguments]

    ARGUMENTS
 ---> --ports PortForward   List of frontend/backend port mappings to forward.
                            Frontend is the port accepting traffic on the host, backend is the service port. (default [])
      --random              Bind each tunnel port to a random port on the host.

```
{{% /details %}}

Now that we have got all the pieces together, let's expose a Service returned by a Dagger Function on a specified host port,
by chaining a call to the `Service` object's `Up()` function:

```bash
dagger call --mod github.com/sagikazarmark/daggerverse/openssh-server@v0.1.0 service up --ports=22022:22
```

Here we print the contents of a File returned by a Dagger Function, by chaining a call to the `File` object's `Contents()` function:

```bash
dagger call --mod github.com/dagger/dagger/dev/ruff@a29dadbb5d9968784847a15fccc5629daf2985ae lint --source=https://github.com/puzzle/puzzle-radicale-auth-ldap report contents
```


### Task {{% param sectionnumber %}}.1: Chain calls

Display and return the contents of the `/etc/os-release` file in a container, by chaining additional calls to the `Container`
object of the `github.com/shykes/daggerverse/wolfi@v0.1.3` module:

{{% details title="show hint" mode-switcher="normalexpertmode" %}}
Have a look at the [WithExec()](https://docs.dagger.io/api/reference/#Container-withExec) and [Stout()](https://docs.dagger.io/api/reference/#Container-stdout) functions.
{{% /details %}}

{{% details title="show solution" mode-switcher="normalexpertmode" %}}
```bash
dagger call --mod github.com/shykes/daggerverse/wolfi@v0.1.3 container with-exec --args="cat","/etc/os-release" stdout
```
{{% /details %}}

Try an alternative approach using [File()](https://docs.dagger.io/api/reference/#definition-File) instead.

{{% details title="show hint" mode-switcher="normalexpertmode" %}}
```bash
dagger call --mod github.com/shykes/daggerverse/wolfi@v0.1.3 container file --help
```
{{% /details %}}

{{% details title="show solution" mode-switcher="normalexpertmode" %}}
```bash
dagger call --mod github.com/shykes/daggerverse/wolfi@v0.1.3 container file --path=/etc/os-release contents
```
{{% /details %}}
