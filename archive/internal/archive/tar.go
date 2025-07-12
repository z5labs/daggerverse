// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package archive

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"log/slog"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/z5labs/sdk-go/try"
)

func tarCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tar",
		Short: "Utilities for working with TAR archives.",
	}

	cmd.AddCommand(
		tarExtractCommand(),
	)

	return cmd
}

func tarExtractCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extract FILE DIR",
		Short: "Extract a TAR archive to a specified directory.",
		Args:  cobra.ExactArgs(2),
		RunE:  extractTar,
	}

	cmd.Flags().Bool("gzip", false, "Enable gzip decompression.")

	return cmd
}

func extractTar(cmd *cobra.Command, args []string) error {
	log := slog.New(slog.NewJSONHandler(cmd.ErrOrStderr(), &slog.HandlerOptions{}))

	err := os.MkdirAll(args[1], os.ModeDir)
	if err != nil {
		log.ErrorContext(cmd.Context(), "failed to create output directory", slog.Any("error", err))
		return err
	}

	f, err := os.Open(args[0])
	if err != nil {
		log.ErrorContext(cmd.Context(), "failed to open file", slog.Any("error", err))
		return err
	}
	defer f.Close()

	var stream io.Reader = f
	isGziped, _ := cmd.Flags().GetBool("gzip")
	if isGziped {
		stream, err = gzip.NewReader(stream)
		if err != nil {
			log.ErrorContext(cmd.Context(), "failed to create gzip reader", slog.Any("error", err))
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
			log.ErrorContext(cmd.Context(), "failed to get header", slog.Any("error", err))
			return err
		}

		if h.FileInfo().IsDir() {
			continue
		}

		dir, _ := path.Split(h.Name)
		if dir != "" {
			err = os.MkdirAll(path.Join(args[1], dir), os.ModeDir)
			if err != nil {
				log.ErrorContext(cmd.Context(), "failed to create output sub dir", slog.Any("error", err))
				return err
			}
		}

		f, err := os.Create(path.Join(args[1], h.Name))
		if err != nil {
			return err
		}

		err = copyTarFile(f, tr)
		if err != nil {
			log.ErrorContext(cmd.Context(), "failed to write tar content to output directory", slog.Any("error", err))
			return err
		}
	}
}

func copyTarFile(wc io.WriteCloser, r io.Reader) (err error) {
	defer try.Close(&err, wc)

	_, err = io.Copy(wc, r)
	return
}
