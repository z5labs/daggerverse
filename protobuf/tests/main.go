package main

import (
	"context"

	"dagger/protobuf-tests/internal/dagger"

	"github.com/sourcegraph/conc/pool"
)

type ProtobufTests struct {
	// +private
	Protobuf *dagger.Protobuf

	// +private
	GoVersion string

	// +private
	GoGrpcVersion string
}

func New(
	// +default="v31.1"
	version string,

	// +default="v1.36.6"
	goVersion string,

	// +default="latest"
	goGrpcVersion string,
) *ProtobufTests {
	return &ProtobufTests{
		Protobuf:      dag.Protobuf(version),
		GoVersion:     goVersion,
		GoGrpcVersion: goGrpcVersion,
	}
}

func (m *ProtobufTests) All(ctx context.Context) error {
	ep := pool.New().WithErrors().WithContext(ctx)

	ep.Go(m.Go(m.GoVersion).All)
	ep.Go(m.GoGrpc(m.GoVersion, m.GoGrpcVersion).All)

	return ep.Wait()
}
