# Dagger Techlab

This guided hands-on training helps dagger.io newbies to get started. It also provides information and labs for more advanced topics.


## Content Sections

The training content resides within the [content](content) directory.

The main part are the labs, which can be found at [content/en/docs](content/en/docs).


## Hugo

This site is built using the static page generator [Hugo](https://gohugo.io/).

The page uses the [docsy theme](https://github.com/google/docsy) which is included as a Hugo Module.
Docsy is being enhanced using [docsy-plus](https://github.com/acend/docsy-plus/) as well as
[docsy-acend](https://github.com/acend/docsy-acend/) and [docsy-puzzle](https://github.com/puzzle/docsy-puzzle/)
for brand specific settings.

After cloning the main repo, you need to initialize the Hugo Module like this:

```bash
hugo mod get
```

The default configuration uses the acend setup from [config/_default](config/_default/config.toml).
Alternatively you can use the Puzzle setup from [config/puzzle](config/puzzle/config.toml), which is enabled with
`--environment puzzle`.


### Docsy theme usage

* [Official docsy documentation](https://www.docsy.dev/docs/)
* [Docsy Plus](https://github.com/acend/docsy-plus/)


### Update submodules for theme updates

Run the following command to update all submodules with their newest upstream version:

```bash
hugo mod get -u
```


### Shortcode usage


#### `onlyWhen` and `onlyWhenNot`

The `onlyWhen` and `onlyWhenNot` shortcodes allow text to be rendered if certain conditions apply.

* `{{% onlyWhen variant1 %}}`: This is only rendered when `enabledModule` in `config.toml` contains `variant1`
* `{{% onlyWhen variant1 variant2 %}}`: This is only rendered when `enabledModule` in `config.toml` contains `variant1` **or** `variant2`
* `{{% onlyWhenNot variant1 %}}`: This is only rendered when `enabledModule` in `config.toml` **does not** contain `variant1`
* `{{% onlyWhenNot variant1 variant2 %}}`: This is only rendered when `enabledModule` in `config.toml` **does not** contain `variant1` **or** `variant2`

In order to only render text if **all** of multiple conditions do not apply simply chain several `onlyWhenNot` shortcodes:

```
{{% onlyWhenNot variant1 %}}
{{% onlyWhenNot variant2 %}}
This is only rendered when `enabledModule` in `config.toml` **does not** contain `variant1` **nor** `variant2`.
{{% /onlyWhen %}}
{{% /onlyWhen %}}
```


## Build production image locally


### Docker

Build the image:

```bash
docker build -t puzzle/dagger-techlab .
```

Run it locally:

> [!WARNING]  
> When building locally and applying changes

```bash
docker run --rm -p 8080:8080 puzzle/dagger-techlab
```


### Buildah and Podman

Build the image:

```bash
buildah build-using-dockerfile [--build-arg TRAINING_HUGO_ENV=...] -t puzzle/dagger-techlab .
```

Run it locally:

```bash
podman run --rm --rmi --publish 8080:8080 localhost/puzzle/dagger-techlab
```

**Note:** Beware that `--rmi` automatically removes the built image when the container stops, so you either have to rebuild it or remove the parameter from the command.


## How to develop locally


### Dagger

Start Hugo with the local content:

```bash
dagger -m workflows/ call local-start --src . up --ports 8080:8080
```


### Docker Compose

You can use `docker-compose`. If you prefer Podman check out [podman-compose](https://github.com/containers/podman-compose).

```bash
docker-compose up
```

Use the following command to set the hugo environment

```bash
HUGO_ENVIRONMENT="something" docker-compose up
```

To rebuild the image if something changed in the `Dockerfile`:

```bash
docker-compose up --build
```


**Tip:** Set the following environment variables for faster builds: `DOCKER_BUILDKIT=1 COMPOSE_DOCKER_CLI_BUILD=1`


### Plain Docker

To develop locally we don't want to rebuild the entire container image every time something changed, and it is also important to use the same hugo versions like in production.
We simply mount the working directory into a running container, where hugo is started in the server mode.

```bash
export HUGO_VERSION=$(grep "FROM docker.io/floryn90/hugo" Dockerfile | sed 's/FROM docker.io\/floryn90\/hugo://g' | sed 's/ AS builder//g')
docker run --rm --publish 8080:8080 -v $(pwd):/src docker.io/floryn90/hugo:${HUGO_VERSION} server -p 8080
```

Use the following command to set the hugo environment

```bash
export HUGO_VERSION=$(grep "FROM docker.io/floryn90/hugo" Dockerfile | sed 's/FROM docker.io\/floryn90\/hugo://g' | sed 's/ AS builder//g')
docker run --rm --publish 8080:8080 -v $(pwd):/src docker.io/floryn90/hugo:${HUGO_VERSION} server --environment=<environment> -p 8080
```


## Linting of Markdown content

Markdown files are linted with <https://github.com/DavidAnson/markdownlint>.
Custom rules are in `.markdownlint.json`.
There's a GitHub Action `.github/workflows/markdownlint.yaml` for CI.
For local checks, you can either use Visual Studio Code with the corresponding extension, or the command line like this:

```shell script
npm install
npm run mdlint
```

Npm not installed? no problem

```bash
export HUGO_VERSION=$(grep "FROM docker.io/floryn90/hugo" Dockerfile | sed 's/FROM docker.io\/floryn90\/hugo://g' | sed 's/ AS builder//g')
docker run --rm -v $(pwd):/src docker.io/floryn90/hugo:${HUGO_VERSION}-ci /bin/bash -c "npm install && npm run mdlint"
```

Automatically fix errors if possible:

```bash
npm run mdlint-fix
```


## Github Actions

Needed GitHub Repository Actions secrets:

* `KUBECONFIG_TEST_AV2`: K8S cluster configuration incl. auth. for the test namespace
* `KUBECONFIG_PROD_AV2`: K8S cluster configuration incl. auth. for the prod namespace

Create the cluster configuration:

Copy the login token of the ServiceAccount in the namespace where the GitHub Action should have access to and export it as `SA_TOKEN`.
Use then this token to create the kubeconfig. Follow these commands.

```bash
export SA_TOKEN=<ServiceAccount Token>
export KUBECONFIG=tmp/.kube/dev_cluster_config
oc login --token=${SA_TOKEN} --server=<K8S-cluster-api-url>
cat tmp/.kube/dev_cluster_config
```

The output of the `cat` command has to be used for the GitHub Action secret.


### Build

The [build action](.github/workflows/build.yaml) is fired on Pull Requests does the following

* builds all PR Versions (Linting and Docker build)
* deploys the built container images to the container registry
* Deploys a PR environment in a k8s test namespace with helm
* Triggers a redeployment
* Comments in the PR where the PR Environments can be found


### PR Cleanup

The [pr-cleanup action](.github/workflows/pr-cleanup.yaml) is fired when Pull Requests are closed and does the following

* Uninstalls PR Helm Release


### Push Main

The [push main action](.github/workflows/push-main.yaml) is fired when a commit is pushed to the main branch (eg. a PR is merged) and does the following, it's very similar to the Build Action

* builds main Versions (Linting and Docker build)
* deploys the built container images to the container registry
* Deploys the main Version on k8s using helm
* Triggers a redeployment


## Helm

Manually deploy the training Release using the following command:

```bash
helm install --repo https://acend.github.io/helm-charts/  <release> acend-training-chart --values helm-chart/values.yaml -n <namespace>
```

For debugging purposes use the `--dry-run` parameter

```bash
helm install --dry-run --repo https://acend.github.io/helm-charts/  <release> acend-training-chart --values helm-chart/values.yaml -n <namespace>
```


## Contributions

If you find errors, bugs or missing information please help us improve and have a look at the [Contribution Guide](CONTRIBUTING.md).
