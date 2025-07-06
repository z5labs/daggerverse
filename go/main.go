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
	// Base container configured for working with Go.
	Container *dagger.Container
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
	if buildCache == nil {
		buildCache = dag.CacheVolume("github.com/z5labs/daggerverse/go:build")
	}
	if moduleCache == nil {
		moduleCache = dag.CacheVolume("github.com/z5labs/daggerverse/go:mod")
	}

	if from == nil {
		from = dag.Container().
			From("golang:"+version).
			WithMountedCache("/root/.cache/go-build", buildCache).
			WithMountedCache("/go/pkg/mod", moduleCache)
	}

	return &Go{
		Container: from,
	}
}
