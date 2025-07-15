// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"io"
	"log/slog"
	"os"
)

//go:generate go run main.go

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{}))

	f, err := os.Create("generated.txt")
	if err != nil {
		log.Error("failed to create file", slog.Any("error", err))
		os.Exit(1)
	}
	defer f.Close()

	_, err = io.WriteString(f, "hello world")
	if err != nil {
		log.Error("failed write file contents", slog.Any("error", err))
		os.Exit(1)
	}
}
