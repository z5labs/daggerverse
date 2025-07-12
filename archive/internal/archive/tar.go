// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package archive

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"io"
	"log/slog"
	"os"
	"path"

	"dagger/archive/internal/dagger"

	"github.com/z5labs/sdk-go/try"
)

func ExtractTar(ctx context.Context, filename, out string, isGziped bool) error {
	_, span := dagger.Tracer().Start(ctx, "archive.ExtractTar")
	defer span.End()

	log := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{}))

	err := os.MkdirAll(out, os.ModeDir)
	if err != nil {
		log.ErrorContext(ctx, "failed to create output directory", slog.Any("error", err))
		return err
	}

	f, err := os.Open(filename)
	if err != nil {
		log.ErrorContext(ctx, "failed to open file", slog.Any("error", err))
		return err
	}
	defer f.Close()

	var stream io.Reader = f
	if isGziped {
		stream, err = gzip.NewReader(stream)
		if err != nil {
			log.ErrorContext(ctx, "failed to create gzip reader", slog.Any("error", err))
			return err
		}
	}

	tr := tar.NewReader(stream)
	for {
		h, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.ErrorContext(ctx, "failed to get header", slog.Any("error", err))
			return err
		}

		if h.FileInfo().IsDir() {
			continue
		}

		dir, _ := path.Split(h.Name)
		if dir != "" {
			err = os.MkdirAll(path.Join(out, dir), os.ModeDir)
			if err != nil {
				log.ErrorContext(ctx, "failed to create output sub dir", slog.Any("error", err))
				return err
			}
		}

		f, err := os.Create(path.Join(out, h.Name))
		if err != nil {
			return err
		}

		err = copyTarFile(f, tr)
		if err != nil {
			log.ErrorContext(ctx, "failed to write tar content to output directory", slog.Any("error", err))
			return err
		}
	}
}

func copyTarFile(wc io.WriteCloser, r io.Reader) (err error) {
	defer try.Close(&err, wc)

	_, err = io.Copy(wc, r)
	return
}
