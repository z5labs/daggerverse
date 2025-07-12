// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"
	"fmt"

	"dagger/protobuf/internal/dagger"
)

type generator struct {
	name   string
	outDir string
	opts   []string
}

type Protoc struct {
	// +private
	Protobuf *Protobuf

	// +private
	Generators []generator
}

// Run protoc.
func (m *Protobuf) Protoc() *Protoc {
	return &Protoc{
		Protobuf: m,
	}
}

// Compile protocol buffer definitions.
func (p *Protoc) Compile(
	ctx context.Context,

	source *dagger.Directory,

	proto []string,

	// Specify the directory in which to search for imports.
	// +default=["."]
	includePath []string,

	// +default=true
	includeWellKnownTypes bool,
) (string, error) {
	args := []string{
		"protoc",
	}
	if includeWellKnownTypes {
		args = append(args, "-I") // TODO
	}

	for _, p := range includePath {
		args = append(args, "-I", p)
	}

	for _, g := range p.Generators {
		args = append(args, fmt.Sprintf("--%s_out", g.name), g.outDir)

		for _, opt := range g.opts {
			args = append(args, fmt.Sprintf("--%s_opt", g.name), opt)
		}
	}

	for _, p := range proto {
		args = append(args, p)
	}

	return p.Protobuf.Container.
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec(args).
		Stdout(ctx)
}
