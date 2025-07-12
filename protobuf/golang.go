// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"dagger/protobuf/internal/dagger"
	"fmt"
)

type Go struct {
	Protobuf *Protobuf
}

// Install the protoc-gen-go plugin.
func (m *Protobuf) Go(
	version string,
) (*Go, error) {
	archiveType := "tar.gz"
	if m.OS == "windows" {
		archiveType = "zip"
	}

	plugin := dag.HTTP(fmt.Sprintf(
		"https://github.com/protocolbuffers/protobuf-go/releases/download/%s/protoc-gen-go.%s.%s.%s.%s",
		version,
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

	withPlugin, err := m.WithPlugin("protoc-gen-go", dir.File("protoc-gen-go"))
	if err != nil {
		return nil, err
	}

	return &Go{
		Protobuf: withPlugin,
	}, nil
}

// Generate Go code.
func (p *Protoc) Go(
	outDir string,
	opt []string,
) *Protoc {
	p.Generators = append(p.Generators, generator{
		name:   "go",
		outDir: outDir,
		opts:   opt,
	})

	return p
}

// Install the protoc-gen-go-grpc plugin.
func (g *Go) Grpc() *Go {
	return g
}

// Generate GRPC Go code.
func (p *Protoc) GoGrpc(
	outDir string,
	opt []string,
) *Protoc {
	p.Generators = append(p.Generators, generator{
		name:   "go-grpc",
		outDir: outDir,
		opts:   opt,
	})

	return p
}
