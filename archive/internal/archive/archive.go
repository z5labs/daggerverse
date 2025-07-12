// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package archive

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
)

func Main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer cancel()

	log := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{}))

	cmd := archiveCommand()

	err := cmd.ExecuteContext(ctx)
	if err != nil {
		log.ErrorContext(ctx, "unexpected error", slog.Any("error", err))
		os.Exit(1)
	}
}

func archiveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "archive",
	}

	cmd.AddCommand(
		zipCommand(),
		tarCommand(),
	)

	return cmd
}
