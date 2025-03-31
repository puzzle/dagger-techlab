---
title: "6. Container interaction and debugging"
weight: 6
sectionnumber: 6
---

## Container interaction

Sometimes things don't go as smoothly as expected.

For such rare cases, Dagger gives us the possibility to start an interactive session to inspect a specific container:

```bash
dagger call backend --context=. terminal
```

This attaches us to `sh` shell inside the container returned by the `bakend` function.
To exit press `ctrl-D` repeatedly.

If you prefer another command, you can pass the optional `cmd` argument:

```bash
dagger call backend --context=. terminal --cmd=printenv
```

In this case it executes the command and terminates afterward.


### Task {{% param sectionnumber %}}.1: Start a live terminal session

Connect to the `frontend` container with a live terminal session using `bash` and list the installed debian packages.

{{% details title="show solution" mode-switcher="normalexpertmode" %}}
```bash
dagger call frontend --context=frontend terminal --cmd=bash
# dpkg -l
```
{{% /details %}}


## Debugging


### Debug flag

In other cases, it would be helpful if we just could get a bit more detailed output.

For this purpose, every Dagger command can be invoked with the `debug` flag:

```bash
dagger call --debug backend --context=.
```


### Interactive flag

Dagger has an interactive debugging feature, which allows users to drop in to an interactive shell when a pipeline run fails, with all the context at the point of failure.

No need to set breakpoints or change your code.

Just do your dagger call and include the `--interactive` flag. E.g:

```bash
dagger call --interactive foo
```

Use the interactive (`-i`) shorthand to enter the container when a problem occurs:

```bash
dagger -i -c 'container | from alpine | with-exec echooo "Daggernaut" | stdout'
```
{{% alert title="Note" color="primary" %}}
There is no `echooo` command available inside the container.\
The call will stop in the alpine container. There you can check the environment and find the right command.\
Press `Ctrl+D` to exit.
{{% /alert %}}


Utilizing the interactive flag, you can also set breakpoints in your pipeline `.terminal()`

```python
import dagger
from dagger import dag, function, object_type


@object_type
class MyModule:
    @function
    async def foo(self) -> dagger.Container:
        return await (
            dag.container()
            .from_("alpine:latest")
            .terminal()
            .with_exec(["sh", "-c", "echo hello world > /foo"])
            .terminal()
        )
```


### Debugging with Dagger Cloud

link:https://dagger.io/cloud[Dagger Cloud] allows you to run your dagger calls online, giving you an overview across all pipelines, both pre-push and post-push.

As a single user you can register for free link:https://dagger.cloud/signup[here].

Give it a try :)

