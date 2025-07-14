// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"

	"dagger/protobuf/internal/dagger"

	"github.com/containerd/platforms"
)

type Protobuf struct {
	Container *dagger.Container

	// +private
	OS string

	// +private
	Arch string
}

func New(
	ctx context.Context,
	version string,
	// +optional
	platform dagger.Platform,
	// +optional
	from *dagger.Container,
) (*Protobuf, error) {
	if platform == "" {
		platform = dagger.Platform(platforms.DefaultString())
	}

	p, err := platforms.Parse(string(platform))
	if err != nil {
		return nil, err
	}

	if from == nil {
		c, err := newContainer(version, p)
		if err != nil {
			return nil, err
		}

		from = c
	}

	return &Protobuf{
		Container: from,
		OS:        p.OS,
		Arch:      p.Architecture,
	}, nil
}

func newContainer(version string, platform platforms.Platform) (*dagger.Container, error) {
	arch, err := mapToProtoArch(platform.Architecture)
	if err != nil {
		return nil, err
	}

	protoc := dag.HTTP(fmt.Sprintf(
		"https://github.com/protocolbuffers/protobuf/releases/download/%s/protoc-%s-%s-%s.zip",
		version,
		strings.TrimPrefix(version, "v"),
		platform.OS,
		arch,
	))

	dir := dag.Archive().Zip().Extract(protoc)

	c := dag.Container().
		From("alpine").
		WithDirectory("/protobuf", dir).
		WithEnvVariable("PATH", "/protobuf/bin/:${PATH}", dagger.ContainerWithEnvVariableOpts{
			Expand: true,
		}).
		WithExec([]string{"chmod", "+x", "/protobuf/bin/protoc"})

	return c, nil
}

func mapToProtoArch(arch string) (string, error) {
	mapping := map[string]string{
		"amd64":   "x86_64",
		"ppc64le": "ppcle_64",
	}

	mappedArch, ok := mapping[arch]
	if !ok {
		return "", errors.New("protoc does not support: " + arch)

	}

	return mappedArch, nil
}

// Register the given binary file as a protoc plugin.
func (m *Protobuf) WithPlugin(name string, bin *dagger.File) (*Protobuf, error) {
	if !strings.HasPrefix(name, "protoc-gen-") {
		return nil, errors.New("plugin name must start with: protoc-gen-")
	}

	binPath := path.Join("/protobuf/bin", name)

	m.Container = m.Container.
		WithFile(binPath, bin).
		WithExec([]string{"chmod", "+x", binPath})

	return m, nil
}

// Copy well-known types, protoc, and plugins to the provided container.
func (m *Protobuf) CopyTo(container *dagger.Container) *dagger.Container {
	return container
}
