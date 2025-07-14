// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"
	"fmt"

	"dagger/protobuf-tests/internal/dagger"

	"github.com/sourcegraph/conc/pool"
)

type Go struct {
	// +private
	Protobuf *dagger.Protobuf
}

func (m *ProtobufTests) Go() *Go {
	return &Go{
		Protobuf: m.Protobuf.Go("v1.36.6").Protobuf(),
	}
}

func (g *Go) All(ctx context.Context) error {
	ep := pool.New().WithErrors().WithContext(ctx)

	ep.Go(g.ProtocTest)
	ep.Go(g.ProtocWithWellKnownTypesTest)
	ep.Go(g.ProtocWithoutWellKnownTypesTest)

	return ep.Wait()
}

func (g *Go) ProtocTest(ctx context.Context) error {
	dir := g.Protobuf.
		Protoc().
		Go("proto-go").
		Compile(
			dag.CurrentModule().Source().Directory("testdata"),
			[]string{"proto/message.proto"},
		).
		Directory("proto-go")

	entries, err := dir.Entries(ctx)
	if err != nil {
		return err
	}

	if len(entries) != 1 {
		return fmt.Errorf("expected only 1 generated file: %v", entries)
	}

	return nil
}

func (g *Go) ProtocWithWellKnownTypesTest(ctx context.Context) error {
	dir := g.Protobuf.
		Protoc().
		Go("proto-go").
		Compile(
			dag.CurrentModule().Source().Directory("testdata"),
			[]string{"proto/well_known.proto"},
		).
		Directory("proto-go")

	entries, err := dir.Entries(ctx)
	if err != nil {
		return err
	}

	if len(entries) != 1 {
		return fmt.Errorf("expected only 1 generated file: %v", entries)
	}

	return nil
}

func (g *Go) ProtocWithoutWellKnownTypesTest(ctx context.Context) error {
	dir := g.Protobuf.
		Protoc().
		Go("proto-go").
		Compile(
			dag.CurrentModule().Source().Directory("testdata"),
			[]string{"proto/well_known.proto"},
			dagger.ProtobufProtocCompileOpts{
				ExcludeWellKnownTypes: true,
			},
		).
		Directory("proto-go")

	entries, err := dir.Entries(ctx)
	if err != nil {
		return err
	}

	if len(entries) != 1 {
		return fmt.Errorf("expected only 1 generated file: %v", entries)
	}

	return nil
}

type GoGrpc struct {
	// +private
	Protobuf *dagger.Protobuf
}

func (m *ProtobufTests) GoGrpc() *GoGrpc {
	return &GoGrpc{
		Protobuf: m.Protobuf.Go("v1.36.6").Grpc().Protobuf(),
	}
}

func (g *GoGrpc) All(ctx context.Context) error {
	return nil
}
