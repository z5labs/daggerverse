package main

import (
	"context"

	"dagger/protobuf-tests/internal/dagger"

	"github.com/sourcegraph/conc/pool"
)

type ProtobufTests struct {
	// +private
	Protobuf *dagger.Protobuf
}

func New(
	// +default="v31.1"
	version string,
) *ProtobufTests {
	return &ProtobufTests{
		Protobuf: dag.Protobuf(version),
	}
}

func (m *ProtobufTests) All(ctx context.Context) error {
	ep := pool.New().WithErrors().WithContext(ctx)

	return ep.Wait()
}
