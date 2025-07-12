package main

import (
	"context"
)

type Archive struct{}

func New(ctx context.Context) *Archive {
	return &Archive{}
}
