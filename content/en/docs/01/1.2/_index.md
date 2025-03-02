---
title: "1.2 What is Dagger?"
weight: 12
sectionnumber: 1.2
description: >
  Introduction to Dagger.
---


## How does dagger work under the hood?


"Dagger is Docker for your CI Pipeline" - Jeremy, propably

Let's start at the beginning. Contaners have been around since the early 2000's, but only experienced a true upswing once Docker was introduced in 2013, as it mad working with containers easier and more intuitiv.
Dagger works with similar principals as Docker, this is not surprising as the founder of Docker is also the founder of Dagger.

TDagger is based on an engine, at its heart the workflow is depicted as a Directed Acyclic Graph (DAG). The dagger engine receives outside communication via an GraphQL-API, which allows for an declarative definition of the separate pipeline steps. The output of the DAG is passed on to BildKit. The whole process is executed in a container, this ensures consistency and portability.

Picture 4


### Host components


* **Input/Output**
* **SDK**
* **Dagger CLI**


### GraphQL API


### Dagger Engine


### BuildKit

