//go:build !darwin

package activity

import (
	"errors"
	"time"
)

type System struct{}

func (System) LastInputAge() (time.Duration, error) {
	return 0, errors.New("stay is supported only on macOS")
}

func (System) AccessibilityTrusted() bool {
	return false
}
