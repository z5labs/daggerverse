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
	// +private
	Ctr *dagger.Container

	// +private
	Source *dagger.Directory
}

func New(
	version string,

	source *dagger.Directory,

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

	from = from.
		WithWorkdir("/src").
		WithMountedDirectory("/src", source)

	if buildCache != nil {
		from = from.WithMountedCache("/root/.cache/go-build", buildCache)
	}
	if moduleCache != nil {
		from = from.WithMountedCache("/go/pkg/mod", moduleCache)
	}

	return &Go{
		Ctr:    from,
		Source: source,
	}
}
