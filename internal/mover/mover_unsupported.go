//go:build !darwin

package mover

import (
	"context"
	"errors"
)

type System struct{}

func (System) MoveAndReturn(ctx context.Context, distance int) error {
	return errors.New("stay is supported only on macOS")
}
