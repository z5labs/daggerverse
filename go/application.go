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
	Library *Library

	// +private
	MainPackagePath string

	// +private
	ContainerScanner ContainerScanner
}

// A set of functions for working with a application written in Go.
func (lib *Library) Application(
	// Path to package main for the application.
	mainPackagePath string,

	// Specify a tool for scanning built application container before publishing.
	// +optional
	containerScanner ContainerScanner,
) *Application {
	if containerScanner == nil {
		containerScanner = dag.Trivy()
	}

	return &Application{
		Library:          lib,
		MainPackagePath:  mainPackagePath,
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
	err := app.Library.Ci(ctx)
	if err != nil {
		return err
	}

	variants, err := app.Build(
		ctx,
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

// Produce container image(s) for application.
func (app *Application) Build(
	ctx context.Context,

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
			b, err := app.Library.Module.Build(
				app.MainPackagePath,
				false,
				ldflags,
				tags,
				trimpath,
				enableCGO,
				platform,
			)
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
