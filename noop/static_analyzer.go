// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"

	"dagger/noop/internal/dagger"
)

type GoStaticAnalyzer struct{}

func (*Noop) GoStaticAnalyzer() *GoStaticAnalyzer {
	return &GoStaticAnalyzer{}
}

func (*GoStaticAnalyzer) StaticAnalysis(
	ctx context.Context,
	ctr *dagger.Container,
	lintReport *dagger.File,
	coverageReport *dagger.File,
) error {
	_, err := lintReport.Sync(ctx)
	if err != nil {
		return err
	}

	_, err = coverageReport.Sync(ctx)
	if err != nil {
		return err
	}

	return nil
}
