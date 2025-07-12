// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"

	"dagger/protobuf-tests/internal/dagger"
)

type Go struct {
	// +private
	Protobuf *dagger.Protobuf
}

func (m *ProtobufTests) Go() *Go {
	return &Go{
		Protobuf: m.Protobuf.Go().Protobuf(),
	}
}

func (g *Go) All(ctx context.Context) error {
	return nil
}

type GoGrpc struct {
	// +private
	Protobuf *dagger.Protobuf
}

func (m *ProtobufTests) GoGrpc() *GoGrpc {
	return &GoGrpc{
		Protobuf: m.Protobuf.Go().Grpc().Protobuf(),
	}
}

func (g *GoGrpc) All(ctx context.Context) error {
	return nil
}
