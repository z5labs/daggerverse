// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"
	"fmt"
	"strings"

	"dagger/protobuf/internal/dagger"
)

type Go struct {
	Protobuf *Protobuf
}

// Install the protoc-gen-go plugin.
func (m *Protobuf) Go(
	version string,
) *Go {
	const name = "protoc-gen-go"

	archiveType := "tar.gz"
	if m.OS == "windows" {
		archiveType = "zip"
	}

	plugin := dag.HTTP(fmt.Sprintf(
		"https://github.com/protocolbuffers/protobuf-go/releases/download/%s/%s.%s.%s.%s.%s",
		version,
		name,
		version,
		m.OS,
		m.Arch,
		archiveType,
	))

	var dir *dagger.Directory
	switch archiveType {
	case "zip":
		dir = dag.Archive().Zip().Extract(plugin)
	case "tar.gz":
		dir = dag.Archive().Tar().Extract(plugin, dagger.ArchiveTarExtractOpts{
			Gzip: true,
		})
	}

	m = m.WithPlugin(name, dir.File(name))

	return &Go{
		Protobuf: m,
	}
}

// Generate Go code.
func (p *Protoc) Go(
	outDir string,

	// +optional
	opt []string,
) *Protoc {
	p.Generators = append(p.Generators, generator{
		Name:   "go",
		OutDir: outDir,
		Opts:   opt,
	})

	return p
}

// Install the protoc-gen-go-grpc plugin.
func (g *Go) Grpc(
	ctx context.Context,

	// +default="latest"
	version string,
) (*Go, error) {
	const name = "protoc-gen-go-grpc"

	c := dag.Container().
		From("golang:alpine").
		WithExec([]string{"go", "install", fmt.Sprintf("google.golang.org/grpc/cmd/%s@%s", name, version)})

	binPath, err := c.
		WithExec([]string{"which", name}).
		Stdout(ctx)
	if err != nil {
		return nil, err
	}

	g.Protobuf = g.Protobuf.WithPlugin(name, c.File(strings.TrimSpace(binPath)))

	return g, nil
}

// Generate GRPC Go code.
func (p *Protoc) GoGrpc(
	outDir string,

	// +optional
	opt []string,
) *Protoc {
	p.Generators = append(p.Generators, generator{
		Name:   "go-grpc",
		OutDir: outDir,
		Opts:   opt,
	})

	return p
}
