//go:build !darwin && !windows

package mover

import (
	"context"
	"errors"
)

type System struct{}

func (System) MoveAndReturn(ctx context.Context, distance int) error {
	return errors.New("keep is supported only on macOS and Windows")
}
