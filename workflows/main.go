// A helper module for workflow functions
//
// This module helps to develop the hugo application.

package main

import (
	"context"
	"dagger/workflows/internal/dagger"
	"fmt"
	"strings"
)

type Workflows struct{}

// Get the Hugo Image tag
func (m *Workflows) HugoTag(
	ctx context.Context,
	src *dagger.Directory,
) (string, error) {
	return dag.Container().
		From("alpine:latest").
		WithMountedDirectory("/mnt", src).
		WithWorkdir("/mnt").
		WithExec([]string{"sh", "-c", "grep \"FROM docker.io/floryn90/hugo\" Dockerfile | sed 's/FROM docker.io\\/floryn90\\/hugo://g' | sed 's/ AS builder//g'"}).
		Stdout(ctx)
}

// Prepares the local dev container with the given source.
func (m *Workflows) LocalDev(
	ctx context.Context,
	src *dagger.Directory,
) (*dagger.Container, error) {
	tag,err := m.HugoTag(ctx, src)
	if err != nil {
		return nil, err
	}
	from := strings.TrimSpace(fmt.Sprintf("docker.io/floryn90/hugo:%s", tag))
	container := dag.Container().
		From(from).
		WithMountedDirectory("/src", src, dagger.ContainerWithMountedDirectoryOpts{Owner: "hugo"}).
		WithExposedPort(8080)
	return container, nil
}

// Builds and runs Hugo from the given source.
func (m *Workflows) LocalStart(
	ctx context.Context,
	src *dagger.Directory,
) (*dagger.Service, error) {
	container, err := m.LocalDev(ctx, src)
	if err != nil {
		return nil, err
	}
	service := container.
		AsService(dagger.ContainerAsServiceOpts{Args: []string{"hugo", "server", "-p", "8080"}})
	return service, nil
}

// Runs lint on the given source.
func (m *Workflows) Lint(
	ctx context.Context,
	src *dagger.Directory,
) (string, error) {
	tag,err := m.HugoTag(ctx, src)
	if err != nil {
		return "", err
	}
	from := strings.TrimSpace(fmt.Sprintf("docker.io/floryn90/hugo:%s", tag))
	return dag.Container().
		From(fmt.Sprintf("%s-ci", from)).
		WithMountedDirectory("/mnt", src, dagger.ContainerWithMountedDirectoryOpts{Owner: "hugo"}).
		WithWorkdir("/mnt").
		WithExec([]string{"/bin/bash", "-c", "npm install && npm run mdlint"}).
		Stdout(ctx)
}
