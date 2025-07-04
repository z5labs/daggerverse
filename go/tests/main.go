package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"dagger/tests/internal/dagger"

	"github.com/sourcegraph/conc/pool"
	"golang.org/x/tools/cover"
)

type Tests struct {
	// +private
	Version string
}

func New(
	version string,
) *Tests {
	return &Tests{
		Version: version,
	}
}

func (m *Tests) All(ctx context.Context) error {
	ep := pool.New().WithErrors().WithContext(ctx)

	ep.Go(m.BuildWithDefaultsTest)
	ep.Go(m.BuildWithOptionsTest)

	ep.Go(m.TestCoverageTest)

	return ep.Wait()
}

func (m *Tests) BuildWithDefaultsTest(ctx context.Context) error {
	g := dag.Go(m.Version, dag.CurrentModule().Source().Directory("src/buildoutput"))

	f := g.Build(".").Output()

	n, err := f.Size(ctx)
	if err != nil {
		return err
	}

	if n == 0 {
		return errors.New("expected non-empty output file")
	}

	return nil
}

func (m *Tests) BuildWithOptionsTest(ctx context.Context) error {
	g := dag.Go(m.Version, dag.CurrentModule().Source().Directory("src/buildoutput"))

	f := g.Build(".").Output()

	withDebugSize, err := f.Size(ctx)
	if err != nil {
		return err
	}

	if withDebugSize == 0 {
		return errors.New("expected non-empty output file")
	}

	f = g.Build(".", dagger.GoBuildOpts{
		Ldflags: []string{"-s", "-w"},
	}).Output()

	withoutDebugSize, err := f.Size(ctx)
	if err != nil {
		return err
	}

	if withoutDebugSize == 0 {
		return errors.New("expected non-empty output file")
	}

	if withDebugSize < withoutDebugSize {
		return fmt.Errorf(
			"expected output file with debug symbols to be larger than without: %d:%d",
			withDebugSize,
			withoutDebugSize,
		)
	}

	return nil
}

func (m *Tests) TestCoverageTest(ctx context.Context) error {
	g := dag.Go(m.Version, dag.CurrentModule().Source().Directory("src/testcoverage"))

	f := g.Test("./...", dagger.GoTestOpts{
		Race: true,
	}).Coverage()

	contents, err := f.Contents(ctx)
	if err != nil {
		return err
	}

	profiles, err := cover.ParseProfilesFromReader(strings.NewReader(contents))
	if err != nil {
		return err
	}

	if len(profiles) != 1 {
		return fmt.Errorf("expected only 1 coverage profile but received: %d", len(profiles))
	}

	profile := profiles[0]
	if profile.Mode != dagger.GoCoverageModeAtomic.Value() {
		return errors.New("expected race option to set coverage mode to atomic: " + profile.Mode)
	}

	if len(profile.Blocks) != 1 {
		return fmt.Errorf("expected only 1 profile block in coverage report: %d", len(profile.Blocks))
	}

	return nil
}
