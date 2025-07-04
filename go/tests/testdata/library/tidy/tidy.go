// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package tidy

import (
	"context"

	"go.opentelemetry.io/otel"
)

func HelloWorld() {
	_, span := otel.Tracer("tidy").Start(context.Background(), "HelloWorld")
	defer span.End()
}
