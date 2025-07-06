// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"dagger/go/internal/dagger"
)

type Mod struct {
	// +private
	Ctr *dagger.Container
}

// Mount a directory containing a Go module.
func (m *Go) Module(
	source *dagger.Directory,

	// +default="."
	path string,
) *Mod {
	return &Mod{
		Ctr: m.Container.WithMountedDirectory("/src", source).WithWorkdir("/src/" + path),
	}
}
