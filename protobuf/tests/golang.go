// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"dagger/protobuf-tests/internal/dagger"

	"github.com/sourcegraph/conc/pool"
)

type Go struct {
	// +private
	Protobuf *dagger.Protobuf
}

func (m *ProtobufTests) Go(
	// +default="v1.36.6"
	version string,
) *Go {
	return &Go{
		Protobuf: m.Protobuf.Go(version).Protobuf(),
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
		Go("proto-go", dagger.ProtobufProtocGoOpts{
			Opt: []string{"paths=source_relative"},
		}).
		Compile(
			dag.CurrentModule().Source().Directory("testdata"),
			[]string{"proto/message.proto"},
		).
		Directory("proto-go").
		Directory("proto")

	entries, err := dir.Entries(ctx)
	if err != nil {
		return err
	}

	if len(entries) != 1 {
		return fmt.Errorf("expected 1 generated file(s): %v", entries)
	}

	generatedFileName := entries[0]
	if generatedFileName != "message.pb.go" {
		return errors.New("unexpected file name for generated file: " + generatedFileName)
	}

	return nil
}

func (g *Go) ProtocWithWellKnownTypesTest(ctx context.Context) error {
	dir := g.Protobuf.
		Protoc().
		Go("proto-go", dagger.ProtobufProtocGoOpts{
			Opt: []string{"paths=source_relative"},
		}).
		Compile(
			dag.CurrentModule().Source().Directory("testdata"),
			[]string{"proto/well_known.proto"},
		).
		Directory("proto-go").
		Directory("proto")

	entries, err := dir.Entries(ctx)
	if err != nil {
		return err
	}

	if len(entries) != 1 {
		return fmt.Errorf("expected 1 generated file(s): %v", entries)
	}

	generatedFileName := entries[0]
	if generatedFileName != "well_known.pb.go" {
		return errors.New("unexpected file name for generated file: " + generatedFileName)
	}

	return nil
}

func (g *Go) ProtocWithoutWellKnownTypesTest(ctx context.Context) error {
	dir := g.Protobuf.
		Protoc().
		Go("proto-go", dagger.ProtobufProtocGoOpts{
			Opt: []string{"paths=source_relative"},
		}).
		Compile(
			dag.CurrentModule().Source().Directory("testdata"),
			[]string{"proto/well_known.proto"},
			dagger.ProtobufProtocCompileOpts{
				ExcludeWellKnownTypes: true,
			},
		).
		Directory("proto-go").
		Directory("proto")

	entries, err := dir.Entries(ctx)
	if err != nil {
		return err
	}

	if len(entries) != 1 {
		return fmt.Errorf("expected 1 generated file(s): %v", entries)
	}

	generatedFileName := entries[0]
	if generatedFileName != "well_known.pb.go" {
		return errors.New("unexpected file name for generated file: " + generatedFileName)
	}

	return nil
}

type GoGrpc struct {
	// +private
	Protobuf *dagger.Protobuf
}

func (m *ProtobufTests) GoGrpc(
	// +default="v1.36.6"
	goVersion string,

	// +default="latest"
	version string,
) *GoGrpc {
	return &GoGrpc{
		Protobuf: m.Protobuf.Go(goVersion).Grpc(dagger.ProtobufGoGrpcOpts{
			Version: version,
		}).Protobuf(),
	}
}

func (g *GoGrpc) All(ctx context.Context) error {
	ep := pool.New().WithErrors().WithContext(ctx)

	ep.Go(g.ProtocTest)
	ep.Go(g.ProtocWithWellKnownTypesTest)
	ep.Go(g.ProtocWithoutWellKnownTypesTest)

	return ep.Wait()
}

func (g *GoGrpc) ProtocTest(ctx context.Context) error {
	dir := g.Protobuf.
		Protoc().
		Go("proto-go", dagger.ProtobufProtocGoOpts{
			Opt: []string{"paths=source_relative"},
		}).
		GoGrpc("proto-go", dagger.ProtobufProtocGoGrpcOpts{
			Opt: []string{"paths=source_relative"},
		}).
		Compile(
			dag.CurrentModule().Source().Directory("testdata"),
			[]string{"proto/service.proto"},
		).
		Directory("proto-go").
		Directory("proto")

	entries, err := dir.Entries(ctx)
	if err != nil {
		return err
	}

	if len(entries) != 2 {
		return fmt.Errorf("expected 2 generated file(s): %v", entries)
	}

	if !slices.Contains(entries, "service.pb.go") {
		return errors.New("missing generated proto file: service.pb.go")
	}
	if !slices.Contains(entries, "service_grpc.pb.go") {
		return errors.New("missing generated grpc proto file: service_grpc.pb.go")
	}

	return nil
}

func (g *GoGrpc) ProtocWithWellKnownTypesTest(ctx context.Context) error {
	dir := g.Protobuf.
		Protoc().
		Go("proto-go", dagger.ProtobufProtocGoOpts{
			Opt: []string{"paths=source_relative"},
		}).
		GoGrpc("proto-go", dagger.ProtobufProtocGoGrpcOpts{
			Opt: []string{"paths=source_relative"},
		}).
		Compile(
			dag.CurrentModule().Source().Directory("testdata"),
			[]string{"proto/well_known_service.proto"},
			dagger.ProtobufProtocCompileOpts{
				ExcludeWellKnownTypes: false,
			},
		).
		Directory("proto-go").
		Directory("proto")

	entries, err := dir.Entries(ctx)
	if err != nil {
		return err
	}

	if len(entries) != 2 {
		return fmt.Errorf("expected 2 generated file(s): %v", entries)
	}

	if !slices.Contains(entries, "well_known_service.pb.go") {
		return errors.New("missing generated proto file: well_known_service.pb.go")
	}
	if !slices.Contains(entries, "well_known_service_grpc.pb.go") {
		return errors.New("missing generated grpc proto file: well_known_service_grpc.pb.go")
	}

	return nil
}

func (g *GoGrpc) ProtocWithoutWellKnownTypesTest(ctx context.Context) error {
	dir := g.Protobuf.
		Protoc().
		Go("proto-go", dagger.ProtobufProtocGoOpts{
			Opt: []string{"paths=source_relative"},
		}).
		GoGrpc("proto-go", dagger.ProtobufProtocGoGrpcOpts{
			Opt: []string{"paths=source_relative"},
		}).
		Compile(
			dag.CurrentModule().Source().Directory("testdata"),
			[]string{"proto/well_known_service.proto"},
			dagger.ProtobufProtocCompileOpts{
				ExcludeWellKnownTypes: true,
			},
		).
		Directory("proto-go").
		Directory("proto")

	entries, err := dir.Entries(ctx)
	if err != nil {
		return err
	}

	if len(entries) != 2 {
		return fmt.Errorf("expected 2 generated file(s): %v", entries)
	}

	if !slices.Contains(entries, "well_known_service.pb.go") {
		return errors.New("missing generated proto file: well_known_service.pb.go")
	}
	if !slices.Contains(entries, "well_known_service_grpc.pb.go") {
		return errors.New("missing generated grpc proto file: well_known_service_grpc.pb.go")
	}

	return nil
}
