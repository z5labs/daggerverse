// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

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

type Test struct {
	// +private
	Go *dagger.Go
}

func (m *Tests) Test() *Test {
	return &Test{
		Go: m.Go,
	}
}

func (t *Test) All(ctx context.Context) error {
	ep := pool.New().WithErrors().WithContext(ctx)

	ep.Go(t.CoverageTest)

	return ep.Wait()
}

func (t *Test) CoverageTest(ctx context.Context) error {
	f := t.Go.Module(dag.CurrentModule().Source().Directory("testdata/testcoverage")).
		Test("./...", dagger.GoModTestOpts{
			Race: true,
		}).
		Coverage()

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
