// A module to support the Puzzle dagger techlab.
//
// This module has been generated via dagger init --sdk go.
//
// The functions are used inside the hands-on lab: https://dagger-techlab.puzzle.ch/

package main

import (
	"context"
	"dagger/mod/internal/dagger"
	"fmt"
	"strings"
)

var defaultFigletContainer = dag.
	Container().
	From("alpine:latest").
	WithExec([]string{
		"apk", "add", "figlet",
	})

type Mod struct{}

// Say hello to the world!
// Code taken from https://github.com/shykes/hello. Thanks!
func (hello *Mod) Hello(
	ctx context.Context,
	// Change the greeting
	// +optional
	// +default="hello"
	greeting string,
	// Change the name
	// +optional
	// +default="world"
	name string,
	// Encode the message in giant multi-character letters
	// +optional
	giant bool,
	// Make the message uppercase, and add more exclamation points
	// +optional
	shout bool,
	// Custom container for running the figlet tool
	// +optional
	figletContainer *dagger.Container,
) (string, error) {
	message := fmt.Sprintf("%s, %s!", greeting, name)
	if shout {
		message = strings.ToUpper(message) + "!!!!!"
	}
	if giant {
		// Run 'figlet' in a container to produce giant letters
		ctr := figletContainer
		if ctr == nil {
			ctr = defaultFigletContainer
		}
		return ctr.
			WithoutEntrypoint(). // clear the entrypoint to make sure 'figlet' is executed
			WithExec([]string{"figlet", message}).
			Stdout(ctx)
	}
	return message, nil
}

// Returns the files of the directory
func (m *Mod) Ls(
	ctx context.Context,
	// directory to list it's files
	dir *dagger.Directory,
	) (string, error) {
	return dag.Container().
		From("alpine:latest").
		WithMountedDirectory("/mnt", dir).
		WithWorkdir("/mnt").
		WithExec([]string{"ls", "-l", "."}).
		Stdout(ctx)
}

// Returns the operating system of the container
func (m *Mod) Os(
	ctx context.Context,
	// container to get is's OS
	ctr *dagger.Container,
	) (string, error) {
	return ctr.
		WithExec([]string{"cat", "/etc/os-release"}).
		Stdout(ctx)
}
