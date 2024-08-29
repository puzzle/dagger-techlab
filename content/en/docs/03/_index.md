---
title: "3. Daggerize an App"
weight: 3
sectionnumber: 3
---


## Daggerize an App


### The Challenge

After we have learned the basic Dagger Functions, we want to apply our new knowledge to solve a real life problem:

We would like to conduct a survey regarding the popularity of the different Dagger SDKs!


### The Candidate

Fortunately, there is a free open-source quiz app called [ClassQuiz](https://classquiz.de/)\
It allows the creation of shareable, fully customizable quizzes and surveys.\
The app is split in a frontend and an api part:

* The frontend is written in type script and uses a redis memcache.
* The backend is mostly written in python, uses a postgreSQL database and meilisearch.

Caddy is used as reverse proxy to keep the parts together.


### The Journey


#### Prerequisites

Check out ClassQuiz:

```bash
git clone https://github.com/mawoka-myblock/ClassQuiz.git
```

Get familiar with the source - take a closer look at the [docker-compse.yaml](https://github.com/mawoka-myblock/ClassQuiz/blob/master/docker-compose.yml),
which is particularly interesting for our purpose.

Since the app has some hard-coded configurations that would interfere with our setup, let's apply the following [patch](config.patch) to the `classquiz/config.py`:

{{< details title="show patch" mode-switcher="normalexpertmode" >}}

{{% readAndHighlight file="config.patch" code="true" lang="patch" highlight="hl_lines=8 18-23 32" %}}

{{< /details >}}

Create first the `config.patch` file in the root of your git folder. Then add the content of the above patch.

```bash
patch classquiz/config.py < config.patch
```

{{% alert title="Note" color="primary" %}}

If patching does not work, overwrite the file `classquiz/config.py` with the content from the following `config.py` file.

{{< details title="show final config.py file" >}}

{{< readfile file="config.py" code="true" lang="Python" >}}

{{< /details >}}

{{% /alert %}}

The app also binds the privileged port `80`, which would be an obstacle as well.\
So let's replace all occurrences of `:80` in `Caddyfile-docker` with `:8081`.\
Additionally the missing protocol has to be added to the last `reverse_proxy` line. Add `http://` in front of `api:80`.

Do it by hand or use the following `sed` commands:

```bash
sed -i 's# api:80# http://api:80#g' Caddyfile-docker
sed -i 's#api:80#api:8081#g' Caddyfile-docker
```

If patching does not work, overwrite the file `Caddyfile-docker` with the content from the following `Caddyfile-docker` file.

{{< details title="show final Caddyfile-docker file" >}}

{{< readfile file="Caddyfile-docker" code="true" >}}

{{< /details >}}


#### Start using Dagger

As we learnt in the first labs, Dagger functions are needed to encapsulate our pipeline functionality.

In Dagger, everything is a Module, therefore the first step is to initialize a Dagger Module.

A new Dagger module in Go, Python or TypeScript can be initialized by running `dagger init` inside the app's root directory,
using the `--source` flag to specify a directory for the module's source code.

We will use the **Python SDK** for this example:

```bash
dagger init --sdk=python --source=./ci
```

This leaves us with a generated `dagger.json` module metadata file, an initial `ci/src/main/__init__.py` source code template, `ci/pyproject.toml` and
`ci/requirements.lock` files, as well as a generated `ci/sdk` folder for local development.
The configuration file sets the name of the module to the name of the current directory, unless an alternative is specified with the `--name` argument.


## Run the App locally

The generated `ci/src/main/__init__.py` is the starting point, which needs to be extended.\
It has already some example functions that are ready to use or extend.

The ClassQuiz repository has two Dockerfile. One to build the frontend and one to build the backend.\
A starting point is to use the Dockerfile for a Docker build.\
The resulting Docker image can be used to run the app inside a container.

As a first step, we could implement a simple `build` function:

* function name: `build`
* argument: `context` - the folder containing the Docker build context, including the Dockerfile
* return: a Dagger Container

{{% alert title="Note" color="primary" %}}
The Dagger Engine has no access to your host computer. Therefore you have to explicitly provide folders as arguments.
{{% /alert %}}

```python
    @function
    def build(self, context: dagger.Directory) -> dagger.Container:
        """Returns a container built with the given context."""
        return (
            dag.container()
            .build(context)
        )
```

The entrypoint to accessing the Dagger API from your own module's code is `dag`, the Dagger client, which is pre-initialized.\
It contains all the core types (like Container, Directory, etc.), as well as bindings to any dependencies your module has declared.

The [Python SDK Reference](https://dagger-io.readthedocs.org/) documents all Dagger API types and functions.\
Our function starts by creating a container (`dag.container()`).
[Here](https://dagger-io.readthedocs.io/en/sdk-python-v0.12.5/client.html#dagger.Container) is the reference to the Python documentation.

The [build](https://dagger-io.readthedocs.io/en/sdk-python-v0.12.5/client.html#dagger.Container.build) executes the Docker build with the given files.

This function allows us to build the frontend and expose it as a [Service](https://docs.dagger.io/manuals/developer/services) to the [localhost](https://docs.dagger.io/manuals/developer/services#expose-services-returned-by-functions-to-the-host) on port 3000:

```bash
dagger call --mod ClassQuiz build --context=./frontend/ with-exposed-port --port=3000 as-service up
```

Here we do the previous explained [Function Chaining](https://docs.dagger.io/manuals/user/chaining).

* Our `build` method returns a Dagger container.
* `with-exposed-port --port=3000` opens the port to the Container (expose)
* `as-service` Turn the container into a Service that runs the app.
* `up` Opens the connection to the app Service. (Creates a tunnel that forwards traffic from the caller’s network to this service.)

{{% alert title="Note" color="primary" %}}

As we are in the root directory of the Dagger module, we do not need to provide the module (`--mod`) option.\
This is the simplified command:

```bash
dagger call build --context=./frontend/ with-exposed-port --port=3000 as-service up
```
{{% /alert %}}

Use `Ctrl +c` to stop the container.

And the backend as well with its context folder:

```bash
dagger call build --context=. with-exposed-port --port=8000 as-service up
```

Use `Ctrl +c` to stop the container.

But if we have a closer look to the console output, we will discover some error messages due to missing configurations and components.

As we have seen before, the two parts of the app depend on several components:

* Redis
* PostgreSQL
* Meilisearch
* Caddy

We have to implement each component as a [Service](https://docs.dagger.io/manuals/developer/services), which then can be used app.
For Redis this could look like this:

```python
    @function
    def redis(self) -> dagger.Service:
        """Returns a redis service from a container built with the given params."""
        return (
            dag.container()
            .from_("redis:alpine")
            .with_exposed_port(6379)
            .as_service()
        )
```

{{% alert title="Note" color="primary" %}}
This Container build does not use a Dockerfile. The Container is defined using the Dagger API.\
The exposing to a Service, that we did with Bash before, is done in the function.
{{% /alert %}}

Add the redis function to your module.


### Task {{% param sectionnumber %}}.1: Implement Services

Add the remaining Services as well. Consult [docker-compse.yaml](https://github.com/mawoka-myblock/ClassQuiz/blob/master/docker-compose.yml)
for the required ports and params.

{{% alert title="Note" color="primary" %}}
A simple check for your function code is to run `dagger functions`. This will build/compile your code including all Dagger dependencies.\
When you see your functions listed, then the syntax is right.
{{% /alert %}}

While the implementations of PostgreSQL and Meilisearch are very similar and quite simple:

{{% details title="show solution" mode-switcher="normalexpertmode" %}}
```python
    @function
    def postgres(self) -> dagger.Service:
        """Returns a postgres database service from a container built with the given params."""
        return (
            dag.container()
            .from_("postgres:14-alpine")
            .with_env_variable("POSTGRES_PASSWORD", "classquiz")
            .with_env_variable("POSTGRES_DB", "classquiz")
            .with_env_variable("POSTGRES_USER", "postgres")
            .with_exposed_port(5432)
            .as_service()
        )

    @function
    def meilisearch(self) -> dagger.Service:
        """Returns a meilisearch service from a container built with the given params."""
        return (
            dag.container()
            .from_("getmeili/meilisearch:v0.28.0")
            .with_exposed_port(7700)
            .as_service()
        )
```
{{% /details %}}

The implementation of Caddy is a bit more sophisticated, as the proxy is our new entry point, which "glues" all the pieces together.

Official documentation about how to [Bind services in functions](https://docs.dagger.io/manuals/developer/services/#bind-services-in-functions)

{{% alert title="Note" color="primary" %}}
Important detail from the docs: The name used for the service binding defines the host name to be used by the function!
{{% /alert %}}

{{% details title="show solution" mode-switcher="normalexpertmode" %}}
```python
    @function
    def proxy(self, context: dagger.Directory, proxy_config: dagger.File) -> dagger.Service:
        """Returns a caddy proxy service encapsulating the front and backend services. This service must be bound to port 8000 in order to match some hard coded configuration: --ports 8000:8080"""
        return (
            dag.container()
            .from_("caddy:alpine")
            .with_service_binding("frontend", self.build(context.directory("frontend")).as_service())
            .with_service_binding("api", self.build(context).as_service())
            .with_file("/etc/caddy/Caddyfile", proxy_config)
            .with_exposed_port(8080)
            .as_service()
        )
```
{{% /details %}}

You could try to run the ClassQuiz app now. But it will not work because of some missing configuration.

{{% details title="show solution" mode-switcher="normalexpertmode" %}}
```
dagger call proxy --context . --proxy-config Caddyfile-docker up --ports 8000:8080
```
{{% /details %}}


### Task {{% param sectionnumber %}}.2: Create separate Front- and Backend functions

Our initial `build` function can be used to create both, front- and backend containers.\
But in fact, the two app parts require different config params and dependencies:\
The frontend only communicates with the api of the backend,
which is encapsulated by the Caddy reverse proxy, while the backend relies on the services we created earlier.

Hints:

* Start with the `backend` and adjust the host part of the urls used in `classquiz/config.py`.
* For `PORT` use the port that you set inside the `Caddyfile-docker` earlier in this lab.

{{% details title="show solution" mode-switcher="normalexpertmode" %}}
```python
    @function
    def backend(self, context: dagger.Directory) -> dagger.Service:
        """Returns a backend service from a container built with the given context, params and service bindings."""
        return (
            dag.container()
            .with_env_variable("MAX_WORKERS", "1")
            .with_env_variable("PORT", "8081")
            .with_service_binding("postgresd", self.postgres())
            .with_service_binding("meilisearchd", self.meilisearch())
            .with_service_binding("redisd", self.redis())
            .build(context)
            .as_service()
        )
```
{{% /details %}}

For convenience, the function returns directly a Service.

And the `frontend`:

{{% details title="show solution" mode-switcher="normalexpertmode" %}}
```python
    @function
    def frontend(self, context: dagger.Directory) -> dagger.Service:
        """Returns a frontend service from a container built with the given context and params."""
        return (
            dag.container()
            .with_env_variable("API_URL", "http://api:8081")
            .with_env_variable("REDIS_URL", "redis://redisd:6379/0?decode_responses=True")
            .build(context)
            .as_service()
        )
```
{{% /details %}}

Now the two service bindings in the `proxy` function can be simplified a bit.

Before:

```python
            .with_service_binding("frontend", self.build(context.directory("frontend")).as_service())
            .with_service_binding("api", self.build(context).as_service())
```

After:

```python
            .with_service_binding("frontend", self.frontend(context.directory("frontend")))
            .with_service_binding("api", self.backend(context))
```

Now we can finally run ClassQuiz locally:

```bash
dagger call proxy --context=. --proxy-config=Caddyfile-docker up --ports=8000:8080
```

And then visit [localhost:8000](http://localhost:8000/) - where, after registering ourselves, we can register, log in and create our survey!


{{% alert title="Note" color="primary" %}}
Sometimes old cookies or session storage corrupts the app.\
To fix this, delete all cookies and session data.
{{% /alert %}}


### Complete Solution

`ci/src/main/__init__.py`:

<!-- markdownlint-capture -->
<!-- markdownlint-disable -->
{{< readfile file="solution/__init__.py" code="true" lang="Python" >}}
<!-- markdownlint-restore -->

`classquiz/config.py`:


{{< readfile file="solution/config.py" code="true" lang="Python" >}}


`Cadyyfile-docker`:

{{< readfile file="solution/Caddyfile-docker" code="true" lang="YAML" >}}
