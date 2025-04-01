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

Fortunately, there is a free open-source quiz app called [ClassQuiz](https://classquiz.de/).\
It allows the creation of shareable, fully customizable quizzes and surveys.\
The app is split in a frontend and an api part:

* The frontend is written in TypeScript and uses a Redis memcache.
* The backend is mostly written in python, uses a PostgreSQL database and Meilisearch.

Caddy is used as a reverse proxy to keep the parts together.


### The Journey


#### Prerequisites

Create a fork of the [ClassQuiz Repo on Github](https://github.com/mawoka-myblock/ClassQuiz), this will come in handy later.

Then, check out your fork of ClassQuiz:

```bash
git clone https://github.com/_your-Github-user_/ClassQuiz.git
```

Get familiar with the source - take a closer look at the [docker-compse.yaml](https://github.com/mawoka-myblock/ClassQuiz/blob/master/docker-compose.yml),
which is particularly interesting for our purpose.

The app binds the privileged port `80`, which would be an obstacle.\
So let's replace all occurrences of `:80` in `Caddyfile-docker` with `:8081`.\
Additionally the missing protocol has to be added to the last `reverse_proxy` line. Add `http://` in front of `api:80`.

Do it by hand or use the following `sed` commands:

```bash
sed -i 's# api:80# http://api:80#g' Caddyfile-docker
sed -i 's#api:80#api:8081#g' Caddyfile-docker
```

If patching does not work, overwrite the file `Caddyfile-docker` with the content from the following `Caddyfile-docker` file.

{{< details title="show final Caddyfile-docker file" >}}

{{< readfile file="solution/Caddyfile-docker" code="true" >}}

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

This leaves us with a generated `dagger.json` module metadata file, an initial `ci/src/class_quiz/main.py` source code template, `ci/pyproject.toml` and
other needed files, as well as a generated `ci/sdk` folder for local development.\
The configuration file sets the name of the module to the name of the current directory, unless an alternative is specified with the `--name` argument.

To check if the module works and what example functions are created, run the `functions` command.

```bash
$ dagger functions
✔ connect 0.3s
✔ load module 1.4s

Name             Description
container-echo   Returns a container that echoes whatever string argument is provided
grep-dir         Returns lines that match a pattern in the files of the provided Directory
```


## Run the App locally

The generated `ci/src/class_quiz/main.py` is the starting point, which needs to be extended.\
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
[Here](https://dagger-io.readthedocs.io/en/latest/client.html#dagger.Container) is the reference to the Python documentation.

The [build](https://dagger-io.readthedocs.io/en/latest/client.html#dagger.Container.build) executes the Docker build with the given files.

This function allows us to build the frontend as Container.\
With function chaining we expose the container as a [Service](https://docs.dagger.io/manuals/developer/services) to the [localhost](https://docs.dagger.io/manuals/developer/services#expose-services-returned-by-functions-to-the-host) on port 3000:

```bash
dagger call --mod ./ci/ build --context=./frontend/ with-exposed-port --port=3000 as-service up
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

{{% alert title="Warning" color="primary" %}}
Unfortunately the Dagger call stops after a while. We have to analyze this!
{{% /alert %}}

If we do not see the relevant logs of the app in the output of the Dagger call, we should change the verbosity.\
Try to make the output more verbose. This is implemented with the `-v` option.

Run the call again with the verbosity option:

```bash
dagger call -v build --context=. with-exposed-port --port=8000 as-service up
```

{{% alert title="Note" color="primary" %}}
If the output does still not contain the needed information of the problem,
increase the verbosity of the Dagger call even more to get to the goal.\
E.g. by adding two more levels at once (`-v` -> `-vvv`)
{{% /alert %}}

If we have a closer look to the console output, we will discover some error messages due to missing configurations.

As we have seen before, the two parts of the app depend on several components:

* Redis
* PostgreSQL
* Meilisearch
* Caddy

We have to implement each component as a [Service](https://docs.dagger.io/manuals/developer/services), which then can be used.
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

Official documentation about how to [Bind services in functions](https://docs.dagger.io/manuals/developer/services/#bind-services-in-functions).

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
```bash
dagger call proxy --context=. --proxy-config=Caddyfile-docker up --ports=8000:8080
```
{{% /details %}}


### Task {{% param sectionnumber %}}.2: Create separate Front- and Backend functions

Our initial `build` function can be used to create both, front- and backend containers.\
But in fact, the two app parts require different config params and dependencies:\
The frontend only communicates with the api of the backend,
which is encapsulated by the Caddy reverse proxy, while the backend relies on the services we created earlier.

Hints:

* Start with the `backend` and pass the required environment variables found in the `docker-compose.yml`.
* For `PORT` use the port that you set inside the `Caddyfile-docker` earlier in this lab.

{{% details title="show required environment variables" mode-switcher="normalexpertmode" %}}
```
MAX_WORKERS
PORT
REDIS
SKIP_EMAIL_VERIFICATION
DB_URL
MAIL_ADDRESS
MAIL_PASSWORD
MAIL_USERNAME
MAIL_SERVER
MAIL_PORT
SECRET_KEY
MEILISEARCH_URL
STORAGE_BACKEND
STORAGE_PATH
```
{{% /details %}}

{{% details title="show solution" mode-switcher="normalexpertmode" %}}
```python
    @function
    def backend(self, context: dagger.Directory) -> dagger.Container:
        """Returns a backend container built with the given context, params and service bindings."""
        return (
            dag.container()
            .with_env_variable("MAX_WORKERS", "1")
            .with_env_variable("PORT", "8081")
            .with_env_variable("REDIS", "redis://redisd:6379/0?decode_responses=True")
            .with_env_variable("SKIP_EMAIL_VERIFICATION", "True")
            .with_env_variable("DB_URL", "postgresql://postgres:classquiz@postgresd:5432/classquiz")
            .with_env_variable("MAIL_ADDRESS", "some@example.org")
            .with_env_variable("MAIL_PASSWORD", "some@example.org")
            .with_env_variable("MAIL_USERNAME", "some@example.org")
            .with_env_variable("MAIL_SERVER", "some.example.org")
            .with_env_variable("MAIL_PORT", "525")
            .with_env_variable("SECRET_KEY", "secret")
            .with_env_variable("MEILISEARCH_URL", "http://meilisearchd:7700")
            .with_env_variable("STORAGE_BACKEND", "local")
            .with_env_variable("STORAGE_PATH", "/app/data")
            .with_service_binding("postgresd", self.postgres())
            .with_service_binding("meilisearchd", self.meilisearch())
            .with_service_binding("redisd", self.redis())
            .build(context)
        )
```
{{% /details %}}

And the `frontend`:

{{% details title="show solution" mode-switcher="normalexpertmode" %}}
```python
    @function
    def frontend(self, context: dagger.Directory) -> dagger.Container:
        """Returns a frontend container built with the given context and params."""
        return (
            dag.container()
            .with_env_variable("API_URL", "http://api:8081")
            .with_env_variable("REDIS_URL", "redis://redisd:6379/0?decode_responses=True")
            .build(context)
        )
```
{{% /details %}}

Now the two service bindings in the `proxy` function needs to be changed accordingly.

Before:

```python
            .with_service_binding("frontend", self.build(context.directory("frontend")).as_service())
            .with_service_binding("api", self.build(context).as_service())
```

After:

```python
            .with_service_binding("frontend", self.frontend(context.directory("frontend")).as_service())
            .with_service_binding("api", self.backend(context).as_service())
```

Now we can run ClassQuiz locally:

```bash
dagger call proxy --context=. --proxy-config=Caddyfile-docker up --ports=8000:8080
```

And then visit [localhost:8000](http://localhost:8000/) - where, after registering ourselves, we can log in and create our survey!


{{% alert title="Note" color="primary" %}}
Sometimes old cookies or session storage corrupts the app, especially when applying changes.\
To fix this, delete all cookies and session data or open it in an incognito tab.
{{% /alert %}}


### Complete Solution

`ci/src/class_quiz/main.py`:

<!-- markdownlint-capture -->
<!-- markdownlint-disable -->
{{< readfile file="solution/main.py" code="true" lang="Python" >}}
<!-- markdownlint-restore -->

`Cadyyfile-docker`:

{{< readfile file="solution/Caddyfile-docker" code="true" lang="YAML" >}}

<!--
TODO: add frontend build by Dagger API 

### Additional / Advanced


    @function
    def build_frontend(self, src: dagger.Directory) -> dagger.Container:
        """Build and publish Docker container"""
        # build container
        builder = (
            dag.container()
            .from_("node:19-bullseye")
            .with_env_variable("API_URL", "http://api:8081")
            .with_env_variable("REDIS_URL", "redis://redisd:6379/0?decode_responses=True")
            # .with_env_variable("ENV API_URL", "https://mawoka.eu")
            # .with_env_variable("ENV REDIS_URL", "redis://localhost:6379")
            # .with_env_variable("ENV VITE_MAPBOX_ACCESS_TOKEN", "pk.eyJ1IjoibWF3b2thIiwiYSI6ImNsMjBob3d4ZjBhcGszYnE0bWp4aXB1ZW4ifQ.IByxV1qeIuEWpHCWsuB88A")
            # .with_env_variable("ENV VITE_HCAPTCHA", "ee81b2a1-acf3-4d20-b2a4-a7ea94c7eba5") 
            .with_directory("/usr/src/app", src)
            .with_workdir("/usr/src/app")
            .with_exec(["corepack", "enable"])
            .with_exec(["corepack", "prepare", "pnpm@8.14.0", "--activate"])
            .with_exec(["pnpm", "i"])
            .with_exec(["pnpm", "run", "build"])
        )

        # runtime image
        runtime_image = (
            dag.container()
            .from_("node:19-bullseye-slim")
            .with_workdir("/app")
            .with_file("/app/package.json", builder.file("/usr/src/app/package.json"))
            .with_file("/app/pnpm-lock.yaml", builder.file("/usr/src/app/pnpm-lock.yaml"))
            .with_exec(["corepack", "enable"])
            .with_exec(["corepack", "prepare", "pnpm@8.14.0", "--activate"])
            .with_exec(["pnpm", "i"])
            .with_directory("/app/", builder.directory("/usr/src/app/build/"))
            .with_directory("/app/node_modules/", builder.directory("/usr/src/app/node_modules/"))
            .with_exec(["pnpm", "run", "run:prod"])
            .with_exposed_port(3000)
        )
        return runtime_image


{{% alert title="Note" color="primary" %}}
When using dagger call, all names (functions, arguments, struct fields, etc) are converted into a shell-friendly "kebab-case" style.
{{% /alert %}}

-->
