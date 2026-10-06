//go:build !darwin && !windows

package activity

import (
	"errors"
	"time"
)

type System struct{}

func (System) LastInputAge() (time.Duration, error) {
	return 0, errors.New("keep is supported only on macOS and Windows")
}

func (System) AccessibilityTrusted() bool {
	return false
}
