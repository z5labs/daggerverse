// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"dagger/go/internal/dagger"

	"github.com/sourcegraph/conc/pool"
)

type ContainerScanner interface {
	DaggerObject

	ScanContainer(ctx context.Context, ctr *dagger.Container) error
}

type Application struct {
	// +private
	Go *Go

	// +private
	Linter Linter

	// +private
	StaticAnalyzer StaticAnalyzer

	// +private
	ContainerScanner ContainerScanner
}

// A set of functions for working with a application written in Go.
func (m *Go) Application(
	// The Go module source code for the library.
	module *dagger.Directory,

	// +optional
	linter Linter,

	// +optional
	staticAnalyzer StaticAnalyzer,

	// +optional
	containerScanner ContainerScanner,
) *Application {
	return &Application{
		Go:               m.WithWorkdir("/src", module),
		Linter:           linter,
		StaticAnalyzer:   staticAnalyzer,
		ContainerScanner: containerScanner,
	}
}

// Run all continuous integration functions.
func (app *Application) Ci(
	ctx context.Context,

	imageRegistry string,

	imageName string,

	imageTags []string,

	registryUsername string,

	registrySecret *dagger.Secret,

	// +default="."
	mainPackagePath string,

	// +default=["-s", "-w"]
	ldflags []string,

	// +optional
	buildTags []string,

	// +default=true
	trimpath bool,

	// +optional
	enableCGO bool,

	// +default=["linux/amd64","linux/arm64"]
	platforms []dagger.Platform,
) error {
	err := app.Generate(ctx, "./...")
	if err != nil {
		return err
	}

	err = app.Tidy(ctx)
	if err != nil {
		return err
	}

	lintReport := app.Lint(ctx)

	coverageReport := app.Test("./...", true)

	err = app.StaticAnalysis(ctx, lintReport, coverageReport)
	if err != nil {
		return err
	}

	variants, err := app.Build(
		ctx,
		mainPackagePath,
		ldflags,
		buildTags,
		trimpath,
		enableCGO,
		platforms,
	)
	if err != nil {
		return err
	}

	err = app.Scan(ctx, variants)
	if err != nil {
		return err
	}

	results, err := app.Publish(
		ctx,
		imageRegistry,
		imageName,
		imageTags,
		variants,
		registryUsername,
		registrySecret,
	)
	if err != nil {
		return err
	}

	for _, result := range results {
		s := result.String()

		fmt.Println("published:", s)
	}

	return nil
}

// Run generate directives and validate no filesystem changes.
func (app *Application) Generate(
	ctx context.Context,

	// +default="./..."
	pkg string,
) error {
	entries, err := app.Go.Generate(pkg, nil).Diff(ctx)
	if err != nil {
		return err
	}

	if len(entries) > 0 {
		return fmt.Errorf("forgot to run go generate")
	}

	return nil
}

// Validate no necessary changes for go.mod or go.sum.
func (app *Application) Tidy(ctx context.Context) error {
	diff, err := app.Go.Tidy(nil).Diff(ctx)
	if err != nil {
		return err
	}

	if len(diff) != 0 {
		return errors.New("forgot to run go mod tidy")
	}

	return nil
}

// Lint source code.
func (app *Application) Lint(ctx context.Context) *dagger.File {
	if app.Linter == nil {
		return &dagger.File{}
	}

	return app.Linter.Lint(ctx, app.Go.Ctr)
}

// Run tests and return coverage report.
func (app *Application) Test(
	// +default="./..."
	pkg string,

	// +default=true
	race bool,
) *dagger.File {
	return app.Go.Test(pkg, nil, true).Coverage(Atomic)
}

// Perform static analysis.
func (app *Application) StaticAnalysis(
	ctx context.Context,
	lintReport *dagger.File,
	coverageReport *dagger.File,
) error {
	if app.StaticAnalyzer == nil {
		return nil
	}

	return app.StaticAnalyzer.StaticAnalysis(
		ctx,
		app.Go.Ctr,
		lintReport,
		coverageReport,
	)
}

// Produce container image(s) for application.
func (app *Application) Build(
	ctx context.Context,

	// +default="."
	pkg string,

	// +default=["-s", "-w"]
	ldflags []string,

	// +optional
	tags []string,

	// +default=true
	trimpath bool,

	// +optional
	enableCGO bool,

	// +default=["linux/amd64", "linux/arm64"]
	platforms []dagger.Platform,
) ([]*dagger.Container, error) {
	buildPool := pool.New().WithErrors().WithContext(ctx)

	containerCh := make(chan *dagger.Container, len(platforms))
	for _, platform := range platforms {
		buildPool.Go(func(ctx context.Context) error {
			b, err := app.Go.Build(pkg, nil, false, ldflags, tags, trimpath, enableCGO, platform)
			if err != nil {
				return err
			}

			f := b.Output()

			c := dag.Container(dagger.ContainerOpts{
				Platform: platform,
			}).
				WithFile("/main", f).
				WithEntrypoint([]string{"/main"})

			select {
			case <-ctx.Done():
				return nil
			case containerCh <- c:
			}

			return nil
		})
	}

	collectPool := pool.New().WithErrors().WithContext(ctx)
	collectPool.Go(func(ctx context.Context) error {
		defer close(containerCh)

		return buildPool.Wait()
	})

	var containers []*dagger.Container
	collectPool.Go(func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case c, ok := <-containerCh:
				if !ok {
					return nil
				}

				containers = append(containers, c)
			}
		}
	})

	err := collectPool.Wait()
	if err != nil {
		return nil, err
	}

	return containers, nil
}

// Create a service based on application container image.
func (app *Application) AsService(
	ctx context.Context,

	// +default="."
	pkg string,

	// +default=["-s", "-w"]
	ldflags []string,

	// +optional
	tags []string,

	// +default=true
	trimpath bool,

	// +optional
	enableCGO bool,

	// +optional
	args []string,

	// +default="linux/amd64"
	platform dagger.Platform,
) (*dagger.Service, error) {
	variants, err := app.Build(
		ctx,
		pkg,
		ldflags,
		tags,
		trimpath,
		enableCGO,
		[]dagger.Platform{
			platform,
		},
	)
	if err != nil {
		return nil, err
	}

	svc := variants[0].AsService(dagger.ContainerAsServiceOpts{
		Args:          args,
		UseEntrypoint: true,
	})

	return svc, nil
}

// Scan application container image(s).
func (app *Application) Scan(
	ctx context.Context,
	platformVariants []*dagger.Container,
) error {
	if app.ContainerScanner == nil {
		return nil
	}

	ep := pool.New().WithErrors().WithContext(ctx)

	for _, ctr := range platformVariants {
		ep.Go(func(ctx context.Context) error {
			return app.ContainerScanner.ScanContainer(ctx, ctr)
		})
	}

	return ep.Wait()
}

type PublishResult struct {
	Registry  string `json:"registry"`
	ImageName string `json:"image_name"`
	Tag       string `json:"tag"`
	Digest    string `json:"digest"`
}

// Format publish result as a string.
func (result *PublishResult) String() string {
	return fmt.Sprintf(
		"%s/%s:%s@%s",
		result.Registry,
		result.ImageName,
		result.Tag,
		result.Digest,
	)
}

// Format publish result as JSON.
func (result *PublishResult) Json() (string, error) {
	var buf bytes.Buffer

	enc := json.NewEncoder(&buf)
	err := enc.Encode(result)

	return buf.String(), err
}

// Publish application container image(s).
func (app *Application) Publish(
	ctx context.Context,

	registry string,

	imageName string,

	imageTags []string,

	platformVariants []*dagger.Container,

	registryUsername string,

	registrySecret *dagger.Secret,
) ([]*PublishResult, error) {
	publishPool := pool.New().WithErrors().WithContext(ctx)

	publishResultCh := make(chan *PublishResult, len(imageTags))
	for _, tag := range imageTags {
		addr := fmt.Sprintf("%s/%s:%s", registry, imageName, tag)

		publishPool.Go(func(ctx context.Context) error {
			fqin, err := dag.Container().
				WithRegistryAuth(registry, registryUsername, registrySecret).
				Publish(ctx, addr, dagger.ContainerPublishOpts{
					PlatformVariants: platformVariants,
				})
			if err != nil {
				return err
			}

			_, digest, found := strings.Cut(fqin, "@")
			if !found {
				return errors.New("malformed publish result: " + fqin)
			}

			result := &PublishResult{
				Registry:  registry,
				ImageName: imageName,
				Tag:       tag,
				Digest:    digest,
			}

			select {
			case <-ctx.Done():
				return nil
			case publishResultCh <- result:
			}

			return nil
		})
	}

	collectPool := pool.New().WithErrors().WithContext(ctx)
	collectPool.Go(func(ctx context.Context) error {
		defer close(publishResultCh)

		return publishPool.Wait()
	})

	results := make([]*PublishResult, len(imageTags))
	collectPool.Go(func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case result, ok := <-publishResultCh:
				if !ok {
					return nil
				}

				results = append(results, result)
			}
		}
	})

	err := collectPool.Wait()
	if err != nil {
		return nil, err
	}

	return results, nil
}
