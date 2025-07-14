// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"fmt"

	"dagger/protobuf/internal/dagger"
)

type generator struct {
	Name   string
	OutDir string
	Opts   []string
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
	source *dagger.Directory,

	proto []string,

	// Specify the directory in which to search for imports.
	// +default=["."]
	includePath []string,

	// +optional
	excludeWellKnownTypes bool,
) (*dagger.Directory, error) {
	args, outDirs := buildCompileArgs(
		excludeWellKnownTypes,
		includePath,
		p.Generators,
		proto,
	)

	c := p.Protobuf.Container.
		WithMountedDirectory("/src", source).
		WithWorkdir("/src")

	for _, outDir := range outDirs {
		c = c.WithExec([]string{"mkdir", "-p", outDir})
	}

	c = c.WithExec(args)

	dir := dag.Directory()
	for _, outDir := range outDirs {
		dir = dir.WithDirectory(outDir, c.Directory(outDir))
	}

	return dir, nil
}

func buildCompileArgs(
	excludeWellKnownTypes bool,
	includePath []string,
	generators []generator,
	proto []string,
) (args []string, outDirs []string) {
	args = []string{
		"protoc",
	}
	if !excludeWellKnownTypes {
		args = append(args, "-I", "/protobuf/include/")
	}

	for _, p := range includePath {
		args = append(args, "-I", p)
	}

	for _, g := range generators {
		outDirs = append(outDirs, g.OutDir)

		args = append(args, fmt.Sprintf("--%s_out", g.Name), g.OutDir)

		for _, opt := range g.Opts {
			args = append(args, fmt.Sprintf("--%s_opt", g.Name), opt)
		}
	}

	args = append(args, proto...)

	return
}
