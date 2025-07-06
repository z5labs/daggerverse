// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"
	"slices"
	"strings"

	"dagger/go/internal/dagger"
)

// Tidy
type Tidy struct {
	// +private
	Ctr *dagger.Container
}

// Ensure go.mod matches the source code in the module.
func (m *Go) Tidy(
	// The Go module source code.
	// +optional
	module *dagger.Directory,
) *Tidy {
	ctr := m.Ctr
	if module != nil {
		ctr = ctr.WithMountedDirectory("/src", module).
			WithWorkdir("/src")
	}

	return &Tidy{
		Ctr: ctr,
	}
}

func (t *Tidy) run(
	ctx context.Context,
	expectedReturnCode dagger.ReturnType,
	args ...string,
) (string, error) {
	cmd := []string{"go", "mod", "tidy"}
	cmd = append(cmd, args...)

	t.Ctr = t.Ctr.WithExec(cmd, dagger.ContainerWithExecOpts{
		Expect: expectedReturnCode,
	})

	return t.Ctr.Stdout(ctx)
}

// Apply neccessary changes for go.mod or go.sum and report changes made.
func (t *Tidy) Report(ctx context.Context) ([]string, error) {
	stdout, err := t.run(ctx, dagger.ReturnTypeSuccess, "-v")
	if err != nil {
		return nil, err
	}

	modules := slices.Collect(strings.Lines(stdout))
	return modules, nil
}

// Only report necessary changes for go.mod or go.sum.
func (t *Tidy) Diff(ctx context.Context) ([]string, error) {
	stdout, err := t.run(ctx, dagger.ReturnTypeAny, "-diff")
	if err != nil {
		return nil, err
	}

	changes := slices.Collect(strings.Lines(stdout))
	return changes, nil
}
