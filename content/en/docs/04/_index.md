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
