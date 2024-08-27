---
title: "2. Daggerize an App"
weight: 2
sectionnumber: 2
---

## 2. Daggerize an App


### The Challenge

After we have learned the basic Dagger Functions, we want to apply our new knowledge to solve a real life problem:
We would like to conduct a survey regarding the popularity of the different Dagger SDKs!


### The Candidate

Fortunately, there is a free open-source quiz app called [ClassQuiz](https://classquiz.de/)
It allows the creation of shareable, fully customizable quizzes and surveys.
The app is split in a frontend and an api part:

* The frontend is written in type script and uses a redis memcache.
* The backend is mostly written in python, uses a postgreSQL database and mailisearch.

Caddy is used as reverse proxy to keep the various parts together.


### The Journey


#### Prerequisites

Check out ClassQuiz:

```bash
git clone https://github.com/mawoka-myblock/ClassQuiz.git
```

Get familiar with the source - take a closer look at the [docker-compse.yaml](https://github.com/mawoka-myblock/ClassQuiz/blob/master/docker-compose.yml),
which is particularly interesting for our purpose.

Since the app has some hard-coded configurations that would interfere with our setup, let's apply the following [patch](config.patch) to the `classquiz/config.py`:

```bash
patch classquiz/config.py < config.patch
```

The app also binds the privileged port `80`, which would be an obstacle as well. So let's replace all occurrences of `:80` in `Caddyfile-docker` with `:8081` or another port of your choice.


#### Initialize a Dagger module

A new Dagger module in Go, Python or TypeScript can be initialized by running `dagger init` inside the app's root directory,
using the `--source` flag to specify a directory for the module's source code.

We will use the **Python SDK** for this example:

```bash
dagger init --sdk=python --source=./ci
```

This leaves us with a generated `dagger.json` module metadata file, an initial `ci/src/main/__init__.py` source code template, `ci/pyproject.toml` and
`ci/requirements.lock` files, as well as a generated `ci/sdk` folder for local development.


## Run the App locally

The generated `ci/src/main/__init__.py` is the starting point, which needs to be extended.

As a first step, we could implement a simple `build` function:

```python
    @function
    def build(self, context: dagger.Directory) -> dagger.Container:
        """Returns a container built with the given context."""
        return (
            dag.container()
            .build(context)
        )
```

This allows us to build the frontend and bind it as a service to the localhost on port 3000:

```bash
dagger -m Classquiz call build --context ./frontend/ with-exposed-port --port 3000 as-service
```

And the backend as well:

```bash
dagger -m Classquiz call build --context . with-exposed-port --port 8000 as-service up
```

But if we have a closer look to the console output, we will discover some error messages due to missing configurations and components.

As we have seen before, the two parts of the app depend on several components:

* Redis
* postgreSQL
* Meilisearch
* Caddy

We have to implement each component as a service, which then can be used app. For Redis this could look like this:

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


### Task 2.1: Implement Services

Add the remaining Services as well. Consult [docker-compse.yaml](https://github.com/mawoka-myblock/ClassQuiz/blob/master/docker-compose.yml)
for the required ports and params.

While the implementations of PostgreSQL and Meilisearch are very similar and quite simple:

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
    def mailisearch(self) -> dagger.Service:
        """Returns a mailisearch service from a container built with the given params."""
        return (
            dag.container()
            .from_("getmeili/meilisearch:v0.28.0")
            .with_exposed_port(7700)
            .as_service()
        )
```

The implementation of Caddy is a bit more sophisticated, as the proxy is our new entry point, which "glues" all the pieces together.

Official documentation about how to [Bind services in functions](https://docs.dagger.io/manuals/developer/services/#bind-services-in-functions)

{{% alert title="Note" color="primary" %}}
Important detail from the docs: The name used for the service binding defines the host name to be used by the function!
{{% /alert %}}

```python
    @function
    def proxy(self, context_backend: dagger.Directory, context_frontend: dagger.Directory, proxy_config: dagger.File) -> dagger.Service:
        """Returns a caddy proxy service encapsulating the front and backend services. This service must be bound to port 8000 in order to match some hard coded configuration: --ports 8000:8080"""
        return (
            dag.container()
            .from_("caddy:alpine")
            .with_service_binding("frontend", self.build(context_frontend).as_service())
            .with_service_binding("api", self.build(context_backend).as_service())
            .with_file("/etc/caddy/Caddyfile", proxy_config)
            .with_exposed_port(8080)
            .as_service()
        )
```


### Task 2.2: Create separate Front- and Backend functions

Our initial `build` function can be used to create both, front- and backend containers.
But in fact, the two app parts require different config params and dependencies: The frontend only communicates with the api of the backend,
which is encapsulated by the Caddy reverse proxy, while the backend relies on the services we created earlier.

Hint: Start with the `backend` and adjust the host part of the urls used in `classquiz/config.py`.

```python
    @function
    def backend(self, context: dagger.Directory) -> dagger.Service:
        """Returns a backend service from a container built with the given context, params and service bindings."""
        return (
            dag.container()
            .with_env_variable("MAX_WORKERS", "1")
            .with_env_variable("PORT", "8081")
            .with_service_binding("postgresd", self.postgres())
            .with_service_binding("mailisearchd", self.mailisearch())
            .with_service_binding("redisd", self.redis())
            .build(context)
            .as_service()
        )
```

For convenience, the function returns directly a Service.

And the `frontend`:

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

Now the two service bindings in the `proxy` function can be simplified a bit. 

Before:

```python
            .with_service_binding("frontend", self.build(context_frontend).as_service())
            .with_service_binding("api", self.build(context_backend).as_service())
```

After:

```python
            .with_service_binding("frontend", self.frontend(context_frontend))
            .with_service_binding("api", self.backend(context_backend))
```

Now we can finally run ClassQuiz locally:

```bash
dagger call proxy --context-frontend=./frontend  --context-backend=.  --proxy-config=Caddyfile-docker up --ports=8000:8080
```

And then visit http://localhost:8000/

Where, after registering ourselves, we can log in and create our survey!


### Complete Solution

`ci/src/main/__init__.py`:

```python
import dagger
from dagger import dag, function, object_type


@object_type
class Classquiz:

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


    @function
    def backend(self, context: dagger.Directory) -> dagger.Service:
        """Returns a backend service from a container built with the given context, params and service bindings."""
        return (
            dag.container()
            .with_env_variable("MAX_WORKERS", "1")
            .with_env_variable("PORT", "8081")
            .with_service_binding("postgresd", self.postgres())
            .with_service_binding("mailisearchd", self.mailisearch())
            .with_service_binding("redisd", self.redis())
            .build(context)
            .as_service()
        )


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
    def mailisearch(self) -> dagger.Service:
        """Returns a mailisearch service from a container built with the given params."""
        return (
            dag.container()
            .from_("getmeili/meilisearch:v0.28.0")
            .with_exposed_port(7700)
            .as_service()
        )


    @function
    def redis(self) -> dagger.Service:
        """Returns a redis service from a container built with the given params."""
        return (
            dag.container()
            .from_("redis:alpine")
            .with_exposed_port(6379)
            .as_service()
        )


    @function
    def proxy(self, context_backend: dagger.Directory, context_frontend: dagger.Directory, proxy_config: dagger.File) -> dagger.Service:
        """Returns a caddy proxy service encapsulating the front and backend services. This service must be bound to port 8000 in order to match some hard coded configuration: --ports 8000:8080"""
        return (
            dag.container()
            .from_("caddy:alpine")
            .with_service_binding("frontend", self.frontend(context_frontend))
            .with_service_binding("api", self.backend(context_backend))
            .with_file("/etc/caddy/Caddyfile", proxy_config)
            .with_exposed_port(8080)
            .as_service()
        )

```

`classquiz/config.py`:

```python
# SPDX-FileCopyrightText: 2023 Marlon W (Mawoka)
#
# SPDX-License-Identifier: MPL-2.0

import re
from functools import lru_cache

from redis import asyncio as redis_lib
import redis as redis_base_lib
from pydantic import BaseSettings, RedisDsn, PostgresDsn, BaseModel
import meilisearch as MeiliSearch
from typing import Optional
from arq import create_pool
from arq.connections import RedisSettings, ArqRedis

from classquiz.storage import Storage


class CustomOpenIDProvider(BaseModel):
    scopes: str = "openid email profile"
    server_metadata_url: str
    client_id: str
    client_secret: str


class Settings(BaseSettings):
    """
    Settings class for the shop app.
    """

    root_address: str = "http://127.0.0.1:8000"
    redis: RedisDsn = "redis://redisd:6379/0?decode_responses=True"
    skip_email_verification: bool = True
    db_url: str | PostgresDsn = "postgresql://postgres:classquiz@postgresd:5432/classquiz"
    hcaptcha_key: str | None = None
    recaptcha_key: str | None = None
    mail_address: str = "some@example.org"
    mail_password: str = "some@example.org"
    mail_username: str = "some@example.org"
    mail_server: str = "some.example.org"
    mail_port: int = "525"
    secret_key: str = "secret"
    access_token_expire_minutes: int = 30
    cache_expiry: int = 86400
    sentry_dsn: str | None
    meilisearch_url: str = "http://mailisearchd:7700"
    meilisearch_index: str = "classquiz"
    google_client_id: Optional[str]
    google_client_secret: Optional[str]
    github_client_id: Optional[str]
    github_client_secret: Optional[str]
    custom_openid_provider: CustomOpenIDProvider | None = None
    telemetry_enabled: bool = True
    free_storage_limit: int = 1074000000
    pixabay_api_key: str | None = None
    mods: list[str] = []
    registration_disabled: bool = False

    # storage_backend
    storage_backend: str | None = "local"

    # if storage_backend == "local":
    storage_path: str | None = "/app/data"

    # if storage_backend == "s3":
    s3_access_key: str | None
    s3_secret_key: str | None
    s3_bucket_name: str = "classquiz"
    s3_base_url: str | None

    class Config:
        env_file = ".env"
        env_file_encoding = "utf-8"
        env_nested_delimiter = "__"


async def initialize_arq():
    # skipcq: PYL-W0603
    global arq
    arq = await create_pool(RedisSettings.from_dsn(settings.redis))


@lru_cache()
def settings() -> Settings:
    return Settings()


pool = redis_lib.ConnectionPool().from_url(settings().redis)

redis: redis_base_lib.client.Redis = redis_lib.Redis(connection_pool=pool)
arq: ArqRedis = ArqRedis(pool_or_conn=pool)
storage: Storage = Storage(
    backend=settings().storage_backend,
    storage_path=settings().storage_path,
    access_key=settings().s3_access_key,
    secret_key=settings().s3_secret_key,
    bucket_name=settings().s3_bucket_name,
    base_url=settings().s3_base_url,
)

meilisearch = MeiliSearch.Client(settings().meilisearch_url)

ALLOWED_TAGS_FOR_QUIZ = ["b", "strong", "i", "em", "small", "mark", "del", "sub", "sup"]

ALLOWED_MIME_TYPES = ["image/png", "video/mp4", "image/jpeg", "image/gif", "image/webp"]

server_regex = rf"^{re.escape(settings().root_address)}/api/v1/storage/download/.{{36}}--.{{36}}$"

```

`Cadyyfile-docker`:

```yaml
# SPDX-FileCopyrightText: 2023 Marlon W (Mawoka)
#
# SPDX-License-Identifier: MPL-2.0

:8080 {
	reverse_proxy * http://frontend:3000
	reverse_proxy /api/* http://api:8081
	reverse_proxy /openapi.json http://api:8081 # Only use if you need to serve the OpenAPI spec
	reverse_proxy /socket.io/* http://api:8081

}
```