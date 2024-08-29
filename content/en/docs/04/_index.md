---
title: "4. Daggerverse and Modules"
weight: 4
sectionnumber: 4
---

## {{% param sectionnumber %}}. Daggerverse and Modules

Until here, we have successfully daggerized an application using the Dagger API.
But there is more: Dagger allows you to reuse Dagger Functions developed by others, which were published to the [Daggerverse](https://daggerverse.dev)!

So let's visit the [Daggerverse](https://daggerverse.dev) and explore it a bit.
Here we find hundreds of ready to use Dagger Modules - and each one of them extends Dagger with one or more additional Functions!
We can also navigate to the (Git) repositories and inspect the source code of each published Module.

{{% alert title="Note" color="primary" %}}
As Dagger Functions can call other functions across languages, the language a Module is written in doesn't matter!
{{% /alert %}}

Let's search for [Trivy](https://trivy.dev/), a very popular open-source vulnerability scann tool
At the time of writing, a search for `trivy` reveals six Modules (+1 just containing examples) -
which leaves us with six different solutions to one problem :)


### Task {{% param sectionnumber %}}.1: Install a Module from Daggerverse

For our task, we chose the [github.com/sagikazarmark/daggerverse/trivy](https://daggerverse.dev/mod/github.com/sagikazarmark/daggerverse/trivy@5b826062b6bc1bfbd619aa5d0fba117190c85aba) Module.

Explore its page, have a look at the available functions which are documented on the left side.

After that, add it to our project by installing it

{{% details title="show solution" mode-switcher="normalexpertmode" %}}
```bash
dagger install github.com/sagikazarmark/daggerverse/trivy@v0.5.0
```
{{% /details %}}

Dagger downloaded the Module and added it as dependency to our `dagger.json`:

```json
{
  "name": "classquiz",
  "sdk": "python",
  "dependencies": [
    {
      "name": "trivy",
      "source": "github.com/sagikazarmark/daggerverse/trivy@5b826062b6bc1bfbd619aa5d0fba117190c85aba"
    }
  ],
  "source": "ci",
  "engineVersion": "v0.12.5"
}
```

This way, all the functions provided by the module are available directly in our code - no need to add further imports or anything like that!

{{% alert title="Note" color="primary" %}}
You may wonder why the dependency contains `trivy@5b826062b6bc1bfbd619aa5d0fba117190c85aba` while we wanted to install `trivy@v0.5.0`?
This is not a mistake: Dagger enforces version pinning, which guarantees that the module version to be installed always remains exactly the same!
{{% /details %}}


## The Challenge

We already know how to spin up our app locally, but now it's time to do some security tests.\
So instead of starting the app, we want to build the container images and scan them for vulnerabilities.

We could simply create and start a [Trivy](https://trivy.dev/) Dagger Container using the Dagger API like we did for Redis & Co.\
But after what we learned previously, we will, of course, use the functions of the [Trivy Module](https://daggerverse.dev/mod/github.com/sagikazarmark/daggerverse/trivy@5b826062b6bc1bfbd619aa5d0fba117190c85aba)!


### Extending our codebase

The Trivy Module has a `container()` function, which expects [Container](https://docs.dagger.io/api/reference/#definition-Container) as argument.
As our existing `frontend()` and `backend()` return [Service](https://docs.dagger.io/api/reference/#definition-Service)s, we need an additional functions.

Since the only difference in creating these containers is the path in which they are built, we will combine them into a single function:

```python
    @function
    async def build(self, context: dagger.Directory) -> dagger.Container:
        """Returns a container built with the given context."""
        return await dag.container().build(context)
```

For better scalability, we have defined the function as asynchronous.

Add this `build` method to your module.


### Task {{% param sectionnumber %}}.2: Add Trivy scan

Now that everything is prepared, it's time to add the actual Trivy scan:

Add a `vulnerability_scan` function, which returns a directory containing the scan results.

{{% details title="show solution" mode-switcher="normalexpertmode" %}}
```python
    @function
    async def vulnerability_scan(self, context: dagger.Directory) -> dagger.Directory:
        """Builds the front- and backend, performs a Trivy scan and returns the directory containing the reports."""
        trivy = dag.trivy()

        directory = (
            dag.directory()
            .with_file("scans/backend.sarif", trivy.container(await self.build(context)).report("sarif"))
            .with_file("scans/frontend.sarif", trivy.container(await self.build(context.directory("frontend"))).report("sarif"))
        )
        return directory
```
{{% /details %}}


### Task {{% param sectionnumber %}}.3: Run the scan from CLI

Let's see if we can run it from the CLI and have a look at the results:

{{% details title="show solution" mode-switcher="normalexpertmode" %}}
```bash
dagger -m Classquiz call vulnerability-scan --context=. export --path=.tmp
```
{{% /details %}}

{{% alert title="Note" color="primary" %}}
When using dagger call, all names (functions, arguments, struct fields, etc) are converted into a shell-friendly "kebab-case" style.
{{% /alert %}}

If everything went well, the scan results should be found in the directory `.tmp/scans/`.


### Complete Solution

`ci/src/main/__init__.py`:

<!-- markdownlint-capture -->
<!-- markdownlint-disable -->
{{< readfile file="solution/__init__.py" code="true" lang="Python" >}}
<!-- markdownlint-restore -->
