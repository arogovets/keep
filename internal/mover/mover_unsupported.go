//go:build !darwin

package mover

import "errors"

type System struct{}

func (System) MoveAndReturn(distance int) error {
	return errors.New("stay is supported only on macOS")
}
