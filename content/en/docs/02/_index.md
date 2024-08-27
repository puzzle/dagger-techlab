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

