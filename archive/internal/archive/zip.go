// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package archive

import (
	"archive/zip"
	"io"
	"log/slog"
	"os"
	"path"

	"github.com/spf13/cobra"
)

func zipCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zip",
		Short: "Utilities for working with ZIP archives.",
	}

	cmd.AddCommand(
		zipExtractCommand(),
	)

	return cmd
}

func zipExtractCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extract FILE DIR",
		Short: "Extract a ZIP archive to a specified directory.",
		Args:  cobra.ExactArgs(2),
		RunE:  extractZip,
	}

	return cmd
}

func extractZip(cmd *cobra.Command, args []string) error {
	log := slog.New(slog.NewJSONHandler(cmd.ErrOrStderr(), &slog.HandlerOptions{}))

	err := os.MkdirAll(args[1], os.ModeDir)
	if err != nil {
		log.ErrorContext(cmd.Context(), "failed to create output directory", slog.Any("error", err))
		return err
	}

	zr, err := zip.OpenReader(args[0])
	if err != nil {
		log.ErrorContext(cmd.Context(), "failed to open file", slog.Any("error", err))
		return err
	}
	defer zr.Close()

	for _, zipFile := range zr.File {
		if !zipFile.FileInfo().IsDir() {
			continue
		}

		err = os.MkdirAll(path.Join(args[1], zipFile.Name), os.ModeDir)
		if err != nil {
			log.ErrorContext(cmd.Context(), "failed to create output dir", slog.Any("error", err))
			return err
		}
	}

	for _, zipFile := range zr.File {
		if zipFile.FileInfo().IsDir() {
			continue
		}

		out, err := os.Create(path.Join(args[1], zipFile.Name))
		if err != nil {
			log.ErrorContext(cmd.Context(), "failed to create output file", slog.Any("error", err))
			return err
		}

		rc, err := zipFile.Open()
		if err != nil {
			log.ErrorContext(cmd.Context(), "failed to open zip content", slog.Any("error", err))
			return err
		}

		_, err = io.Copy(out, rc)
		if err != nil {
			log.ErrorContext(cmd.Context(), "failed to write zip content to output directory", slog.Any("error", err))
			return err
		}
	}

	return nil
}
