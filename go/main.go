// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"dagger/go/internal/dagger"
)

// Go
type Go struct {
	Ctr *dagger.Container
}

func New(
	// Alpine default does not come with gcc so cgo in unsupported by default.
	// +default="alpine"
	version string,

	// +optional
	from *dagger.Container,

	// +optional
	buildCache *dagger.CacheVolume,

	// +optional
	moduleCache *dagger.CacheVolume,
) *Go {
	if from == nil {
		from = dag.Container().From("golang:" + version)
	}

	if buildCache == nil {
		buildCache = dag.CacheVolume("github.com/z5labs/daggerverse/go:build")
	}
	if moduleCache == nil {
		moduleCache = dag.CacheVolume("github.com/z5labs/daggerverse/go:mod")
	}

	from = from.WithMountedCache("/root/.cache/go-build", buildCache)
	from = from.WithMountedCache("/go/pkg/mod", moduleCache)

	return &Go{
		Ctr: from,
	}
}

// WithWorkdir
func (m *Go) WithWorkdir(
	// +default="/src"
	path string,

	// The Go module source code.
	module *dagger.Directory,
) *Go {
	m.Ctr = m.Ctr.WithMountedDirectory(path, module).WithWorkdir(path)

	return m
}
