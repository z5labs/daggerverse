// A generated module for Trivy functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"fmt"
	"strings"

	"dagger/trivy/internal/dagger"
)

type Trivy struct {
	Ctr *dagger.Container

	// +private
	Severity []string

	// +private
	Format string
}

func New(
	// +default="latest"
	imageTag string,

	// +default=["UNKNOWN", "LOW", "MEDIUM", "HIGH", "CRITICAL"]
	severity []string,

	// +default="table"
	format string,
) *Trivy {
	ctr := dag.Container().
		From("aquasec/trivy:"+imageTag).
		WithMountedCache("/root/.cache/trivy", dag.CacheVolume("github.com/z5labs/daggerverse/trivy"))

	return &Trivy{
		Ctr:      ctr,
		Severity: severity,
		Format:   format,
	}
}

// Scan a container for vulnerabilities.
func (m *Trivy) ScanContainer(ctx context.Context, ctr *dagger.Container) error {
	platform, err := ctr.Platform(ctx)
	if err != nil {
		return err
	}

	stdout, err := m.Ctr.
		WithMountedFile("/scan/ctr.tar", ctr.AsTarball()).
		WithExec([]string{
			"trivy",
			"image",
			"--quiet",
			"--platform",
			string(platform),
			"--severity",
			strings.Join(m.Severity, ","),
			"--format",
			m.Format,
			"--input",
			"/scan/ctr.tar",
		}).
		Stdout(ctx)
	if err != nil {
		return err
	}

	fmt.Println(stdout)

	return nil
}
