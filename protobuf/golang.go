// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import "fmt"

type Go struct {
	Protobuf *Protobuf
}

// Install the protoc-gen-go plugin.
func (m *Protobuf) Go(
	version string,
) (*Go, error) {
	archiveType := "tar.gz"
	if m.Plaform.OS == "windows" {
		archiveType = "zip"
	}

	plugin := dag.HTTP(fmt.Sprintf(
		"https://github.com/protocolbuffers/protobuf-go/releases/download/%s/protoc-gen-go.%s.%s.%s.%s",
		version,
		m.Plaform.OS,
		m.Plaform.Architecture,
		archiveType,
	))

	withPlugin, err := m.WithPlugin("protoc-gen-go", plugin)
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
