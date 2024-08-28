---
title: "3. Pipeline integration"
weight: 3
sectionnumber: 3
---

## {{% param sectionnumber %}}. Pipeline integration


### The Challenge

After we have learned how to spin up an app locally, we want to integrate it into CI-Pipeline.


### Vulnerability scan

So instead of starting the app, we want to build the container images and scan them for vulnerabilities.
A very popular open-source tool for this task is [Trivy](https://trivy.dev/).

We could simply create and start a Trivy Dagger Container like we did for Redis in the previous Lab.

But wait - as it is such a popular tool, maybe someone already did this before and shared its solution?

So let's visit [daggerverse.dev](https://daggerverse.dev). Here we can find hundreds of ready to use Dagger Modules.

At the time of writing, a search for `trivy` reveals six different Modules (+1 just containing examples) -
which leaves us with six solutions to one problem :)


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
