---
title: "5. Pipeline integration"
weight: 5
sectionnumber: 5
---

## Pipeline integration


### The Challenge

It is time to write pipelines. This is why we are here.

Good, that we already have some functions, that we can reuse - but let's
add another one which executes the Python tests:

```python
    @function
    async def pytest(self, context: dagger.Directory) -> str:
        """Run pytest and return its output."""
        return await (
            self.backend(context)
            .with_exec(["pip", "install", "--upgrade", "pip"])
            .with_exec(["pip", "install", "--upgrade", "pytest"])
            .with_exec(["pytest", "classquiz/tests/", "--ignore=classquiz/tests/test_server.py"])
            .stdout()
        )
```

The tests are run on the backend, therefore we need to build the backend first.


### Task {{% param sectionnumber %}}.1: Create CI function

Add a `ci` function, that first runs the python tests and then returns a directory containing the scan results of the security scan.

{{% details title="show solution" mode-switcher="normalexpertmode" %}}
```python
    @function
    async def ci(self, context: dagger.Directory) -> dagger.Directory:
        """Run all pipeline stages."""
        await self.pytest(context)
        return await self.vulnerability_scan(context)
```
{{% /details %}}


### Task {{% param sectionnumber %}}.2: Run the ci from CLI

Let's see if we can run it from the CLI and have a look at the results:

{{% details title="show solution" mode-switcher="normalexpertmode" %}}
```bash
dagger call ci --context=. export --path=.tmp
```
{{% /details %}}

If everything went well, the scan results should again be found in the directory `.tmp/scans/`.


### Task {{% param sectionnumber %}}.3: Add GitHub action

As final step, we need to call the `ci` function on every push to the repository.

Have a look at [Dagger for GitHub](https://github.com/marketplace/actions/dagger-for-github) first
and then add the action to __your fork__ on GitHub. Keep it simple and trigger the pipeline with every push on every branch.

{{% details title="show solution" mode-switcher="normalexpertmode" %}}
<!-- markdownlint-capture -->
<!-- markdownlint-disable -->
{{< readfile file="solution/dagger.yml" code="true" lang="Yaml" >}}
<!-- markdownlint-restore -->
{{% /details %}}

Add or alter something, push it to the repo and see if the action runs as expected.


### Complete Solution


`ci/src/class_quiz/main.py`:

<!-- markdownlint-capture -->
<!-- markdownlint-disable -->
{{< readfile file="solution/main.py" code="true" lang="Python" >}}
<!-- markdownlint-restore -->

`dagger.yml`:

<!-- markdownlint-capture -->
<!-- markdownlint-disable -->
{{< readfile file="solution/dagger.yml" code="true" lang="Yaml" >}}
<!-- markdownlint-restore -->
