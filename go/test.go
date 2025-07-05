// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import "dagger/go/internal/dagger"

type Test struct {
	// +private
	Ctr *dagger.Container

	// +private
	Pkg string

	// +private
	Race bool
}

// Test
func (m *Go) Test(
	pkg string,

	// The Go module source code.
	// +optional
	module *dagger.Directory,

	// +optional
	race bool,
) *Test {
	ctr := m.Ctr
	if module != nil {
		ctr = ctr.WithMountedDirectory("/src", module).
			WithWorkdir("/src")
	}

	if race {
		ctr = ctr.WithEnvVariable("CGO_ENABLED", "1")
	}

	return &Test{
		Ctr:  ctr,
		Pkg:  pkg,
		Race: race,
	}
}

type CoverageMode string

const (
	Atomic CoverageMode = "atomic"
	Count  CoverageMode = "count"
	Set    CoverageMode = "set"
)

// Coverage
func (t *Test) Coverage(
	// +default="set"
	mode CoverageMode,
) *dagger.File {
	args := []string{
		"go",
		"test",
		"-coverprofile",
		"/tmp/cover.out",
	}

	if t.Race {
		args = append(args, "-race")
		mode = Atomic
	}

	args = append(args, "-covermode", string(mode))

	args = append(args, t.Pkg)

	out := dag.Directory()

	return t.Ctr.
		WithDirectory("/tmp", out).
		WithExec(args).
		File("/tmp/cover.out")
}
