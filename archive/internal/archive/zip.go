// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package archive

import (
	"archive/zip"
	"context"
	"io"
	"log/slog"
	"os"
	"path"

	"dagger/archive/internal/dagger"

	"github.com/z5labs/sdk-go/try"
)

func ExtractZip(ctx context.Context, filename, out string) error {
	_, span := dagger.Tracer().Start(ctx, "archive.ExtractZip")
	defer span.End()

	log := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{}))

	err := os.MkdirAll(out, os.ModeDir)
	if err != nil {
		log.ErrorContext(ctx, "failed to create output directory", slog.Any("error", err))
		return err
	}

	zr, err := zip.OpenReader(filename)
	if err != nil {
		log.ErrorContext(ctx, "failed to open file", slog.Any("error", err))
		return err
	}
	defer zr.Close()

	for _, zipFile := range zr.File {
		if !zipFile.FileInfo().IsDir() {
			continue
		}

		err = os.MkdirAll(path.Join(out, zipFile.Name), os.ModeDir)
		if err != nil {
			log.ErrorContext(ctx, "failed to create output dir", slog.Any("error", err))
			return err
		}
	}

	for _, zipFile := range zr.File {
		if zipFile.FileInfo().IsDir() {
			continue
		}

		out, err := os.Create(path.Join(out, zipFile.Name))
		if err != nil {
			log.ErrorContext(ctx, "failed to create output file", slog.Any("error", err))
			return err
		}

		rc, err := zipFile.Open()
		if err != nil {
			log.ErrorContext(ctx, "failed to open zip content", slog.Any("error", err))
			return err
		}

		err = copyZipFile(out, rc)
		if err != nil {
			log.ErrorContext(ctx, "failed to write zip content to output directory", slog.Any("error", err))
			return err
		}
	}

	return nil
}

func copyZipFile(wc io.WriteCloser, rc io.ReadCloser) (err error) {
	defer try.Close(&err, wc)
	defer try.Close(&err, rc)

	_, err = io.Copy(wc, rc)
	return
}
