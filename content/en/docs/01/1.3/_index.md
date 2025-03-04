---
title: "1.3 Dagger is more..."
weight: 13
sectionnumber: 1.3
description: >
  Dagger is more than push and flow.
---


## Dagger is more than push and flow


### More solutions for CI

Dagger solves more than the "push-and-pray" problem in CI. 
It also introduces a way for developers to directly write tests in the coding language they are comfortable in by utilizing software developement kits (SDKs).
At the moment Dagger offers SDKs for the following coding languages: Go, Python, Typescript, PHP and the newest addition (in beta) [Java](https://dagger.io/blog/java-sdk).

The Dagger SDKs translate your code/functions internally into GraphQL API calls to communicate with the DAG in the Dagger RunTime.

To increase the quality and efficiency, tandardized Dagger functions are shared in the [Daggerverse](https://daggerverse.dev/), so you don't have to write your own code to for example, run a k3s server that can be accessed both locally and in your pipelines or to utilize helm and many more.

Running your CI pipeline locally with Dagger also helps you be more efficient by cutting down on the pipeline wait times by utilizing caching. Like this you can avoid unnecessary rebuilds and test reruns when nothing has changed. Only rerunning the parts of your pipeline that you changed.

To increase colaboration and visualize your pipelines you can use Dagger Cloud.


### Dagger for AI agents

Picture Robots

connect to LLM endpoint of your choice to access your Dagger objects

