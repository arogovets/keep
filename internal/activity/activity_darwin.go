//go:build darwin

package activity

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework ApplicationServices
#include <ApplicationServices/ApplicationServices.h>
#include <stdbool.h>

static double seconds_since_last_input(void) {
	return CGEventSourceSecondsSinceLastEventType(kCGEventSourceStateCombinedSessionState, kCGAnyInputEventType);
}

static bool accessibility_trusted(void) {
	return AXIsProcessTrusted();
}
*/
import "C"

import "time"

type System struct{}

func (System) LastInputAge() (time.Duration, error) {
	seconds := float64(C.seconds_since_last_input())
	return time.Duration(seconds * float64(time.Second)), nil
}

func (System) AccessibilityTrusted() bool {
	return bool(C.accessibility_trusted())
}
