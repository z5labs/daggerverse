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

// Run tests within a Go module.
func (m *Mod) Test(
	pkg string,

	// +optional
	race bool,
) *Test {
	ctr := m.Ctr
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

// Retrieve a coverage profile from running the tests.
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
