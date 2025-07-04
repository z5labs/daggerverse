// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"strings"

	"dagger/go/internal/dagger"

	"github.com/containerd/platforms"
)

type Build struct {
	// +private
	Ctr *dagger.Container

	// +private
	Pkg string

	// +private
	Race bool

	// +private
	Ldflags []string

	// +private
	Tags []string

	// +private
	Trimpath bool
}

// Build
func (m *Go) Build(
	pkg string,

	// The Go module source code.
	// +optional
	module *dagger.Directory,

	// +optional
	race bool,

	// +optional
	ldflags []string,

	// +optional
	tags []string,

	// +optional
	trimpath bool,

	// +optional
	enableCGO bool,

	// +optional
	platform dagger.Platform,
) (*Build, error) {
	ctr := m.Ctr

	if module != nil {
		ctr = ctr.WithMountedDirectory("/src", module).
			WithWorkdir("/src")
	}

	if platform != "" {
		p, err := platforms.Parse(string(platform))
		if err != nil {
			return nil, err
		}

		ctr = ctr.
			WithEnvVariable("GOARCH", p.Architecture).
			WithEnvVariable("GOOS", p.OS)
	}

	cgoEnabled := "0"
	if enableCGO {
		cgoEnabled = "1"
	}

	ctr = ctr.WithEnvVariable("CGO_ENABLED", cgoEnabled)

	b := &Build{
		Ctr:      ctr,
		Pkg:      pkg,
		Race:     race,
		Ldflags:  ldflags,
		Tags:     tags,
		Trimpath: trimpath,
	}

	return b, nil
}

// Output
func (b *Build) Output() *dagger.File {
	args := []string{
		"go",
		"build",
		"-o",
		"/tmp/main",
	}

	if b.Race {
		args = append(args, "-race")
	}

	if b.Trimpath {
		args = append(args, "-trimpath")
	}

	if len(b.Ldflags) > 0 {
		args = append(args, "-ldflags", strings.Join(b.Ldflags, " "))
	}

	if len(b.Tags) > 0 {
		args = append(args, "-tags", strings.Join(b.Tags, ","))
	}

	args = append(args, b.Pkg)

	out := dag.Directory()

	return b.Ctr.
		WithDirectory("/tmp", out).
		WithExec(args).
		File("/tmp/main")
}
