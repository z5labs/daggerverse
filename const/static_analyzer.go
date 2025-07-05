// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"
	"errors"

	"dagger/const/internal/dagger"
)

type GoStaticAnalyzer struct {
	// +private
	AlwaysFail bool
}

func (*Const) GoStaticAnalyzer(
	// +optional
	alwaysFail bool,
) *GoStaticAnalyzer {
	return &GoStaticAnalyzer{
		AlwaysFail: alwaysFail,
	}
}

func (g *GoStaticAnalyzer) StaticAnalysis(
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

	if g.AlwaysFail {
		return errors.New("static analysis failed")
	}

	return nil
}
