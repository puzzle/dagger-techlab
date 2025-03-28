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

Just do your dagger call and include the `--interactive` flag.

```bash
dagger call --interactive foo
```

Utilizing the interactive flag, you can also set breakpoints in your pipeline

