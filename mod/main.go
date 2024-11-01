// A module to support the Puzzle dagger techlab.
//
// The functions are used inside the hands-on lab: https://dagger-techlab.puzzle.ch/

package main

import (
	"context"
	"dagger/mod/internal/dagger"
	"errors"
	"strings"
)

type DaggerTechlabModule struct{}

type LintRun struct {
	// +private
	Source *dagger.Directory
}

// Say hello to the world!
// Calls external module Hello https://github.com/shykes/hello
func (m *DaggerTechlabModule) Hello(
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
	) (string, error) {
    return dag.Hello().
        Hello(ctx, dagger.HelloHelloOpts{Greeting: greeting, Name: name, Giant: giant, Shout: shout})
}

// Returns the files of the directory
func (m *DaggerTechlabModule) Ls(
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
func (m *DaggerTechlabModule) Os(
	ctx context.Context,
	// container to get it's OS
	ctr *dagger.Container,
	) (string, error) {
	return ctr.
		WithExec([]string{"cat", "/etc/os-release"}).
		Stdout(ctx)
}

// Returns the answer to everything when the password is right
func (m *DaggerTechlabModule) Unlock(
	ctx context.Context,
	password *dagger.Secret,
	) (string, error) {
		passwordText, err := password.Plaintext(ctx)
		if err != nil {
			return "", err
		}
	passwordTextClean := strings.TrimSpace(passwordText)
	if passwordTextClean == "MySuperSecret" {
		return "You unlocked the secret. The answer is 42!", nil
	}
	return "", errors.New("Nice try ;-) Provide right password to unlock the secret.")
}

// Returns a service that runs an OpenSSH server
// Calls external module OpensshServer https://github.com/sagikazarmark/daggerverse/tree/main/openssh-server
func (m *DaggerTechlabModule) SshService(
	// +optional
	// +default=22
    port int,
    ) *dagger.Service {
	return dag.OpensshServer().
	    Service(dagger.OpensshServerServiceOpts{Port: port})
}

// Lint a Python codebase
// Calls external module Ruff https://github.com/dagger/dagger/tree/main/modules/ruff
func (m *DaggerTechlabModule) Lint(
	source *dagger.Directory,
) *LintRun {
	return &LintRun{
		Source: source,
	}
}

// Returns a JSON report file for this run
func (run LintRun) Report() *dagger.File {
	return dag.Ruff().
	    Lint(run.Source).
	    Report()
}

// Build a Wolfi Linux container
// Calls external module Wolfi https://github.com/shykes/daggerverse/tree/main/wolfi
func (m *DaggerTechlabModule) Wolfi() *dagger.Container {
    return dag.Wolfi().
        Container()
}
